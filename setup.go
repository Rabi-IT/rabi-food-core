//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
)

type dep struct {
	name string
	pkg  string
}

var deps = []dep{
	{"Task runner", "github.com/go-task/task/v3/cmd/task@latest"},
	{"golangci-lint", "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"},
	{"gotestsum", "gotest.tools/gotestsum@latest"},
	{"mockery", "github.com/vektra/mockery/v2@latest"},
	{"swaggo", "github.com/swaggo/swag/cmd/swag@latest"},
}

func main() {
	fmt.Println("🚀 Installing development dependencies...")

	for _, d := range deps {
		fmt.Printf("📦 Installing %s...\n", d.name)
		run("go", "install", d.pkg)
	}

	fmt.Println("📁 Installing Go modules...")
	runDir("api", "go", "mod", "download")

	fmt.Println()
	fmt.Println("✅ All dependencies installed successfully!")
	fmt.Println()
	fmt.Println("🎯 Next steps:")
	fmt.Println("   task dev    # Start development environment")
	fmt.Println("   task test   # Run tests")
	fmt.Println("   task lint   # Run linter")
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to run %q: %v\n", name, err)
		os.Exit(1)
	}
}

func runDir(dir, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to run %q in %q: %v\n", name, dir, err)
		os.Exit(1)
	}
}
