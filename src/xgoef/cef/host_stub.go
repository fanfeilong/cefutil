//go:build !cgo || !linux

package cef

// NewNativeHost returns nil when cgo/linux CEF host is unavailable.
func NewNativeHost() *NativeHost { return nil }

// NativeHost is a stub so non-cgo builds still typecheck call sites.
type NativeHost struct{}

func (h *NativeHost) ExecuteProcess() (int, bool)          { return -1, false }
func (h *NativeHost) InitAPIHash() error                   { return ErrBindingsMissing }
func (h *NativeHost) Initialize(Settings) error            { return ErrBindingsMissing }
func (h *NativeHost) CreateBrowser(string, WindowBounds) error {
	return ErrBindingsMissing
}
func (h *NativeHost) RunLoop()  {}
func (h *NativeHost) Shutdown() {}
