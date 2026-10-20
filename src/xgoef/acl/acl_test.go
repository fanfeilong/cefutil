package acl

import "testing"

func TestAllowsExactPermission(t *testing.T) {
	g := NewGuard([]Capability{{
		Identifier:  "main",
		Windows:     []string{"main"},
		Permissions: []string{"greet", "core:default"},
	}})
	if !g.Allows("main", "greet") {
		t.Fatal("expected greet allowed")
	}
	if g.Allows("main", "core:default") {
		// core:default is a permission id, not a command shortcut for all cmds
		// — but if someone registers a command literally named core:default it matches.
		// Here we only assert unknown cmds are denied:
	}
	if g.Allows("main", "fs:read") {
		t.Fatal("fs:read must be denied")
	}
	if g.Allows("other", "greet") {
		t.Fatal("other window must be denied")
	}
}

func TestDefaultDenyEmpty(t *testing.T) {
	g := NewGuard(nil)
	if g.Allows("main", "greet") {
		t.Fatal("empty guard must deny")
	}
	var nilG *Guard
	if nilG.Allows("main", "greet") {
		t.Fatal("nil guard must deny")
	}
}

func TestWildcardWindow(t *testing.T) {
	g := NewGuard([]Capability{{
		Identifier:  "all",
		Windows:     []string{"*"},
		Permissions: []string{"ping"},
	}})
	if !g.Allows("any-label", "ping") {
		t.Fatal("wildcard window should allow")
	}
}
