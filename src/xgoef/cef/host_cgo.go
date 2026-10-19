//go:build cgo && linux

package cef

/*
#cgo CFLAGS: -std=c11 -DCEF_X11 -D_GNU_SOURCE
#cgo LDFLAGS: -lX11 -lpthread -ldl -lrt -lstdc++

#include <stdlib.h>
#include "shim/xgoef_host.h"
*/
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

// NativeHost is the cgo-backed Host linking libcef via xgoef_cef.
type NativeHost struct {
	argv []*C.char
}

// NewNativeHost prepares argv for CEF (caller should include --url=).
func NewNativeHost() *NativeHost {
	return &NativeHost{}
}

func (h *NativeHost) cArgs() (C.int, **C.char) {
	if len(h.argv) == 0 {
		args := os.Args
		h.argv = make([]*C.char, len(args)+1)
		for i, a := range args {
			h.argv[i] = C.CString(a)
		}
		h.argv[len(args)] = nil
	}
	return C.int(len(h.argv) - 1), (**C.char)(unsafe.Pointer(&h.argv[0]))
}

// SetArgs replaces process argv used for CEF (must include program name).
func (h *NativeHost) SetArgs(args []string) {
	for _, p := range h.argv {
		if p != nil {
			C.free(unsafe.Pointer(p))
		}
	}
	h.argv = make([]*C.char, len(args)+1)
	for i, a := range args {
		h.argv[i] = C.CString(a)
	}
	h.argv[len(args)] = nil
}

func (h *NativeHost) ExecuteProcess() (int, bool) {
	argc, argv := h.cArgs()
	code := int(C.xgoef_cef_execute_process(argc, argv))
	if code >= 0 {
		return code, true
	}
	return -1, false
}

func (h *NativeHost) InitAPIHash() error {
	C.xgoef_cef_api_hash()
	return nil
}

func (h *NativeHost) Initialize(settings Settings) error {
	argc, argv := h.cArgs()
	var cache *C.char
	if settings.RootCachePath != "" {
		cache = C.CString(settings.RootCachePath)
		defer C.free(unsafe.Pointer(cache))
	}
	var sub *C.char
	if settings.SubprocessPath != "" {
		sub = C.CString(settings.SubprocessPath)
		defer C.free(unsafe.Pointer(sub))
	}
	noSandbox := C.int(0)
	if settings.NoSandbox {
		noSandbox = 1
	}
	if rc := C.xgoef_cef_initialize(argc, argv, cache, noSandbox, sub); rc != 0 {
		return fmt.Errorf("cef: initialize failed (exit code %d)", int(rc))
	}
	return nil
}

func (h *NativeHost) CreateBrowser(url string, bounds WindowBounds) error {
	// Browser is created in OnContextInitialized from argv --url=.
	_ = url
	_ = bounds
	return nil
}

func (h *NativeHost) RunLoop() { C.xgoef_cef_run_loop() }

func (h *NativeHost) Shutdown() { C.xgoef_cef_shutdown() }
