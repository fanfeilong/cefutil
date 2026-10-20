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

func (b *NativeBridge) EmitTo(windowLabel, jsonUTF8 string) error { return nil }

func (b *NativeBridge) EmitAll(jsonUTF8 string) error { return nil }

func (b *NativeBridge) SetPendingLabel(label string) {}

func (b *NativeBridge) CreateWindow(spec WindowSpec) error { return nil }

func (b *NativeBridge) CloseWindow(label string) error { return nil }

func (b *NativeBridge) FocusWindow(label string) error { return nil }
