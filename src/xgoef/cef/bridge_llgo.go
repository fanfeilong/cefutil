//go:build llgo

package cef

import (
	"fmt"
	"unsafe"

	"github.com/goplus/lib/c"
)

var _ = unsafe.Sizeof(0)

//llgo:cdecl
type ipcQueryFn func(queryID int64, windowLabel, request *c.Char)

// NativeBridge is the LLGo-backed IPC bridge to the C++ MessageRouter shim.
type NativeBridge struct {
	handler QueryHandler
}

// activeBridge is the process-wide bridge registered with the C shim.
// C function pointers cannot carry a Go receiver, so we keep one global.
var activeBridge *NativeBridge

// NewNativeBridge prepares the CEF IPC bridge (call SetQueryHandler before RunLoop).
func NewNativeBridge() *NativeBridge {
	return &NativeBridge{}
}

// SetQueryHandler registers the Go callback invoked for each cefQuery.
func (b *NativeBridge) SetQueryHandler(h QueryHandler) {
	b.handler = h
	activeBridge = b
	if h == nil {
		xgoefIPCSetHandler(nil)
		return
	}
	xgoefIPCSetHandler(nativeOnQuery)
}

//llgo:cdecl
func nativeOnQuery(queryID int64, windowLabel, request *c.Char) {
	b := activeBridge
	if b == nil || b.handler == nil {
		xgoefIPCRespond(queryID, 0, c.AllocaCStr("xgoef host bridge not connected"))
		return
	}
	label := "main"
	if windowLabel != nil {
		if s := c.GoString(windowLabel); s != "" {
			label = s
		}
	}
	req := ""
	if request != nil {
		req = c.GoString(request)
	}
	b.handler(queryID, label, req)
}

// Respond completes a pending query (ok → Success JSON; else Failure message).
func (b *NativeBridge) Respond(queryID int64, ok bool, resultUTF8 string) {
	flag := c.Int(0)
	if ok {
		flag = 1
	}
	xgoefIPCRespond(queryID, flag, c.AllocaCStr(resultUTF8))
}

func (b *NativeBridge) EmitTo(windowLabel, jsonUTF8 string) error {
	if rc := xgoefIPCEmitTo(c.AllocaCStr(windowLabel), c.AllocaCStr(jsonUTF8)); rc != 0 {
		return fmt.Errorf("cef: emit_to %q failed (%d)", windowLabel, int(rc))
	}
	return nil
}

func (b *NativeBridge) EmitAll(jsonUTF8 string) error {
	if rc := xgoefIPCEmitAll(c.AllocaCStr(jsonUTF8)); rc != 0 {
		return fmt.Errorf("cef: emit_all failed (%d)", int(rc))
	}
	return nil
}

func (b *NativeBridge) SetPendingLabel(label string) {
	xgoefWindowSetPendingLabel(c.AllocaCStr(label))
}

func (b *NativeBridge) CreateWindow(spec WindowSpec) error {
	if spec.Label == "" {
		spec.Label = "main"
	}
	if spec.Width <= 0 {
		spec.Width = 960
	}
	if spec.Height <= 0 {
		spec.Height = 640
	}
	rc := xgoefWindowCreate(
		c.AllocaCStr(spec.Label),
		c.AllocaCStr(spec.Title),
		c.AllocaCStr(spec.URL),
		c.Int(spec.Width),
		c.Int(spec.Height),
	)
	if rc != 0 {
		return fmt.Errorf("cef: create window %q failed (%d)", spec.Label, int(rc))
	}
	return nil
}

func (b *NativeBridge) CloseWindow(label string) error {
	if rc := xgoefWindowClose(c.AllocaCStr(label)); rc != 0 {
		return fmt.Errorf("cef: close window %q failed (%d)", label, int(rc))
	}
	return nil
}

func (b *NativeBridge) FocusWindow(label string) error {
	if rc := xgoefWindowFocus(c.AllocaCStr(label)); rc != 0 {
		return fmt.Errorf("cef: focus window %q failed (%d)", label, int(rc))
	}
	return nil
}

//go:linkname xgoefIPCSetHandler C.xgoef_ipc_set_handler
func xgoefIPCSetHandler(fn ipcQueryFn)

//go:linkname xgoefIPCRespond C.xgoef_ipc_respond
func xgoefIPCRespond(queryID int64, ok c.Int, result *c.Char)

//go:linkname xgoefIPCEmitTo C.xgoef_ipc_emit_to
func xgoefIPCEmitTo(label, json *c.Char) c.Int

//go:linkname xgoefIPCEmitAll C.xgoef_ipc_emit_all
func xgoefIPCEmitAll(json *c.Char) c.Int

//go:linkname xgoefWindowSetPendingLabel C.xgoef_window_set_pending_label
func xgoefWindowSetPendingLabel(label *c.Char)

//go:linkname xgoefWindowCreate C.xgoef_window_create
func xgoefWindowCreate(label, title, url *c.Char, width, height c.Int) c.Int

//go:linkname xgoefWindowClose C.xgoef_window_close
func xgoefWindowClose(label *c.Char) c.Int

//go:linkname xgoefWindowFocus C.xgoef_window_focus
func xgoefWindowFocus(label *c.Char) c.Int
