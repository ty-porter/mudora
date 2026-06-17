// Command build builds the mudora-web WASM binary and stages wasm_exec.js
// alongside it. Run with: go run ./cmd/mudora-web/build [-version X.Y.Z]
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	version := flag.String("version", "dev", "version string baked into the build via -ldflags")
	flag.Parse()

	_, thisFile, _, _ := runtime.Caller(0)
	root := filepath.Dir(filepath.Dir(thisFile))
	webDir := filepath.Join(root, "web")

	ldflags := "-X github.com/alttpr-mudora/mudora/internal.Version=" + *version
	if err := run("go", "build", "-ldflags", ldflags, "-o", filepath.Join(webDir, "main.wasm"), root); err != nil {
		fail("build failed: %v", err)
	}

	goroot, err := output("go", "env", "GOROOT")
	if err != nil {
		fail("go env GOROOT failed: %v", err)
	}

	src := filepath.Join(strings.TrimSpace(goroot), "lib", "wasm", "wasm_exec.js")
	data, err := os.ReadFile(src)
	if err != nil {
		fail("reading wasm_exec.js failed: %v", err)
	}
	dst := filepath.Join(webDir, "wasm_exec.js")
	_ = os.Chmod(dst, 0o644)
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		fail("writing wasm_exec.js failed: %v", err)
	}

	fmt.Println("Built", filepath.Join(webDir, "main.wasm"))
	fmt.Println("Serve", webDir, "with any static file server, e.g.:")
	fmt.Println("  python -m http.server 8080   (run from", webDir+")")
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func output(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	return string(out), err
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
