
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func RunGo(code string) {
	tmpFile := "nk_tmp.go"
	os.WriteFile(tmpFile, []byte(code), 0644)
	cmd := exec.Command("go", "run", tmpFile)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.Run()
	os.Remove(tmpFile)
}

func BuildGo(code string, outputName string, verbose bool) {
	if verbose {
		fmt.Println("Creating temporary build directory...")
	}

	tmpDir, _ := os.MkdirTemp("", "nk_build_*")
	defer os.RemoveAll(tmpDir)

	if verbose {
		fmt.Printf("Build directory: %s\n", tmpDir)
	}

	tmpFile := filepath.Join(tmpDir, "main.go")
	os.WriteFile(tmpFile, []byte(code), 0644)

	if verbose {
		fmt.Printf("Writing Go code to: %s\n", tmpFile)
	}

	modCmd := exec.Command("go", "mod", "init", "nk_gen")
	modCmd.Dir = tmpDir

	if verbose {
		fmt.Println("Initializing Go module...")
	}

	modCmd.Run()

	absOutput, _ := filepath.Abs(outputName)

	if verbose {
		fmt.Printf("Building to: %s\n", absOutput)
		fmt.Println("Running 'go build'...")
	}

	buildCmd := exec.Command("go", "build", "-o", absOutput, "main.go")
	buildCmd.Dir = tmpDir
	buildCmd.Stdout, buildCmd.Stderr = os.Stdout, os.Stderr
	buildCmd.Run()

	if verbose {
		fmt.Printf("Build complete! Executable saved to: %s\n", absOutput)
	}
}
