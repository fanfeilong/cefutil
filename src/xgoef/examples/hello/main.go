// Hello is the minimal xgoef application.
//
// Build with scripts/build-xgoef-hello.sh (needs CEF_ROOT + LLGo; no cgo).
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fanfeilong/cefutil/src/xgoef/cef"
	"github.com/fanfeilong/cefutil/src/xgoef/conf"
	"github.com/fanfeilong/cefutil/src/xgoef/ipc"
	xruntime "github.com/fanfeilong/cefutil/src/xgoef/runtime"
)

func main() {
	// CEF helper processes (GPU, renderer, …) re-enter this executable.
	// Hand them off before the Go runtime does more work.
	runtime.LockOSThread()
	host := cef.NewNativeHost()
	if host == nil {
		fatal(cef.ErrBindingsMissing)
	}
	host.SetArgs(os.Args)
	if err := host.InitAPIHash(); err != nil {
		fatal(err)
	}
	if code, helper := host.ExecuteProcess(); helper {
		os.Exit(code)
	}

	cfgPath := "xgoef.conf.json"
	if len(os.Args) > 1 && !isCEFSwitch(os.Args[1]) {
		cfgPath = os.Args[1]
	} else if _, err := os.Stat(cfgPath); err != nil {
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

	base := filepath.Dir(abs)
	if cfg.CEF.CacheDir != "" && !filepath.IsAbs(cfg.CEF.CacheDir) {
		cfg.CEF.CacheDir = filepath.Join(base, cfg.CEF.CacheDir)
	}
	if len(cfg.App.Windows) > 0 {
		w := &cfg.App.Windows[0]
		if w.URL == "" {
			w.URL = "file://" + filepath.Join(base, "web", "index.html")
		} else if len(w.URL) >= 7 && w.URL[:7] == "file://" {
			p := w.URL[7:]
			if !filepath.IsAbs(p) {
				w.URL = "file://" + filepath.Join(base, p)
			}
		}
	}

	app := xruntime.New(cfg, host)
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
		if ee, ok := err.(*xruntime.ExitError); ok {
			os.Exit(ee.Code)
		}
		fatal(err)
	}
}

func isCEFSwitch(s string) bool {
	return len(s) > 1 && s[0] == '-'
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "xgoef-hello: %v\n", err)
	os.Exit(1)
}
