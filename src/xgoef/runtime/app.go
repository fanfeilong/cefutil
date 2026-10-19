package runtime

import (
	"fmt"

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

// Run executes the CEF lifecycle.
// Order matches docs/05-lifecycle.md and src/cefdemo/main_linux.cc.
func (a *App) Run() error {
	if code, helper := a.Host.ExecuteProcess(); helper {
		// Secondary process: do not initialize the browser stack.
		return &ExitError{Code: code}
	}
	if err := a.Host.InitAPIHash(); err != nil {
		return err
	}

	settings := cef.Settings{
		NoSandbox:     true,
		RootCachePath: a.Config.CEF.CacheDir,
	}
	if a.Config.CEF.NoSandbox != nil {
		settings.NoSandbox = *a.Config.CEF.NoSandbox
	}

	if err := a.Plugins.Setup(a); err != nil {
		return fmt.Errorf("plugin setup: %w", err)
	}
	if err := a.Host.Initialize(settings); err != nil {
		return err
	}

	win := a.Config.FirstWindow()
	url := win.URL
	if url == "" {
		url = "about:blank"
	}
	bounds := cef.WindowBounds{Width: win.Width, Height: win.Height}
	if bounds.Width <= 0 {
		bounds.Width = 960
	}
	if bounds.Height <= 0 {
		bounds.Height = 640
	}
	if err := a.Host.CreateBrowser(url, bounds); err != nil {
		a.Host.Shutdown()
		return err
	}

	a.Host.RunLoop()
	a.Plugins.Teardown(a)
	a.Host.Shutdown()
	return nil
}

// ExitError means a CEF helper process finished with an exit code.
type ExitError struct {
	Code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("cef helper exit code %d", e.Code)
}
