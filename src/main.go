package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: nekokit <files.nk> [--build] [--output name] [--update]")
		return
	}

	var files []string
	var isBuild bool
	var customOutput string

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch {
		case arg == "--update":
			updateSystem()
			return
		case arg == "--build":
			isBuild = true
		case arg == "--output" && i+1 < len(os.Args):
			customOutput = os.Args[i+1]
			i++
		case strings.HasSuffix(arg, ".nk"):
			files = append(files, arg)
		}
	}

	for _, file := range files {
		source, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("error %s: %v\n", file, err)
			continue
		}

		goCode := Transpile(string(source))

		if isBuild {
			outName := customOutput
			if outName == "" || len(files) > 1 {
				outName = strings.TrimSuffix(file, ".nk")
			}
			BuildGo(goCode, outName)
		} else {
			RunGo(goCode)
		}
	}
}

func updateSystem() {
	cmd := exec.Command("sh", "-c", "curl -fsSL https://raw.githubusercontent.com/NekoKatoriChan/NekoKit/main/install.sh | sh")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	cmd.Run()
}
