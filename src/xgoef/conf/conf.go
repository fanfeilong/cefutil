package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

// File is the on-disk xgoef.conf.json shape (Tauri-like).
type File struct {
	ProductName  string         `json:"productName"`
	Version      string         `json:"version"`
	Identifier   string         `json:"identifier"`
	Build        Build          `json:"build"`
	App          AppSection     `json:"app"`
	CEF          CEFSection     `json:"cef"`
	Plugins      []string       `json:"plugins"`
	Capabilities map[string]any `json:"capabilities"`
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

// Load reads and decodes a config file.
func Load(path string) (File, error) {
	var f File
	data, err := os.ReadFile(path)
	if err != nil {
		return f, err
	}
	if err := json.Unmarshal(data, &f); err != nil {
		return f, fmt.Errorf("conf: %w", err)
	}
	if f.ProductName == "" {
		return f, fmt.Errorf("conf: productName is required")
	}
	if f.CEF.RootEnv == "" {
		f.CEF.RootEnv = "CEF_ROOT"
	}
	if f.Build.FrontendDist == "" {
		f.Build.FrontendDist = "web"
	}
	return f, nil
}

// FirstWindow returns the first window or a default.
func (f File) FirstWindow() Window {
	if len(f.App.Windows) == 0 {
		return Window{Title: f.ProductName, Width: 960, Height: 640, URL: "about:blank"}
	}
	return f.App.Windows[0]
}
