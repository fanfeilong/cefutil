package ipc

import "sync"

// Listener receives event payloads from host or webview.
type Listener func(payload any)

// Bus is a simple in-process event hub (phase 1).
// Later, Emit will also fan-out to the webview via CEF process messages.
type Bus struct {
	mu   sync.RWMutex
	subs map[string][]Listener
}

// NewBus creates an event bus.
func NewBus() *Bus {
	return &Bus{subs: make(map[string][]Listener)}
}

// Listen registers a listener for event name.
func (b *Bus) Listen(name string, l Listener) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs[name] = append(b.subs[name], l)
}

// Emit notifies all listeners of name.
func (b *Bus) Emit(name string, payload any) {
	b.mu.RLock()
	list := append([]Listener(nil), b.subs[name]...)
	b.mu.RUnlock()
	for _, l := range list {
		l(payload)
	}
}
