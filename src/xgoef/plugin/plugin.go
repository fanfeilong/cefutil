package plugin

import "fmt"

// Plugin mirrors Tauri's plugin lifecycle at a minimal level.
// Setup/Teardown receive the *runtime.App as any to avoid an import cycle.
type Plugin interface {
	Name() string
	Setup(app any) error
	Teardown(app any) error
}

// Registry holds plugins in registration order.
type Registry struct {
	list []Plugin
}

// NewRegistry creates an empty plugin registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Add appends a plugin.
func (r *Registry) Add(p Plugin) {
	r.list = append(r.list, p)
}

// Setup runs all plugins.
func (r *Registry) Setup(app any) error {
	for _, p := range r.list {
		if err := p.Setup(app); err != nil {
			return fmt.Errorf("plugin %s: %w", p.Name(), err)
		}
	}
	return nil
}

// Teardown runs plugins in reverse order.
func (r *Registry) Teardown(app any) {
	for i := len(r.list) - 1; i >= 0; i-- {
		_ = r.list[i].Teardown(app)
	}
}
