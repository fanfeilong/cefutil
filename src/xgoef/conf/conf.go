package conf

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fanfeilong/cefutil/src/xgoef/acl"
)

// File is the on-disk xgoef.conf.json shape (Tauri-like).
type File struct {
	ProductName  string           `json:"productName"`
	Version      string           `json:"version"`
	Identifier   string           `json:"identifier"`
	Build        Build            `json:"build"`
	App          AppSection       `json:"app"`
	CEF          CEFSection       `json:"cef"`
	Plugins      []string         `json:"plugins"`
	Capabilities []acl.Capability `json:"capabilities"`
}

// Build holds frontend paths.
type Build struct {
	FrontendDist string `json:"frontendDist"`
	DevURL       string `json:"devUrl"`
}

// AppSection holds window definitions.
type AppSection struct {
	Windows []Window `json:"windows"`
}

// Window is one browser window.
type Window struct {
	Label  string `json:"label"`
	Title  string `json:"title"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

// CEFSection points at the binary distribution.
type CEFSection struct {
	RootEnv   string `json:"rootEnv"`
	NoSandbox *bool  `json:"noSandbox"`
	CacheDir  string `json:"cacheDir"`
}

// rawFile accepts both the P0 capabilities array and the legacy map form.
type rawFile struct {
	ProductName  string          `json:"productName"`
	Version      string          `json:"version"`
	Identifier   string          `json:"identifier"`
	Build        Build           `json:"build"`
	App          AppSection      `json:"app"`
	CEF          CEFSection      `json:"cef"`
	Plugins      []string        `json:"plugins"`
	Capabilities json.RawMessage `json:"capabilities"`
}

// Load reads and decodes a config file.
func Load(path string) (File, error) {
	var f File
	data, err := os.ReadFile(path)
	if err != nil {
		return f, err
	}
	var raw rawFile
	if err := json.Unmarshal(data, &raw); err != nil {
		return f, fmt.Errorf("conf: %w", err)
	}
	f.ProductName = raw.ProductName
	f.Version = raw.Version
	f.Identifier = raw.Identifier
	f.Build = raw.Build
	f.App = raw.App
	f.CEF = raw.CEF
	f.Plugins = raw.Plugins

	caps, err := parseCapabilities(raw.Capabilities)
	if err != nil {
		return f, err
	}
	f.Capabilities = caps

	if f.ProductName == "" {
		return f, fmt.Errorf("conf: productName is required")
	}
	if f.CEF.RootEnv == "" {
		f.CEF.RootEnv = "CEF_ROOT"
	}
	if f.Build.FrontendDist == "" {
		f.Build.FrontendDist = "web"
	}
	for i := range f.App.Windows {
		if f.App.Windows[i].Label == "" {
			f.App.Windows[i].Label = "main"
		}
	}
	return f, nil
}

func parseCapabilities(raw json.RawMessage) ([]acl.Capability, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	// Prefer Tauri-like array form.
	var arr []acl.Capability
	if err := json.Unmarshal(raw, &arr); err == nil {
		return arr, nil
	}
	// Legacy scaffold: { "main": { "commands": ["greet"] } }
	var legacy map[string]struct {
		Commands []string `json:"commands"`
	}
	if err := json.Unmarshal(raw, &legacy); err != nil {
		return nil, fmt.Errorf("conf: capabilities: %w", err)
	}
	out := make([]acl.Capability, 0, len(legacy))
	for id, v := range legacy {
		out = append(out, acl.Capability{
			Identifier:  id,
			Windows:     []string{id},
			Permissions: append([]string(nil), v.Commands...),
		})
	}
	return out, nil
}

// FirstWindow returns the first window or a default.
func (f File) FirstWindow() Window {
	if len(f.App.Windows) == 0 {
		return Window{Label: "main", Title: f.ProductName, Width: 960, Height: 640, URL: "about:blank"}
	}
	w := f.App.Windows[0]
	if w.Label == "" {
		w.Label = "main"
	}
	return w
}
