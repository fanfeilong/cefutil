//go:build !llgo

package cef

// NewNativeHost returns nil when not building with LLGo.
func NewNativeHost() *NativeHost { return nil }

// NativeHost is a stub so `go test` (gc) still typechecks call sites.
type NativeHost struct{}

func (h *NativeHost) SetArgs([]string) {}

func (h *NativeHost) ExecuteProcess() (int, bool) { return -1, false }
func (h *NativeHost) InitAPIHash() error          { return ErrBindingsMissing }
func (h *NativeHost) Initialize(Settings) error   { return ErrBindingsMissing }
func (h *NativeHost) CreateBrowser(string, WindowBounds) error {
	return ErrBindingsMissing
}
func (h *NativeHost) RunLoop()  {}
func (h *NativeHost) Shutdown() {}
