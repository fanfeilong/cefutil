package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fanfeilong/cefutil/src/xgoef/cef"
	"github.com/fanfeilong/cefutil/src/xgoef/conf"
	"github.com/fanfeilong/cefutil/src/xgoef/ipc"
	"github.com/fanfeilong/cefutil/src/xgoef/plugin"
)

// App is the Tauri-like application handle.
type App struct {
	Config  conf.File
	Host    cef.Host
	Router  *ipc.Router
	Plugins *plugin.Registry
}

// New builds an App with defaults. Host may be nil (uses cef.Unimplemented).
func New(cfg conf.File, host cef.Host) *App {
	if host == nil {
		host = cef.Unimplemented{}
	}
	return &App{
		Config:  cfg,
		Host:    host,
		Router:  ipc.NewRouter(),
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

func buildArgs(a *App) []string {
	win := a.Config.FirstWindow()
	url := ResolveURL(win.URL)
	args := append([]string{}, os.Args...)
	args = append(args, "--url="+url)
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
