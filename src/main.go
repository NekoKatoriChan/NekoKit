package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {

	if len(os.Args) >= 2 && os.Args[1] == "update" {
		cmd := exec.Command(
			"sh", "-c",
			"curl -fsSL https://raw.githubusercontent.com/NekoKatoriChan/NekoKit/refs/heads/main/install.sh | sh",
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if len(os.Args) < 2 {
		fmt.Println("usage: nekokit <file.nk>")
		return
	}

	filename := os.Args[1]

	source, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	goCode := Transpile(string(source))
	RunGo(goCode)
}
