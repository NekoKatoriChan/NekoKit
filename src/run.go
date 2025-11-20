package main

import (
	"os"
	"os/exec"
)

func RunGo(code string) {
	tmpFile := "nk_tmp.go"

	err := os.WriteFile(tmpFile, []byte(code), 0644)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "run", tmpFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	os.Remove(tmpFile)
}
