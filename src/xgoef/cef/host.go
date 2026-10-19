// Package cef isolates the Chromium Embedded Framework C API.
//
// The runnable Host is compiled with LLGo (see host_llgo.go) and links the
// small C shim in cef/shim via LLGo's C ABI — not cgo. llcppg.cfg remains the
// path toward generating broader capi bindings later.
package cef

// Settings mirrors the subset of cef_settings_t xgoef needs in phase 1.
type Settings struct {
	NoSandbox      bool
	RootCachePath  string
	SubprocessPath string // native helper for GPU/renderer (not the Go binary)
}

// WindowBounds is the preferred content size for the first browser window.
type WindowBounds struct {
	Width  int
	Height int
}

// Host is the ergonomic surface runtime uses.
// The LLGo build satisfies this via go:linkname to xgoef_cef_* (C ABI).
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

// ErrBindingsMissing is returned when not building with LLGo / CEF host.
var ErrBindingsMissing = errBindings("cef: NativeHost requires LLGo (llgo build); see scripts/build-xgoef-hello.sh")

type errBindings string

func (e errBindings) Error() string { return string(e) }
