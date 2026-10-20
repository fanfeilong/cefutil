package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fanfeilong/cefutil/src/xgoef/acl"
	"github.com/fanfeilong/cefutil/src/xgoef/cef"
	"github.com/fanfeilong/cefutil/src/xgoef/conf"
	"github.com/fanfeilong/cefutil/src/xgoef/ipc"
	"github.com/fanfeilong/cefutil/src/xgoef/plugin"
)

// App is the Tauri-like application handle.
type App struct {
	Config  conf.File
	Host    cef.Host
	Bridge  cef.Bridge
	Router  *ipc.Router
	Events  *ipc.Bus
	ACL     *acl.Guard
	Plugins *plugin.Registry
}

// New builds an App with defaults. Host may be nil (uses cef.Unimplemented).
// Bridge may be nil (uses cef.NewNativeBridge).
func New(cfg conf.File, host cef.Host) *App {
	if host == nil {
		host = cef.Unimplemented{}
	}
	bridge := cef.NewNativeBridge()
	bus := ipc.NewBus()
	bus.SetEmitter(bridge)
	return &App{
		Config:  cfg,
		Host:    host,
		Bridge:  bridge,
		Router:  ipc.NewRouter(),
		Events:  bus,
		ACL:     acl.NewGuard(cfg.Capabilities),
		Plugins: plugin.NewRegistry(),
	}
}

// Use registers a plugin (Tauri-like builder hook).
func (a *App) Use(p plugin.Plugin) *App {
	a.Plugins.Add(p)
	return a
}

// Handle registers an IPC command on the app router.
func (a *App) Handle(name string, h ipc.Handler) *App {
	a.Router.Handle(name, h)
	return a
}

// Emit sends an event to all webviews (and host listeners).
func (a *App) Emit(name string, payload any) error {
	return a.Events.Emit(name, payload)
}

// EmitTo sends an event to one window label.
func (a *App) EmitTo(windowLabel, name string, payload any) error {
	return a.Events.EmitTo(windowLabel, name, payload)
}

// CreateWindow opens another browser window with the given label.
func (a *App) CreateWindow(spec cef.WindowSpec) error {
	if a.Bridge == nil {
		return fmt.Errorf("runtime: no bridge")
	}
	if spec.URL != "" {
		spec.URL = ResolveURL(spec.URL)
	}
	return a.Bridge.CreateWindow(spec)
}

// CloseWindow closes the window with the given label.
func (a *App) CloseWindow(label string) error {
	if a.Bridge == nil {
		return fmt.Errorf("runtime: no bridge")
	}
	return a.Bridge.CloseWindow(label)
}

// FocusWindow focuses the window with the given label.
func (a *App) FocusWindow(label string) error {
	if a.Bridge == nil {
		return fmt.Errorf("runtime: no bridge")
	}
	return a.Bridge.FocusWindow(label)
}

// Run executes the browser-process side of the CEF lifecycle.
// Helper processes must call Host.ExecuteProcess at the start of main
// (before Run); when Run is entered, this is already the browser process.
// Browser URL is passed as `--url=` for OnContextInitialized.
func (a *App) Run() error {
	if nh, ok := a.Host.(interface{ SetArgs([]string) }); ok {
		nh.SetArgs(buildArgs(a))
	}

	settings := cef.Settings{
		NoSandbox:     true,
		RootCachePath: a.Config.CEF.CacheDir,
	}
	if a.Config.CEF.NoSandbox != nil {
		settings.NoSandbox = *a.Config.CEF.NoSandbox
	}
	if settings.RootCachePath != "" {
		_ = os.MkdirAll(settings.RootCachePath, 0o755)
	}
	settings.SubprocessPath = resolveSubprocess()

	if err := a.Plugins.Setup(a); err != nil {
		return fmt.Errorf("plugin setup: %w", err)
	}

	a.wireIPC()

	if err := a.Host.Initialize(settings); err != nil {
		return err
	}

	win := a.Config.FirstWindow()
	url := ResolveURL(win.URL)
	bounds := cef.WindowBounds{Width: win.Width, Height: win.Height}
	if err := a.Host.CreateBrowser(url, bounds); err != nil {
		a.Host.Shutdown()
		return err
	}

	a.Host.RunLoop()
	a.Plugins.Teardown(a)
	a.Host.Shutdown()
	return nil
}

func (a *App) wireIPC() {
	if a.Bridge == nil {
		return
	}
	label := a.Config.FirstWindow().Label
	if label == "" {
		label = "main"
	}
	a.Bridge.SetPendingLabel(label)

	router := a.Router
	guard := a.ACL
	bridge := a.Bridge
	events := a.Events

	bridge.SetQueryHandler(func(queryID int64, windowLabel, requestJSON string) {
		wl := windowLabel
		if wl == "" {
			wl = label
		}
		resp := handleMessage(guard, router, events, wl, requestJSON)
		bridge.Respond(queryID, true, resp)
	})
}

// handleMessage routes invoke / event envelopes from the frontend.
func handleMessage(guard *acl.Guard, router *ipc.Router, events *ipc.Bus, windowLabel, requestJSON string) string {
	var req struct {
		Type    string          `json:"type"`
		ID      int             `json:"id"`
		Cmd     string          `json:"cmd"`
		Args    json.RawMessage `json:"args"`
		Event   string          `json:"event"`
		Payload json.RawMessage `json:"payload"`
	}
	if err := json.Unmarshal([]byte(requestJSON), &req); err != nil {
		return mustResponse(0, false, nil, "invalid invoke JSON: "+err.Error())
	}
	switch req.Type {
	case "", "invoke":
		return handleInvoke(guard, router, windowLabel, req.ID, req.Cmd, req.Args)
	case "event":
		if events != nil && req.Event != "" {
			var payload any
			if len(req.Payload) > 0 {
				_ = json.Unmarshal(req.Payload, &payload)
			}
			// JS → host: local listeners only (avoid echo loop to all webviews).
			events.NotifyLocal(req.Event, payload)
		}
		return mustResponse(req.ID, true, map[string]any{"ok": true}, "")
	default:
		return mustResponse(req.ID, false, nil, "unsupported message type "+req.Type)
	}
}

func handleInvoke(guard *acl.Guard, router *ipc.Router, windowLabel string, id int, cmd string, argsJSON json.RawMessage) string {
	if cmd == "" {
		return mustResponse(id, false, nil, "missing cmd")
	}
	if guard == nil || !guard.Allows(windowLabel, cmd) {
		return mustResponse(id, false, nil, (&acl.DenyError{Window: windowLabel, Command: cmd}).Error())
	}
	args := []byte(argsJSON)
	if len(args) == 0 {
		args = []byte("{}")
	}
	out, err := router.Invoke(cmd, args)
	if err != nil {
		return mustResponse(id, false, nil, err.Error())
	}
	var payload any
	if err := json.Unmarshal(out, &payload); err != nil {
		return mustResponse(id, false, nil, "marshal response: "+err.Error())
	}
	return mustResponse(id, true, payload, "")
}

func mustResponse(id int, ok bool, payload any, errMsg string) string {
	msg := map[string]any{
		"type": "response",
		"id":   id,
		"ok":   ok,
	}
	if ok {
		msg["payload"] = payload
	} else {
		msg["error"] = errMsg
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return `{"type":"response","id":0,"ok":false,"error":"response encode failed"}`
	}
	return string(b)
}

func buildArgs(a *App) []string {
	win := a.Config.FirstWindow()
	url := ResolveURL(win.URL)
	label := win.Label
	if label == "" {
		label = "main"
	}
	args := append([]string{}, os.Args...)
	args = append(args, "--url="+url, "--xgoef-label="+label)
	// LLGo (and Go) hosts are poor CEF GPU/renderer subprocesses.
	// Prefer a native helper via SubprocessPath; these flags keep the
	// browser process alive when the helper still cannot use GPU.
	if os.Getenv("XGOEF_ENABLE_GPU") != "1" {
		args = append(args, "--disable-gpu", "--in-process-gpu")
	}
	return args
}

func resolveSubprocess() string {
	if p := os.Getenv("XGOEF_HELPER"); p != "" {
		if abs, err := filepath.Abs(p); err == nil {
			return abs
		}
		return p
	}
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	candidate := filepath.Join(filepath.Dir(exe), "xgoef-helper")
	if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
		return candidate
	}
	return ""
}

// ResolveURL turns relative file paths into absolute file:// URLs.
func ResolveURL(u string) string {
	if u == "" {
		return "about:blank"
	}
	switch {
	case u == "about:blank", strings.HasPrefix(u, "http://"), strings.HasPrefix(u, "https://"):
		return u
	case strings.HasPrefix(u, "file://"):
		path := strings.TrimPrefix(u, "file://")
		abs, err := filepath.Abs(path)
		if err != nil {
			return u
		}
		return "file://" + abs
	default:
		abs, err := filepath.Abs(u)
		if err != nil {
			return u
		}
		return "file://" + abs
	}
}

// ExitError means a CEF helper process finished with an exit code.
type ExitError struct {
	Code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("cef helper exit code %d", e.Code)
}
