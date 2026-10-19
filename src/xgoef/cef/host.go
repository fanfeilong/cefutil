// Package cef isolates the Chromium Embedded Framework C API.
//
// Bindings are intended to be produced by llcppg (see llcppg.cfg).
// Until generation is wired to a local CEF_ROOT, runtime depends on the
// Host interface below so the rest of xgoef stays testable.
package cef

// Settings mirrors the subset of cef_settings_t xgoef needs in phase 1.
type Settings struct {
	NoSandbox    bool
	RootCachePath string
}

// WindowBounds is the preferred content size for the first browser window.
type WindowBounds struct {
	Width  int
	Height int
}

// Host is the ergonomic surface runtime uses.
// A future llcppg-backed implementation will satisfy this interface
// by calling into libcef via LLGo C ABI.
type Host interface {
	// ExecuteProcess runs CEF secondary-process logic.
	// If this process is a helper, ok is true and code is the process exit code.
	ExecuteProcess() (code int, helper bool)

	// InitAPIHash must run before any other CEF call (C API requirement).
	InitAPIHash() error

	// Initialize / RunLoop / Shutdown follow docs/05-lifecycle.md.
	Initialize(settings Settings) error
	CreateBrowser(url string, bounds WindowBounds) error
	RunLoop()
	Shutdown()
}

// Unimplemented is a stub Host used until llcppg bindings exist.
type Unimplemented struct{}

func (Unimplemented) ExecuteProcess() (int, bool) { return -1, false }

func (Unimplemented) InitAPIHash() error {
	return ErrBindingsMissing
}

func (Unimplemented) Initialize(Settings) error { return ErrBindingsMissing }

func (Unimplemented) CreateBrowser(string, WindowBounds) error { return ErrBindingsMissing }

func (Unimplemented) RunLoop() {}

func (Unimplemented) Shutdown() {}

// ErrBindingsMissing is returned when CEF symbols are not linked yet.
var ErrBindingsMissing = errBindings("cef: llcppg bindings not generated; set CEF_ROOT and run llcppg in cef/")

type errBindings string

func (e errBindings) Error() string { return string(e) }
