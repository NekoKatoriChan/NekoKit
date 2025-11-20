package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: nk <file.nk>")
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
