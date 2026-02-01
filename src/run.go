package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func RunGo(code string) {
	tmpFile := "nk_tmp.go"
	os.WriteFile(tmpFile, []byte(code), 0644)
	cmd := exec.Command("go", "run", tmpFile)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	cmd.Run()
	os.Remove(tmpFile)
}

func BuildGo(code string, outputName string) {
	tmpDir, _ := os.MkdirTemp("", "nk_build_*")
	defer os.RemoveAll(tmpDir)
	tmpFile := filepath.Join(tmpDir, "main.go")
	os.WriteFile(tmpFile, []byte(code), 0644)
	modCmd := exec.Command("go", "mod", "init", "nk_gen")
	modCmd.Dir = tmpDir
	modCmd.Run()
	absOutput, _ := filepath.Abs(outputName)
	buildCmd := exec.Command("go", "build", "-o", absOutput, "main.go")
	buildCmd.Dir = tmpDir
	buildCmd.Stdout, buildCmd.Stderr = os.Stdout, os.Stderr
	buildCmd.Run()
}
