// Package acl enforces Tauri-like capabilities for IPC commands.
package acl

import "fmt"

// Capability binds a set of permissions to one or more window labels.
type Capability struct {
	Identifier  string   `json:"identifier"`
	Windows     []string `json:"windows"`
	Permissions []string `json:"permissions"`
}

// Guard checks whether a window may invoke a command.
// Default deny: no matching capability ⇒ reject.
type Guard struct {
	caps []Capability
}

// NewGuard builds a Guard from the configured capability list.
func NewGuard(caps []Capability) *Guard {
	out := make([]Capability, len(caps))
	copy(out, caps)
	return &Guard{caps: out}
}

// Allows reports whether windowLabel may invoke cmd.
// Permission ids match command names (e.g. "greet") or the wildcard "core:default"
// which does not grant arbitrary commands by itself — only exact permission
// strings authorize a command.
func (g *Guard) Allows(windowLabel, cmd string) bool {
	if g == nil || cmd == "" {
		return false
	}
	if windowLabel == "" {
		windowLabel = "main"
	}
	for _, cap := range g.caps {
		if !windowMatches(cap.Windows, windowLabel) {
			continue
		}
		for _, p := range cap.Permissions {
			if p == cmd {
				return true
			}
		}
	}
	return false
}

// DenyError is returned when ACL rejects an invoke.
type DenyError struct {
	Window  string
	Command string
}

func (e *DenyError) Error() string {
	return fmt.Sprintf("acl: command %q not permitted for window %q", e.Command, e.Window)
}

func windowMatches(windows []string, label string) bool {
	if len(windows) == 0 {
		// Empty windows list ⇒ capability applies to all labels.
		return true
	}
	for _, w := range windows {
		if w == label || w == "*" {
			return true
		}
	}
	return false
}
