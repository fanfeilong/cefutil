//go:build llgo

package cef

import (
	"unsafe"

	"github.com/goplus/lib/c"
)

var _ = unsafe.Sizeof(0)

//llgo:cdecl
type ipcQueryFn func(queryID int64, request *c.Char)

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
func nativeOnQuery(queryID int64, request *c.Char) {
	b := activeBridge
	if b == nil || b.handler == nil {
		xgoefIPCRespond(queryID, 0, c.AllocaCStr("xgoef host bridge not connected"))
		return
	}
	req := ""
	if request != nil {
		req = c.GoString(request)
	}
	b.handler(queryID, req)
}

// Respond completes a pending query (ok → Success JSON; else Failure message).
func (b *NativeBridge) Respond(queryID int64, ok bool, resultUTF8 string) {
	flag := c.Int(0)
	if ok {
		flag = 1
	}
	// C++ IpcRespond copies into std::string before any async post.
	xgoefIPCRespond(queryID, flag, c.AllocaCStr(resultUTF8))
}

//go:linkname xgoefIPCSetHandler C.xgoef_ipc_set_handler
func xgoefIPCSetHandler(fn ipcQueryFn)

//go:linkname xgoefIPCRespond C.xgoef_ipc_respond
func xgoefIPCRespond(queryID int64, ok c.Int, result *c.Char)
