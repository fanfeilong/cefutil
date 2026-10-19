// Command xgoef is a Tauri-like CLI placeholder (create / dev / build).
//
// Real builds still depend on cef-download skills and llgo/llcppg.
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "help", "-h", "--help":
		usage()
	case "info":
		fmt.Println("xgoef: XGo + LLGo + CEF client framework (scaffold)")
		fmt.Println("see src/xgoef/RESEARCH.md and DESIGN.md")
	case "create", "dev", "build":
		fmt.Fprintf(os.Stderr, "xgoef %s: not implemented yet (bindings + runtime first)\n", os.Args[1])
		os.Exit(1)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `xgoef — scaffold CLI for the XGo CEF framework

Usage:
  xgoef info
  xgoef create <name>   (planned)
  xgoef dev             (planned)
  xgoef build           (planned)

Prerequisites for real builds: CEF_ROOT, llgo, llcppg.
`)
}
