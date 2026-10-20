package ipc

import "testing"

type memEmitter struct {
	to   []string
	all  []string
	last string
}

func (m *memEmitter) EmitTo(label, jsonUTF8 string) error {
	m.to = append(m.to, label)
	m.last = jsonUTF8
	return nil
}

func (m *memEmitter) EmitAll(jsonUTF8 string) error {
	m.all = append(m.all, jsonUTF8)
	m.last = jsonUTF8
	return nil
}

func TestBusEmitFansOut(t *testing.T) {
	bus := NewBus()
	em := &memEmitter{}
	bus.SetEmitter(em)
	local := 0
	bus.Listen("tick", func(payload any) { local++ })
	if err := bus.Emit("tick", map[string]any{"n": 1}); err != nil {
		t.Fatal(err)
	}
	if local != 1 || len(em.all) != 1 {
		t.Fatalf("local=%d all=%v", local, em.all)
	}
	if err := bus.EmitTo("main", "tick", 2); err != nil {
		t.Fatal(err)
	}
	if len(em.to) != 1 || em.to[0] != "main" {
		t.Fatalf("%v", em.to)
	}
}

func TestNotifyLocalNoFanOut(t *testing.T) {
	bus := NewBus()
	em := &memEmitter{}
	bus.SetEmitter(em)
	n := 0
	bus.Listen("ping", func(payload any) { n++ })
	bus.NotifyLocal("ping", "x")
	if n != 1 || len(em.all) != 0 || len(em.to) != 0 {
		t.Fatalf("n=%d em=%+v", n, em)
	}
}
