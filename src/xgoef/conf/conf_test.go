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
	if f.CEF.RootEnv != "CEF_ROOT" {
		t.Fatalf("rootEnv=%q", f.CEF.RootEnv)
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
