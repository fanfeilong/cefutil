package ipc

import (
	"encoding/json"
	"testing"
)

func TestRouterInvoke(t *testing.T) {
	r := NewRouter()
	r.Handle("greet", func(ctx Context, argsJSON []byte) (any, error) {
		var args struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(argsJSON, &args); err != nil {
			return nil, err
		}
		return "hi " + args.Name, nil
	})

	out, err := r.Invoke("greet", []byte(`{"name":"xgo"}`))
	if err != nil {
		t.Fatal(err)
	}
	var s string
	if err := json.Unmarshal(out, &s); err != nil {
		t.Fatal(err)
	}
	if s != "hi xgo" {
		t.Fatalf("got %q", s)
	}
}

func TestRouterUnknown(t *testing.T) {
	r := NewRouter()
	if _, err := r.Invoke("nope", nil); err == nil {
		t.Fatal("expected error")
	}
}
