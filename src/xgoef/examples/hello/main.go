// Hello is the minimal xgoef application shape.
//
// With CEF bindings absent, Run returns cef.ErrBindingsMissing.
// After llcppg generation and a real cef.Host, this becomes the
// entry equivalent of src/cefdemo, written in Go/XGo style.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fanfeilong/cefutil/src/xgoef/conf"
	"github.com/fanfeilong/cefutil/src/xgoef/ipc"
	"github.com/fanfeilong/cefutil/src/xgoef/runtime"
)

func main() {
	cfgPath := "xgoef.conf.json"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	} else if _, err := os.Stat(cfgPath); err != nil {
		// Allow `go run ./examples/hello` from the module root.
		cfgPath = filepath.Join("examples", "hello", "xgoef.conf.json")
	}
	abs, err := filepath.Abs(cfgPath)
	if err != nil {
		fatal(err)
	}
	cfg, err := conf.Load(abs)
	if err != nil {
		fatal(err)
	}

	app := runtime.New(cfg, nil)
	app.Handle("greet", func(ctx ipc.Context, argsJSON []byte) (any, error) {
		var args struct {
			Name string `json:"name"`
		}
		if len(argsJSON) > 0 {
			_ = json.Unmarshal(argsJSON, &args)
		}
		if args.Name == "" {
			args.Name = "world"
		}
		return fmt.Sprintf("hello, %s (from xgoef)", args.Name), nil
	})

	if err := app.Run(); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "xgoef-hello: %v\n", err)
	os.Exit(1)
}
