package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: nekokit <files.nk> [--build] [--output name] [--update] [--verbose]")
		return
	}

	var files []string
	var isBuild bool
	var verbose bool
	var customOutput string

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch {
		case arg == "--update":
			updateSystem()
			return
		case arg == "--verbose" || arg == "-v":
			verbose = true
		case arg == "--build":
			isBuild = true
		case arg == "--output" && i+1 < len(os.Args):
			customOutput = os.Args[i+1]
			i++
		case strings.HasSuffix(arg, ".nk"):
			files = append(files, arg)
		}
	}

	if verbose {
		fmt.Printf("Found %d file(s) to process\n", len(files))
		if isBuild {
			fmt.Println("Mode: BUILD")
		} else {
			fmt.Println("Mode: RUN")
		}
	}

	for _, file := range files {
		if verbose {
			fmt.Printf("Processing file: %s\n", file)
		}

		source, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("error %s: %v\n", file, err)
			continue
		}

		if verbose {
			fmt.Printf("File read successfully (%d bytes)\n", len(source))
			fmt.Println("Starting transpilation...")
		}

		goCode := TranspileWithVerbose(string(source), verbose)

		if verbose {
			fmt.Printf("Transpilation complete! Generated %d lines of Go code\n", len(strings.Split(goCode, "\n")))
			
			// Append to log file
			logFile := strings.TrimSuffix(file, ".nk") + ".log"
			logContent := "\n=== GENERATED GO CODE ===\n" + goCode + "\n=== END GO CODE ===\n"
			existingLog, _ := os.ReadFile(logFile)
			os.WriteFile(logFile, append(existingLog, []byte(logContent)...), 0644)
			fmt.Printf("Log appended to: %s\n", logFile)
			fmt.Println("\n=== GENERATED GO CODE ===")
			fmt.Println(goCode)
			fmt.Println("=== END GO CODE ===\n")
		}

		if isBuild {
			outName := customOutput
			if outName == "" || len(files) > 1 {
				outName = strings.TrimSuffix(file, ".nk")
			}
			if verbose {
				fmt.Printf("Building executable: %s\n", outName)
			}
			BuildGo(goCode, outName, verbose)
		} else {
			if verbose {
				fmt.Println("Running generated code...")
				fmt.Println("=" + strings.Repeat("=", 49))
			}
			RunGo(goCode)
			if verbose {
				fmt.Println("=" + strings.Repeat("=", 49))
				fmt.Println("Execution complete!")
			}
		}
	}
}

func updateSystem() {
	cmd := exec.Command("sh", "-c", "curl -fsSL https://raw.githubusercontent.com/NekoKatoriChan/NekoKit/main/install.sh | sh")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	cmd.Run()
}
