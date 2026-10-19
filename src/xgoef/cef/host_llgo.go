//go:build llgo

package cef

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/goplus/lib/c"
)

// LLGoPackage links the native CEF host shim + libcef.
// Paths come from pkg-config (see scripts/build-xgoef-hello.sh).
//
// Do not put -Wl,-rpath,$ORIGIN (or any $VAR) in the .pc Libs line: LLGo runs
// os.Expand on pkg-config output, which empties $ORIGIN and leaves
// "-Wl,-rpath,", which then consumes the following -lffi as an rpath entry.
// LLGo already turns every -L into a runtime -rpath.
const LLGoPackage = "link: $(pkg-config --libs xgoef-cef);"

// NativeHost is the LLGo-backed Host (no cgo).
// argv is heap-allocated C strings kept for CEF main_args lifetime.
type NativeHost struct {
	args []string
	argc c.Int
	argv **c.Char
	cstr []*c.Char // owned; freed on SetArgs / replace
}

// NewNativeHost prepares a CEF host for LLGo builds.
func NewNativeHost() *NativeHost {
	return &NativeHost{args: append([]string(nil), os.Args...)}
}

// SetArgs replaces process argv used for CEF (must include program name).
func (h *NativeHost) SetArgs(args []string) {
	h.args = append([]string(nil), args...)
	h.rebuildArgv()
}

func (h *NativeHost) freeArgv() {
	for _, p := range h.cstr {
		if p != nil {
			c.Free(c.Pointer(p))
		}
	}
	h.cstr = nil
	if h.argv != nil {
		c.Free(c.Pointer(h.argv))
		h.argv = nil
	}
	h.argc = 0
}

func (h *NativeHost) rebuildArgv() {
	h.freeArgv()
	n := len(h.args)
	// argv has n entries + trailing NULL (CEF main_args convention).
	slotBytes := uintptr(n+1) * unsafe.Sizeof((*c.Char)(nil))
	raw := c.Malloc(slotBytes)
	c.Memset(raw, 0, slotBytes)
	slots := unsafe.Slice((**c.Char)(raw), n+1)
	h.cstr = make([]*c.Char, n)
	for i, s := range h.args {
		dup := c.Strdup(c.AllocaCStr(s))
		h.cstr[i] = dup
		slots[i] = dup
	}
	slots[n] = nil
	h.argv = (**c.Char)(raw)
	h.argc = c.Int(n)
}

func (h *NativeHost) cArgs() (c.Int, **c.Char) {
	if h.argv == nil {
		h.rebuildArgv()
	}
	return h.argc, h.argv
}

func (h *NativeHost) ExecuteProcess() (int, bool) {
	argc, argv := h.cArgs()
	code := int(xgoefExecuteProcess(argc, argv))
	if code >= 0 {
		return code, true
	}
	return -1, false
}

func (h *NativeHost) InitAPIHash() error {
	xgoefAPIHash()
	return nil
}

func (h *NativeHost) Initialize(settings Settings) error {
	argc, argv := h.cArgs()
	var cache *c.Char
	if settings.RootCachePath != "" {
		cache = c.Strdup(c.AllocaCStr(settings.RootCachePath))
		defer c.Free(c.Pointer(cache))
	}
	var sub *c.Char
	if settings.SubprocessPath != "" {
		sub = c.Strdup(c.AllocaCStr(settings.SubprocessPath))
		defer c.Free(c.Pointer(sub))
	}
	noSandbox := c.Int(0)
	if settings.NoSandbox {
		noSandbox = 1
	}
	if rc := xgoefInitialize(argc, argv, cache, noSandbox, sub); rc != 0 {
		return fmt.Errorf("cef: initialize failed (exit code %d)", int(rc))
	}
	return nil
}

func (h *NativeHost) CreateBrowser(url string, bounds WindowBounds) error {
	_ = url
	_ = bounds
	return nil
}

func (h *NativeHost) RunLoop()  { xgoefRunLoop() }
func (h *NativeHost) Shutdown() { xgoefShutdown() }

//go:linkname xgoefAPIHash C.xgoef_cef_api_hash
func xgoefAPIHash()

//go:linkname xgoefExecuteProcess C.xgoef_cef_execute_process
func xgoefExecuteProcess(argc c.Int, argv **c.Char) c.Int

//go:linkname xgoefInitialize C.xgoef_cef_initialize
func xgoefInitialize(argc c.Int, argv **c.Char, rootCache *c.Char, noSandbox c.Int, subprocess *c.Char) c.Int

//go:linkname xgoefRunLoop C.xgoef_cef_run_loop
func xgoefRunLoop()

//go:linkname xgoefShutdown C.xgoef_cef_shutdown
func xgoefShutdown()
