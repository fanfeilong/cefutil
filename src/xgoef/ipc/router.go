package ipc

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Context is passed to command handlers (window id, etc. later).
type Context struct {
	Command string
}

// Handler is a Tauri-like command implementation.
// argsJSON is the raw JSON object/array from the frontend.
// Return value is JSON-encoded for the response payload.
type Handler func(ctx Context, argsJSON []byte) (any, error)

// Router dispatches invoke(command) to registered handlers.
type Router struct {
	mu       sync.RWMutex
	handlers map[string]Handler
}

// NewRouter creates an empty command router.
func NewRouter() *Router {
	return &Router{handlers: make(map[string]Handler)}
}

// Handle registers or replaces a command.
func (r *Router) Handle(name string, h Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[name] = h
}

// Invoke runs a command by name (used by the CEF message bridge).
func (r *Router) Invoke(name string, argsJSON []byte) ([]byte, error) {
	r.mu.RLock()
	h, ok := r.handlers[name]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("ipc: unknown command %q", name)
	}
	out, err := h(Context{Command: name}, argsJSON)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return []byte("null"), nil
	}
	return json.Marshal(out)
}

// Commands returns registered command names (for capabilities checks).
func (r *Router) Commands() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.handlers))
	for n := range r.handlers {
		names = append(names, n)
	}
	return names
}
