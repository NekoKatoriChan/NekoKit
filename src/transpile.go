package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// this is the transpiler for making silly game scripts into proper go code!
// completely groomed for Termux, modern Go, and lots of catnip! prrr~

// meow meow helper functions start here!

// STRING INTERPOLATION:
// turns "$variable" thingies into fmt.Sprintf stuff~
// pretty clever, nya? turns "hello $name" into proper go format strings~
func transpileStringInterpolation(arg string) string {
	if strings.HasPrefix(arg, "\"") && strings.HasSuffix(arg, "\"") {
		content := arg[1 : len(arg)-1]

		var vars []string
		var formatStr strings.Builder
		formatStr.WriteString("\"")
		hasVars := false

		i := 0
		for i < len(content) {
			if content[i] == '$' && i+1 < len(content) {
				j := i + 1
				for j < len(content) && isAlphaNumericOrUnderscore(content[j]) {
					j++
				}
				varName := content[i+1 : j]
				vars = append(vars, varName)
				formatStr.WriteString("%v")
				hasVars = true
				i = j
			} else {
				if content[i] == '"' {
					formatStr.WriteString("\\\"")
				} else {
					formatStr.WriteRune(rune(content[i]))
				}
				i++
			}
		}
		formatStr.WriteString("\"")

		// if we have variables, use fmt.Sprintf; otherwise just return the escaped string
		if hasVars {
			result := "fmt.Sprintf(" + formatStr.String()
			for _, v := range vars {
				result += ", " + v
			}
			result += ")"
			return result
		}
		return formatStr.String()
	}
	return arg
}

// removes the $ from variable names, like when we see "$s" it becomes just "s"
// purrfect for cleaning up our variable syntax~ works as a grooming tool, meow!
func stripDollar(s string) string {
	if strings.HasPrefix(s, "$") {
		return s[1:]
	}
	return s
}

// checks if a byte is a letter, number, or underscore~
// basically asking "is this a valid part of a variable name?" nya!
func isAlphaNumericOrUnderscore(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

// this expands the "~/" path into the actual home directory path
// heavily optimized for Termux! ($HOME points to /data/data/com.termux/files/home)
func expandPath(pathStr string) string {
	if strings.HasPrefix(pathStr, "~/") {
		return "os.ExpandEnv(\"$HOME" + pathStr[1:] + "\")"
	}
	return "\"" + pathStr + "\""
}

// grabs just the filename from a path and removes the extension
// much cleaner now using Go's built-in filepath tools! no more messy string splitting!
func extractFilename(pathStr string) string {
	base := filepath.Base(pathStr)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

// this is where all the magic happens, turning cute game script into proper go code
// it's like translating cat into human language~ very sophisticated, nya!
func Transpile(input string) string {
	return transpileInternal(input, false)
}

// Transpile with verbose output tracking~
// meow! shows you all the transformations happening
func TranspileWithVerbose(input string, verbose bool) string {
	return transpileInternal(input, verbose)
}

func transpileInternal(input string, verbose bool) string {
	lines := strings.Split(input, "\n")
	var out []string

	// tracking two separate output bowls: mainOut for main(), blockOut for named funcs~
	// do NOT mix the kibbles!
	var mainOut []string
	var blockOut []string
	inNamedBlock := false

	mainIndent := 1
	blockIndent := 1
	savedMainIndent := 1

	emitLine := func(s string) {
		if inNamedBlock {
			blockOut = append(blockOut, s)
		} else {
			mainOut = append(mainOut, s)
		}
	}

	getIndent := func() string {
		if inNamedBlock {
			return strings.Repeat("    ", blockIndent)
		}
		return strings.Repeat("    ", mainIndent)
	}

	bumpIndent := func(n int) {
		if inNamedBlock {
			blockIndent += n
		} else {
			mainIndent += n
		}
	}

	hasRead := false
	hasRandom := false
	hasFileOps := false
	hasRun := false
	hasFmt := false
	hasSusu := false
	hasStrings := false

	mainDeclared := make(map[string]bool)
	blockDeclared := make(map[string]bool)

	getDeclared := func() map[string]bool {
		if inNamedBlock {
			return blockDeclared
		}
		return mainDeclared
	}

	runCmdCount := 0

	gameBlocks := make(map[string]bool)
	callonceBlocks := make(map[string]bool)

	stats := map[string]int{
		"dialogs":      0,
		"menus":        0,
		"prompts":      0,
		"gameloops":    0,
		"variables":    0,
		"conditionals": 0,
		"file_ops":     0,
		"damage":       0,
		"heal":         0,
		"score":        0,
		"level":        0,
		"run":          0,
	}

	// first pass: hunt for treats! find all named game blocks
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if strings.HasSuffix(line, " start") {
			blockName := strings.TrimSuffix(line, " start")
			if blockName != "gameloop" && !strings.Contains(blockName, " ") {
				gameBlocks[blockName] = true
			}
		}
	}

	// second pass: find callonce blocks~
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "callonce ") {
			blockName := strings.TrimPrefix(line, "callonce ")
			callonceBlocks[blockName] = true
		}
	}

	hasCallonce := len(callonceBlocks) > 0

	// third pass: scan for imports! we only take the toys we need~
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "read -p ") || strings.HasPrefix(line, "read ") {
			hasRead = true
			hasFmt = true
			hasStrings = true
		}
		if strings.HasPrefix(line, "random ") {
			hasRandom = true
		}
		if strings.HasPrefix(line, "load ") || strings.HasPrefix(line, "save ") {
			hasFileOps = true
		}
		if strings.HasPrefix(line, "run ") {
			hasRun = true
		}
		if strings.HasPrefix(line, "susu") {
			hasSusu = true
		}
		if strings.HasPrefix(line, "clear") || strings.HasPrefix(line, "border ") ||
			strings.HasPrefix(line, "dialog ") || strings.HasPrefix(line, "menu ") ||
			strings.HasPrefix(line, "prompt ") || strings.HasPrefix(line, "stat ") ||
			strings.HasPrefix(line, "write ") || strings.HasPrefix(line, "writeln ") {
			hasFmt = true
		}
		if strings.HasPrefix(line, "border ") {
			hasStrings = true
		}
		if strings.HasPrefix(line, "peek ") {
			hasFileOps = true
		}
	}

	out = append(out,
		"package main",
		"import (",
	)

	if hasFmt {
		out = append(out, "    \"fmt\"")
	}
	if hasRead || hasFileOps || hasSusu {
		out = append(out, "    \"os\"")
	}
	if hasRead {
		out = append(out, "    \"bufio\"")
	}
	if hasStrings {
		out = append(out, "    \"strings\"")
	}
	if hasRandom {
		// Modern Go (1.20+) handles seeding automatically! Purrfectly clean~
		out = append(out, "    \"math/rand\"")
	}
	if hasRun {
		out = append(out, "    \"os/exec\"")
	}
	out = append(out, ")", "")

	if hasRead {
		out = append(out, "// reader is package-level so all blocks can read from the terminal bowl~")
		out = append(out, "var reader = bufio.NewReader(os.Stdin)")
		out = append(out, "")
	}

	if hasCallonce {
		out = append(out, "// tracks which blocks have already been called, meow!")
		out = append(out, "var callonceTracker = make(map[string]bool)")
		out = append(out, "")
	}

	out = append(out, "func main() {")

	if verbose {
		fmt.Println("Import analysis:")
		if hasFmt { fmt.Printf("         - fmt: yes\n") }
		if hasRead || hasFileOps || hasSusu { fmt.Printf("         - os: yes\n") }
		if hasRead { fmt.Printf("         - bufio, strings: yes\n") }
		if hasRandom { fmt.Printf("         - math/rand: yes (no manual seed needed in modern Go!)\n") }
		if hasRun { fmt.Printf("         - os/exec: yes\n") }
		if hasCallonce { fmt.Printf("         - callonce map: yes\n") }
	}

	// here comes the fun part! parsing all the silly game commands~
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		indent := getIndent()
		declared := getDeclared()

		switch {
		case line == "gameloop start":
			stats["gameloops"]++
			emitLine(indent + "for {")
			bumpIndent(1)

		case line == "gameloop end":
			bumpIndent(-1)
			emitLine(getIndent() + "}")

		case strings.HasSuffix(line, " start"):
			blockName := strings.TrimSuffix(line, " start")
			if gameBlocks[blockName] {
				savedMainIndent = mainIndent
				inNamedBlock = true
				blockIndent = 1
				blockDeclared = make(map[string]bool)
				blockOut = append(blockOut, "")
				blockOut = append(blockOut, "// game block '"+blockName+"'~ enter at your own risk! nya~")
				blockOut = append(blockOut, "func "+blockName+"() {")
			}

		case strings.HasSuffix(line, " end"):
			blockName := strings.TrimSuffix(line, " end")
			if gameBlocks[blockName] {
				blockIndent--
				blockOut = append(blockOut, "}")
				blockOut = append(blockOut, "")
				inNamedBlock = false
				mainIndent = savedMainIndent
				blockDeclared = make(map[string]bool)
			}

		case strings.HasPrefix(line, "call "):
			blockName := strings.TrimPrefix(line, "call ")
			if gameBlocks[blockName] {
				emitLine(indent + blockName + "()")
			}

		case strings.HasPrefix(line, "callonce "):
			blockName := strings.TrimPrefix(line, "callonce ")
			if gameBlocks[blockName] && callonceBlocks[blockName] {
				emitLine(indent + "if !callonceTracker[\"" + blockName + "\"] {")
				emitLine(indent + "    " + blockName + "()")
				emitLine(indent + "    callonceTracker[\"" + blockName + "\"] = true")
				emitLine(indent + "} else {")
				emitLine(indent + "    // nya! already called this block once, no second helpings~")
				emitLine(indent + "}")
			}

		case line == "clear":
			// A purrfect ANSI sequence for Termux screens!
			emitLine(indent + `fmt.Print("\033[2J\033[H")`)

		case strings.HasPrefix(line, "border "):
			boxType := strings.TrimPrefix(line, "border ")
			switch boxType {
			case "top", "bot":
				emitLine(indent + `fmt.Println(strings.Repeat("═", 50))`)
			case "mid":
				emitLine(indent + `fmt.Println(strings.Repeat("─", 50))`)
			}

		case strings.HasPrefix(line, "dialog "):
			stats["dialogs"]++
			dialogText := strings.TrimPrefix(line, "dialog ")
			goCode := transpileStringInterpolation(dialogText)
			emitLine(indent + `fmt.Println("════════════════════════════════════════════════════")`)
			emitLine(indent + `fmt.Println(" " + ` + goCode + ")")
			emitLine(indent + `fmt.Println("════════════════════════════════════════════════════")`)

		case strings.HasPrefix(line, "menu "):
			stats["menus"]++
			menuContent := strings.TrimPrefix(line, "menu ")
			options := strings.Split(menuContent, ",")
			emitLine(indent + `fmt.Println("[MENU]")`)
			for i, opt := range options {
				emitLine(indent + `fmt.Println("  (` + fmt.Sprintf("%d", i+1) + `) ` + strings.TrimSpace(opt) + `")`)
			}

		case strings.HasPrefix(line, "prompt "):
			stats["prompts"]++
			promptText := strings.TrimPrefix(line, "prompt ")
			goCode := transpileStringInterpolation(promptText)
			emitLine(indent + `fmt.Print(" > " + ` + goCode + ")")

		case strings.HasPrefix(line, "stat "):
			rest := strings.TrimPrefix(line, "stat ")
			var statNameRaw, statValueRaw string

			if strings.HasPrefix(rest, "\"") {
				closeIdx := strings.Index(rest[1:], "\"")
				if closeIdx > 0 {
					statNameRaw = rest[:closeIdx+2]
					statValueRaw = strings.TrimSpace(rest[closeIdx+2:])
				}
			} else {
				spaceIdx := strings.Index(rest, " ")
				if spaceIdx > 0 {
					statNameRaw = rest[:spaceIdx]
					statValueRaw = strings.TrimSpace(rest[spaceIdx+1:])
				}
			}

			if statNameRaw != "" && statValueRaw != "" {
				statValue := stripDollar(statValueRaw)
				var formatStr string
				if strings.HasPrefix(statNameRaw, "\"") && strings.HasSuffix(statNameRaw, "\"") {
					innerName := statNameRaw[1 : len(statNameRaw)-1]
					formatStr = "\"" + innerName + ": %v\\n\""
				} else {
					formatStr = "\"" + statNameRaw + ": %v\\n\""
				}
				emitLine(indent + "fmt.Printf(" + formatStr + ", " + statValue + ")")
			}

		// file ops! heavily improved for modern go (os.ReadFile / os.WriteFile)
		case strings.HasPrefix(line, "load "):
			stats["file_ops"]++
			rest := strings.TrimPrefix(line, "load ")
			parts := strings.Fields(rest)
			var varName, filePath string
			
			if len(parts) == 1 {
				filePath = parts[0]
				varName = extractFilename(filePath)
			} else if len(parts) >= 2 {
				varName = parts[0]
				filePath = strings.Join(parts[1:], " ")
			} else {
				continue
			}
			
			pathExpr := expandPath(filePath)
			dataVar := varName + "Data"
			
			if !declared[dataVar] {
				emitLine(indent + dataVar + ", _ := os.ReadFile(" + pathExpr + ")")
				declared[dataVar] = true
			} else {
				emitLine(indent + dataVar + ", _ = os.ReadFile(" + pathExpr + ")")
			}
			
			if !declared[varName] {
				emitLine(indent + varName + " := string(" + dataVar + ")")
				declared[varName] = true
			} else {
				emitLine(indent + varName + " = string(" + dataVar + ")")
			}

		case strings.HasPrefix(line, "save "):
			stats["file_ops"]++
			parts := strings.SplitN(strings.TrimPrefix(line, "save "), " ", 2)
			if len(parts) == 2 {
				varName := parts[0]
				filePath := parts[1]
				pathExpr := expandPath(filePath)
				emitLine(indent + "os.WriteFile(" + pathExpr + ", []byte(" + varName + "), 0644)")
			}

		case strings.HasPrefix(line, "damage "):
			stats["damage"]++
			parts := strings.Fields(strings.TrimPrefix(line, "damage "))
			if len(parts) >= 2 {
				emitLine(indent + parts[0] + " -= " + stripDollar(parts[1]))
			}

		case strings.HasPrefix(line, "heal "):
			stats["heal"]++
			parts := strings.Fields(strings.TrimPrefix(line, "heal "))
			if len(parts) >= 2 {
				emitLine(indent + parts[0] + " += " + stripDollar(parts[1]))
			}

		case strings.HasPrefix(line, "score "), strings.HasPrefix(line, "level "):
			// Combined these two since the behavior is essentially identical! Meow~
			isScore := strings.HasPrefix(line, "score ")
			if isScore { stats["score"]++ } else { stats["level"]++ }
			
			prefix := "level "
			if isScore { prefix = "score " }
			
			parts := strings.Fields(strings.TrimPrefix(line, prefix))
			if len(parts) >= 2 {
				varName := parts[0]
				amount := stripDollar(parts[1])
				if !declared[varName] {
					emitLine(indent + varName + " := " + amount)
					declared[varName] = true
				} else {
					operator := " = "
					if isScore { operator = " += " } // score adds, level sets~
					emitLine(indent + varName + operator + amount)
				}
			}

		case strings.HasPrefix(line, "reset "):
			emitLine(indent + strings.TrimPrefix(line, "reset ") + " = 0")

		case strings.HasPrefix(line, "random "):
			parts := strings.Fields(strings.TrimPrefix(line, "random "))
			if len(parts) >= 2 {
				varName := parts[0]
				maxVal := parts[1]
				if !declared[varName] {
					emitLine(indent + varName + " := rand.Intn(" + maxVal + ")")
					declared[varName] = true
				} else {
					emitLine(indent + varName + " = rand.Intn(" + maxVal + ")")
				}
			}

		case strings.HasPrefix(line, "if "):
			stats["conditionals"]++
			condition := strings.TrimPrefix(line, "if ")
			emitLine(indent + "if " + condition + " {")
			bumpIndent(1)

		case strings.HasPrefix(line, "peek "):
			stats["conditionals"]++
			rest := strings.TrimPrefix(line, "peek ")
			
			var filePath string
			if strings.HasPrefix(rest, "\"") {
				endQuote := strings.Index(rest[1:], "\"")
				if endQuote != -1 {
					filePath = rest[1 : endQuote+1]
				}
			} else {
				braceIdx := strings.Index(rest, "{")
				if braceIdx > 0 {
					filePath = strings.TrimSpace(rest[:braceIdx])
				}
			}
			
			if filePath != "" {
				pathExpr := expandPath(filePath)
				emitLine(indent + "if _, err := os.Stat(" + pathExpr + "); err == nil {")
				bumpIndent(1)
			}

		case line == "} else {":
			bumpIndent(-1)
			emitLine(getIndent() + "} else {")
			bumpIndent(1)

		case line == "}":
			bumpIndent(-1)
			emitLine(getIndent() + "}")

		case strings.HasPrefix(line, "write "):
			goCode := transpileStringInterpolation(strings.TrimPrefix(line, "write "))
			emitLine(indent + "fmt.Print(" + goCode + ")")

		case strings.HasPrefix(line, "writeln "):
			goCode := transpileStringInterpolation(strings.TrimPrefix(line, "writeln "))
			emitLine(indent + "fmt.Println(" + goCode + ")")

		case strings.HasPrefix(line, "read -p "):
			rest := strings.TrimPrefix(line, "read -p ")
			promptEndIdx := strings.LastIndex(rest, "\"")
			if promptEndIdx > 0 && strings.HasPrefix(rest, "\"") {
				prompt := rest[:promptEndIdx+1]
				varName := strings.TrimSpace(rest[promptEndIdx+1:])
				promptCode := transpileStringInterpolation(prompt)

				if !declared[varName] {
					emitLine(indent + varName + ` := ""`)
					declared[varName] = true
				}
				emitLine(indent + "fmt.Print(" + promptCode + ")")
				emitLine(indent + varName + `, _ = reader.ReadString('\n')`)
				emitLine(indent + varName + " = strings.TrimSpace(" + varName + ")")
				stats["variables"]++
			}

		case strings.HasPrefix(line, "read "):
			varName := strings.TrimPrefix(line, "read ")
			if !declared[varName] {
				emitLine(indent + varName + ` := ""`)
				declared[varName] = true
			}
			emitLine(indent + varName + `, _ = reader.ReadString('\n')`)
			emitLine(indent + varName + " = strings.TrimSpace(" + varName + ")")
			stats["variables"]++

		case strings.HasPrefix(line, "give "):
			rest := strings.TrimPrefix(line, "give ")
			parts := strings.SplitN(rest, "=", 2)
			if len(parts) == 2 {
				varName := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				goValue := transpileStringInterpolation(value)
				if !strings.HasPrefix(value, "\"") {
					goValue = stripDollar(value)
				}

				if !declared[varName] {
					emitLine(indent + varName + " := " + goValue)
					emitLine(indent + "_ = " + varName) // keeps Go compiler happy like a well-fed kitty!
					declared[varName] = true
					stats["variables"]++
				} else {
					emitLine(indent + varName + " = " + goValue)
				}
			}

		case strings.HasPrefix(line, "susu"):
			emitLine(indent + "os.Exit(0)")

		case strings.HasPrefix(line, "run "):
			stats["run"]++
			cmd := strings.TrimPrefix(line, "run ")
			escapedCmd := strings.ReplaceAll(cmd, "\"", "\\\"")
			cmdVar := fmt.Sprintf("_cmd%d", runCmdCount)
			runCmdCount++
			// Using standard "sh" which Termux gracefully intercepts natively!
			emitLine(indent + cmdVar + ` := exec.Command("sh", "-c", "` + escapedCmd + `")`)
			emitLine(indent + cmdVar + ".Stdout = os.Stdout")
			emitLine(indent + cmdVar + ".Stderr = os.Stderr")
			emitLine(indent + cmdVar + ".Run()")

		case strings.HasPrefix(line, "create "):
			stats["file_ops"]++
			filePath := strings.TrimPrefix(line, "create ")
			pathExpr := expandPath(filePath)
			emitLine(indent + "os.WriteFile(" + pathExpr + ", []byte(\"\"), 0644)")
		}
	}

	out = append(out, mainOut...)
	out = append(out, "}") // properly close main() before the named blok functions!
	out = append(out, blockOut...)

	if verbose {
		fmt.Println("\n=== TRANSPILATION SUMMARY ===")
		fmt.Printf("Variables declared: %d\n", stats["variables"])
		fmt.Printf("Game loops: %d\n", stats["gameloops"])
		fmt.Printf("File operations: %d\n", stats["file_ops"])
		fmt.Println("=======================")
	}

	return strings.Join(out, "\n")
}
