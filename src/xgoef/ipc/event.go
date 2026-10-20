package ipc

import (
	"encoding/json"
	"sync"
)

// Listener receives event payloads from host or webview.
type Listener func(payload any)

// Emitter delivers a JSON event envelope to one or more webviews.
type Emitter interface {
	EmitTo(windowLabel, jsonUTF8 string) error
	EmitAll(jsonUTF8 string) error
}

// Bus is an event hub: in-process listeners + optional webview fan-out.
type Bus struct {
	mu      sync.RWMutex
	subs    map[string][]Listener
	emitter Emitter
}

// NewBus creates an event bus.
func NewBus() *Bus {
	return &Bus{subs: make(map[string][]Listener)}
}

// SetEmitter attaches the CEF bridge used for cross-process delivery.
func (b *Bus) SetEmitter(e Emitter) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.emitter = e
}

// Listen registers a host-side listener for event name.
func (b *Bus) Listen(name string, l Listener) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs[name] = append(b.subs[name], l)
}

// Emit notifies in-process listeners and all webviews.
func (b *Bus) Emit(name string, payload any) error {
	b.NotifyLocal(name, payload)
	return b.deliver("", name, payload, true)
}

// EmitTo notifies in-process listeners and one window label.
func (b *Bus) EmitTo(windowLabel, name string, payload any) error {
	b.NotifyLocal(name, payload)
	return b.deliver(windowLabel, name, payload, false)
}

// NotifyLocal runs host-side listeners only (no webview fan-out).
func (b *Bus) NotifyLocal(name string, payload any) {
	b.mu.RLock()
	list := append([]Listener(nil), b.subs[name]...)
	b.mu.RUnlock()
	for _, l := range list {
		l(payload)
	}
}

func (b *Bus) deliver(windowLabel, name string, payload any, all bool) error {
	b.mu.RLock()
	em := b.emitter
	b.mu.RUnlock()
	if em == nil {
		return nil
	}
	env := map[string]any{
		"type":    "event",
		"event":   name,
		"payload": payload,
	}
	raw, err := json.Marshal(env)
	if err != nil {
		return err
	}
	if all {
		return em.EmitAll(string(raw))
	}
	return em.EmitTo(windowLabel, string(raw))
}
