package conf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadHelloConfig(t *testing.T) {
	path := filepath.Join("..", "examples", "hello", "xgoef.conf.json")
	f, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if f.ProductName != "xgoef-hello" {
		t.Fatalf("productName=%q", f.ProductName)
	}
	w := f.FirstWindow()
	if w.Width != 960 || w.Height != 640 {
		t.Fatalf("window=%+v", w)
	}
	if w.Label != "main" {
		t.Fatalf("label=%q", w.Label)
	}
	if f.CEF.RootEnv != "CEF_ROOT" {
		t.Fatalf("rootEnv=%q", f.CEF.RootEnv)
	}
	if len(f.Capabilities) == 0 {
		t.Fatal("expected capabilities")
	}
	found := false
	for _, c := range f.Capabilities {
		for _, p := range c.Permissions {
			if p == "greet" {
				found = true
			}
		}
	}
	if !found {
		t.Fatalf("capabilities=%+v", f.Capabilities)
	}
}

func TestLoadCapabilitiesArray(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ok.json")
	body := `{
	  "productName":"t",
	  "app":{"windows":[{"label":"main","title":"t","width":1,"height":1}]},
	  "capabilities":[{"identifier":"main","windows":["main"],"permissions":["greet"]}]
	}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	f, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.Capabilities) != 1 || f.Capabilities[0].Permissions[0] != "greet" {
		t.Fatalf("%+v", f.Capabilities)
	}
}

func TestLoadMissingProductName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte(`{"app":{}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("expected error")
	}
}
