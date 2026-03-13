
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func RunGo(code string) {
	// Create temp file in system temp directory with unique name
	tmpFile, err := os.CreateTemp("", "nk_*.go")
	if err != nil {
		fmt.Printf("error creating temp file: %v\n", err)
		return
	}
	tmpFileName := tmpFile.Name()
	defer os.Remove(tmpFileName)
	
	// Write code to temp file with error checking
	_, err = tmpFile.WriteString(code)
	if err != nil {
		fmt.Printf("error writing to temp file: %v\n", err)
		tmpFile.Close()
		return
	}
	tmpFile.Close()
	
	// Run the generated Go code
	cmd := exec.Command("go", "run", tmpFileName)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("error running generated code: %v\n", err)
		return
	}
}

func BuildGo(code string, outputName string, verbose bool) {
	if verbose {
		fmt.Println("Creating temporary build directory...")
	}

	tmpDir, err := os.MkdirTemp("", "nk_build_*")
	if err != nil {
		fmt.Printf("error creating build directory: %v\n", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	if verbose {
		fmt.Printf("Build directory: %s\n", tmpDir)
	}

	tmpFile := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(tmpFile, []byte(code), 0644)
	if err != nil {
		fmt.Printf("error writing Go code: %v\n", err)
		return
	}

	if verbose {
		fmt.Printf("Writing Go code to: %s\n", tmpFile)
	}

	modCmd := exec.Command("go", "mod", "init", "nk_gen")
	modCmd.Dir = tmpDir

	if verbose {
		fmt.Println("Initializing Go module...")
	}

	if err := modCmd.Run(); err != nil {
		fmt.Printf("error initializing module: %v\n", err)
		return
	}

	absOutput, err := filepath.Abs(outputName)
	if err != nil {
		fmt.Printf("error resolving output path: %v\n", err)
		return
	}

	if verbose {
		fmt.Printf("Building to: %s\n", absOutput)
		fmt.Println("Running 'go build'...")
	}

	buildCmd := exec.Command("go", "build", "-o", absOutput, "main.go")
	buildCmd.Dir = tmpDir
	buildCmd.Stdout, buildCmd.Stderr = os.Stdout, os.Stderr
	
	if err := buildCmd.Run(); err != nil {
		fmt.Printf("error building: %v\n", err)
		return
	}

	if verbose {
		fmt.Printf("Build complete! Executable saved to: %s\n", absOutput)
	}
}
