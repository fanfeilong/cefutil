package runtime

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/fanfeilong/cefutil/src/xgoef/acl"
	"github.com/fanfeilong/cefutil/src/xgoef/ipc"
)

func TestHandleInvokeGreet(t *testing.T) {
	r := ipc.NewRouter()
	r.Handle("greet", func(ctx ipc.Context, argsJSON []byte) (any, error) {
		var args struct {
			Name string `json:"name"`
		}
		_ = json.Unmarshal(argsJSON, &args)
		return "hello, " + args.Name + " (from xgoef)", nil
	})
	g := acl.NewGuard([]acl.Capability{{
		Identifier:  "main",
		Windows:     []string{"main"},
		Permissions: []string{"greet"},
	}})
	out := handleMessage(g, r, nil, "main", `{"type":"invoke","id":1,"cmd":"greet","args":{"name":"xgo"}}`)
	var resp struct {
		OK      bool   `json:"ok"`
		Payload string `json:"payload"`
		Error   string `json:"error"`
	}
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatal(err)
	}
	if !resp.OK || resp.Payload != "hello, xgo (from xgoef)" {
		t.Fatalf("%s", out)
	}
}

func TestHandleInvokeACLDeny(t *testing.T) {
	r := ipc.NewRouter()
	r.Handle("greet", func(ctx ipc.Context, argsJSON []byte) (any, error) {
		return "ok", nil
	})
	g := acl.NewGuard([]acl.Capability{{
		Identifier:  "main",
		Windows:     []string{"main"},
		Permissions: []string{"greet"},
	}})
	out := handleMessage(g, r, nil, "main", `{"type":"invoke","id":2,"cmd":"secret","args":{}}`)
	if !strings.Contains(out, `"ok":false`) || !strings.Contains(out, "not permitted") {
		t.Fatalf("%s", out)
	}
}

func TestHandleEventToHost(t *testing.T) {
	bus := ipc.NewBus()
	got := ""
	bus.Listen("ping", func(payload any) {
		got, _ = payload.(string)
	})
	out := handleMessage(nil, ipc.NewRouter(), bus, "main",
		`{"type":"event","id":3,"event":"ping","payload":"from-js"}`)
	if !strings.Contains(out, `"ok":true`) {
		t.Fatalf("%s", out)
	}
	if got != "from-js" {
		t.Fatalf("got %q", got)
	}
}
