//go:build !llgo

package cef

// NativeBridge is a no-op stub for gc builds / unit tests.
type NativeBridge struct {
	handler QueryHandler
}

// NewNativeBridge returns a stub bridge when not building with LLGo.
func NewNativeBridge() *NativeBridge {
	return &NativeBridge{}
}

func (b *NativeBridge) SetQueryHandler(h QueryHandler) { b.handler = h }

func (b *NativeBridge) Respond(queryID int64, ok bool, resultUTF8 string) {}
