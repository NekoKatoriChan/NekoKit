package main

import (
	"fmt"
	"strings"
)

// this is the transpiler for making silly game scripts into proper go code!
// made with love and fish treats 

// meow meow helper functions start here!


// STRING INTERPOLATION:
// this one turns "$variable" thingies into fmt.Sprintf stuff
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
				for j < len(content) && (isAlphaNumericOrUnderscore(content[j])) {
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

		// if we have variables, use fmt.Sprintf; otherwise just escape strings
		if hasVars {
			result := "fmt.Sprintf(" + formatStr.String()
			for _, v := range vars {
				result += ", " + v
			}
			result += ")"
			return result
		} else {
			return formatStr.String()
		}
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
// cats love having their home paths expanded, it's cozy~
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		return "os.ExpandEnv(\"$HOME" + path[1:] + "\")"
	}
	return "\"" + path + "\""
}

// grabs just the filename from a path and removes the extension
// like finding the treat inside the wrapper!
func extractFilename(path string) string {
	parts := strings.Split(path, "/")
	filename := parts[len(parts)-1]
	dotIdx := strings.LastIndex(filename, ".")
	if dotIdx > 0 {
		return filename[:dotIdx]
	}
	return filename
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
	indentLevel := 1
	hasRead := false
	hasRandom := false
	hasFileOps := false
	hasRun := false
	hasFmt := false
	hasSusu := false
	hasJson := false
	declared := make(map[string]bool)
	gameBlocks := make(map[string]bool)         // tracks named blocks like Game1~ 
	callonceBlocks := make(map[string]bool)     // tracks blocks that can only be called once!

	// tracking for verbose output~ nya!
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

	// first pass: find all named game blocks like "Game1 start"~ like hunting for treats!
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if strings.HasSuffix(line, " start") {
			blockName := strings.TrimSuffix(line, " start")
			// check if it's a named block (not gameloop)
			if blockName != "gameloop" && !strings.Contains(blockName, " ") {
				gameBlocks[blockName] = true
			}
		}
	}

	// second pass: check which blocks use callonce (marked with !callonce syntax)
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "callonce ") {
			blockName := strings.TrimPrefix(line, "callonce ")
			callonceBlocks[blockName] = true
		}
	}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "read -p ") {
			hasRead = true
			hasFmt = true
		} else if strings.HasPrefix(line, "read ") {
			hasRead = true
		}
		if strings.HasPrefix(line, "random ") {
			hasRandom = true
		}
		if strings.HasPrefix(line, "load ") || strings.HasPrefix(line, "save ") {
			hasFileOps = true
		}
		if strings.HasPrefix(line, "loadall ") || strings.HasPrefix(line, "saveall ") {
			hasJson = true
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
		if strings.HasPrefix(line, "peek ") {
			hasFileOps = true
		}
	}

	// figuring out what imports we need, like gathering all our favorite toys before playtime!
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
		out = append(out,
			"    \"bufio\"",
			"    \"strings\"",
		)
	}

	if hasRandom {
		out = append(out,
			"    \"math/rand\"",
			"    \"time\"",
		)
	}

	if hasFileOps {
		out = append(out,
			"    \"io/ioutil\"",
		)
	}

	if hasJson {
		out = append(out,
			"    \"encoding/json\"",
		)
	}

	if hasRun {
		out = append(out,
			"    \"os/exec\"",
		)
	}

	out = append(out,
		")",
		"",
		"// GameState holds all your precious game variables!",
		"type GameState struct {",
	)

	var stateVars []string
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		switch {
		case strings.HasPrefix(line, "give "):
			rest := strings.TrimPrefix(line, "give ")
			parts := strings.SplitN(rest, "=", 2)
			if len(parts) == 2 {
				varName := strings.TrimSpace(parts[0])
				if varName == "gameState" && hasJson {
					fmt.Printf("⚠️  WARNING: 'gameState' is reserved for JSON persistence!\n")
					fmt.Printf("    Renaming to 'gs_state' to avoid struct collision\n")
					varName = "gs_state"
				}
				if !contains(stateVars, varName) {
					stateVars = append(stateVars, varName)
				}
			}
		case strings.HasPrefix(line, "read "):
			varName := strings.TrimSpace(strings.TrimPrefix(line, "read "))
			varName = strings.TrimPrefix(varName, "-p ")
			if !strings.Contains(varName, "\"") && !contains(stateVars, varName) {
				stateVars = append(stateVars, varName)
			}
		case strings.HasPrefix(line, "random "):
			parts := strings.Fields(strings.TrimPrefix(line, "random "))
			if len(parts) >= 1 && !contains(stateVars, parts[0]) {
				stateVars = append(stateVars, parts[0])
			}
		case strings.HasPrefix(line, "score "):
			parts := strings.Fields(strings.TrimPrefix(line, "score "))
			if len(parts) >= 1 && !contains(stateVars, parts[0]) {
				stateVars = append(stateVars, parts[0])
			}
		case strings.HasPrefix(line, "level "):
			parts := strings.Fields(strings.TrimPrefix(line, "level "))
			if len(parts) >= 1 && !contains(stateVars, parts[0]) {
				stateVars = append(stateVars, parts[0])
			}
		}
	}

	for _, varName := range stateVars {
		capitalizedName := strings.ToUpper(varName[:1]) + varName[1:]
		out = append(out, "    "+capitalizedName+" interface{} `json:\""+varName+"\"`")
	}
	// add callonce tracking fields~ so we know which blocks already ran!
	for blockName := range callonceBlocks {
		out = append(out, "    "+blockName+"_called bool `json:\""+blockName+"_called\"`")
	}
	out = append(out, "}")
	out = append(out, "")

	out = append(out, "func main() {",
	)

	if hasRandom {
		out = append(out,
			"    rand.Seed(time.Now().UnixNano())",
		)
	}

	if hasRead {
		out = append(out,
			"    reader := bufio.NewReader(os.Stdin)",
		)
	}

	if hasJson {
		out = append(out,
			"    gameState := &GameState{}",
		)
		// initialize all callonce flags to false~ so they can be called! meow~
		for blockName := range callonceBlocks {
			out = append(out, "    gameState."+blockName+"_called = false")
		}
	}

	if verbose {
		fmt.Println("Import analysis:")
		if hasFmt {
			fmt.Printf("         - fmt: yes\n")
		}
		if hasRead || hasFileOps || hasSusu {
			fmt.Printf("         - os: yes\n")
		}
		if hasRead {
			fmt.Printf("         - bufio, strings: yes (read command detected)\n")
		}
		if hasRandom {
			fmt.Printf("         - math/rand, time: yes (random command detected)\n")
		}
		if hasFileOps {
			fmt.Printf("         - io/ioutil: yes (file operations detected)\n")
		}
		if hasJson {
			fmt.Printf("         - encoding/json: yes (bulk save/load detected)\n")
		}
		if hasRun {
			fmt.Printf("         - os/exec: yes (run command detected)\n")
		}
	}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		indent := strings.Repeat("    ", indentLevel)
		// here comes the fun part! parsing all the silly game commands~
		// paw through each command type and turn it into go code, m3ow.
		switch {
		case line == "gameloop start":
			stats["gameloops"]++
			out = append(out, indent+"for {")
			indentLevel++
			if verbose {
				fmt.Printf("Detected: gameloop (total: %d)\n", stats["gameloops"])
			}
		case line == "gameloop end":
			indentLevel--
			out = append(out, indent+"}")
		// named game blocks like Game1, Boss, etc~ can call 'em whenever!
		case strings.HasSuffix(line, " start"):
			blockName := strings.TrimSuffix(line, " start")
			if gameBlocks[blockName] {
				out = append(out, "")
				out = append(out, "func "+blockName+"() {")
				indentLevel++
				if verbose {
					fmt.Printf("Detected: game block start %s\n", blockName)
				}
			}
		case strings.HasSuffix(line, " end"):
			blockName := strings.TrimSuffix(line, " end")
			if gameBlocks[blockName] {
				indentLevel--
				out = append(out, "}")
				out = append(out, "")
			}
		// call a named game block! like calling a friend over to play~
		case strings.HasPrefix(line, "call "):
			blockName := strings.TrimPrefix(line, "call ")
			if gameBlocks[blockName] {
				out = append(out, indent+blockName+"()")
				if verbose {
					fmt.Printf("Detected: call to game block %s\n", blockName)
				}
			}
		// callonce! call a block only once~ like a special treat!
		case strings.HasPrefix(line, "callonce "):
			blockName := strings.TrimPrefix(line, "callonce ")
			if gameBlocks[blockName] && callonceBlocks[blockName] {
				out = append(out, indent+"if !gameState."+blockName+"_called {")
				out = append(out, indent+"    "+blockName+"()")
				out = append(out, indent+"    gameState."+blockName+"_called = true")
				out = append(out, indent+"} else {")
				out = append(out, indent+"    // nya! already called this block once~")
				out = append(out, indent+"}")
				if verbose {
					fmt.Printf("Detected: callonce to game block %s\n", blockName)
				}
			}
			out = append(out, indent+"fmt.Print(\"\\033[2J\\033[H\")")
		case strings.HasPrefix(line, "border "):
			boxType := strings.TrimPrefix(line, "border ")
			switch boxType {
			case "top":
				out = append(out, indent+"fmt.Println(strings.Repeat(\"═\", 50))")
			case "mid":
				out = append(out, indent+"fmt.Println(strings.Repeat(\"─\", 50))")
			case "bot":
				out = append(out, indent+"fmt.Println(strings.Repeat(\"═\", 50))")
			}
		case strings.HasPrefix(line, "dialog "):
			stats["dialogs"]++
			dialogText := strings.TrimPrefix(line, "dialog ")
			goCode := transpileStringInterpolation(dialogText)
			out = append(out, indent+"fmt.Println(\"════════════════════════════════════════════════════\")")
			out = append(out, indent+"fmt.Println(\" \" + "+goCode+")")
			out = append(out, indent+"fmt.Println(\"════════════════════════════════════════════════════\")")
			if verbose {
				fmt.Printf("Detected: dialog (total: %d)\n", stats["dialogs"])
			}
		case strings.HasPrefix(line, "menu "):
			stats["menus"]++
			menuContent := strings.TrimPrefix(line, "menu ")
			options := strings.Split(menuContent, ",")
			out = append(out, indent+"fmt.Println(\"[MENU]\")")
			for i, opt := range options {
				out = append(out, indent+"fmt.Println(\"  ("+fmt.Sprintf("%d", i+1)+") " + strings.TrimSpace(opt)+"\")")
			}
			if verbose {
				fmt.Printf("Detected: menu (total: %d)\n", stats["menus"])
			}
		case strings.HasPrefix(line, "prompt "):
			stats["prompts"]++
			promptText := strings.TrimPrefix(line, "prompt ")
			goCode := transpileStringInterpolation(promptText)
			out = append(out, indent+"fmt.Print(\" > \" + "+goCode+")")
			if verbose {
				fmt.Printf("Detected: prompt (total: %d)\n", stats["prompts"])
			}
		case strings.HasPrefix(line, "inventory "):
			itemName := strings.TrimPrefix(line, "inventory ")
			parts := strings.Fields(itemName)
			if len(parts) >= 2 {
				varName := parts[0]
				action := parts[1]
				item := parts[2]
				if action == "add" {
					out = append(out, indent+varName+" += \" [\" + "+item+" + \"]\"")
				} else if action == "remove" {
					out = append(out, indent+varName+" = strings.ReplaceAll("+varName+", \" [\" + "+item+" + \"]\", \"\")")
				}
			}
		case strings.HasPrefix(line, "stat "):
			rest := strings.TrimPrefix(line, "stat ")
			var statNameRaw string
			var statValueRaw string

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
				out = append(out, indent+"fmt.Printf("+formatStr+", "+statValue+")")
			}
		// loading and saving files! like hiding your toys and finding them later
		case strings.HasPrefix(line, "load "):
			stats["file_ops"]++
			filePath := strings.TrimPrefix(line, "load ")
			pathExpr := expandPath(filePath)
			varName := extractFilename(filePath)
			out = append(out,
				indent+varName+"Data, _ := ioutil.ReadFile("+pathExpr+")",
				indent+varName+" := string("+varName+"Data)",
			)
			declared[varName] = true
			if verbose {
				fmt.Printf("Detected: load from %s (total file ops: %d)\n", filePath, stats["file_ops"])
			}
		case strings.HasPrefix(line, "save "):
			stats["file_ops"]++
			parts := strings.SplitN(strings.TrimPrefix(line, "save "), " ", 2)
			if len(parts) == 2 {
				varName := parts[0]
				filePath := parts[1]
				pathExpr := expandPath(filePath)
				out = append(out, indent+"ioutil.WriteFile("+pathExpr+", []byte("+varName+"), 0644)")
				if verbose {
					fmt.Printf("Detected: save to %s (total file ops: %d)\n", filePath, stats["file_ops"])
				}
			}
		// BULK SAVE
		case strings.HasPrefix(line, "saveall "):
			stats["file_ops"]++
			filePath := strings.TrimPrefix(line, "saveall ")
			pathExpr := expandPath(filePath)
			out = append(out,
				indent+"for _, varName := range []string{"+stringifyVarList(stateVars)+"} {",
				indent+"    switch varName {",
			)
			for _, varName := range stateVars {
				capitalizedName := strings.ToUpper(varName[:1]) + varName[1:]
				out = append(out, indent+"    case \""+varName+"\":")
				out = append(out, indent+"        gameState."+capitalizedName+" = "+varName)
			}
			out = append(out,
				indent+"    }",
				indent+"}",
				indent+"jsonData, _ = json.MarshalIndent(gameState, \"\", \"  \") // reusing jsonData, like a good maid recycling~",
				indent+"ioutil.WriteFile("+pathExpr+", jsonData, 0644)",
			)
			if verbose {
				fmt.Printf("Detected: saveall to %s (total file ops: %d)\n", filePath, stats["file_ops"])
			}
		// BULK LOAD!
		case strings.HasPrefix(line, "loadall "):
			stats["file_ops"]++
			filePath := strings.TrimPrefix(line, "loadall ")
			pathExpr := expandPath(filePath)
			out = append(out,
				indent+"jsonData, _ := ioutil.ReadFile("+pathExpr+")",
				indent+"json.Unmarshal(jsonData, gameState)",
				indent+"// Restore all variables from state",
			)
			for _, varName := range stateVars {
				capitalizedName := strings.ToUpper(varName[:1]) + varName[1:]
				out = append(out, indent+"if gameState."+capitalizedName+" != nil {")
				out = append(out, indent+"    switch v := gameState."+capitalizedName+".(type) {")
				out = append(out, indent+"    case float64:")
				out = append(out, indent+"        "+varName+" = int(v)")
				out = append(out, indent+"    // strings from json stay as interface{}, meow!")
				out = append(out, indent+"    }")
				out = append(out, indent+"}")
			}
			if verbose {
				fmt.Printf("Detected: loadall from %s (total file ops: %d)\n", filePath, stats["file_ops"])
			}
		case strings.HasPrefix(line, "damage "):
			stats["damage"]++
			parts := strings.Fields(strings.TrimPrefix(line, "damage "))
			if len(parts) >= 2 {
				varName := parts[0]
				damageAmount := stripDollar(parts[1])
				out = append(out, indent+varName+" -= "+damageAmount)
				if verbose {
					fmt.Printf("Detected: damage %s by %s\n", varName, damageAmount)
				}
			}
		case strings.HasPrefix(line, "heal "):
			stats["heal"]++
			parts := strings.Fields(strings.TrimPrefix(line, "heal "))
			if len(parts) >= 2 {
				varName := parts[0]
				healAmount := stripDollar(parts[1])
				out = append(out, indent+varName+" += "+healAmount)
				if verbose {
					fmt.Printf("Detected: heal %s by %s\n", varName, healAmount)
				}
			}
		case strings.HasPrefix(line, "score "):
			stats["score"]++
			parts := strings.Fields(strings.TrimPrefix(line, "score "))
			if len(parts) >= 2 {
				varName := parts[0]
				scoreAmount := stripDollar(parts[1])
				if !declared[varName] {
					out = append(out, indent+varName+" := "+scoreAmount)
					declared[varName] = true
				} else {
					out = append(out, indent+varName+" += "+scoreAmount)
				}
				if verbose {
					fmt.Printf("Detected: score operation on %s\n", varName)
				}
			}
		case strings.HasPrefix(line, "level "):
			stats["level"]++
			parts := strings.Fields(strings.TrimPrefix(line, "level "))
			if len(parts) >= 2 {
				varName := parts[0]
				levelAmount := stripDollar(parts[1])
				if !declared[varName] {
					out = append(out, indent+varName+" := "+levelAmount)
					declared[varName] = true
				} else {
					out = append(out, indent+varName+" = "+levelAmount)
				}
				if verbose {
					fmt.Printf("Detected: level operation on %s\n", varName)
				}
			}
		case strings.HasPrefix(line, "reset "):
			varName := strings.TrimPrefix(line, "reset ")
			out = append(out, indent+varName+" = 0")
		case strings.HasPrefix(line, "random "):
			parts := strings.Fields(strings.TrimPrefix(line, "random "))
			if len(parts) >= 2 {
				varName := parts[0]
				maxVal := parts[1]
				out = append(out, indent+varName+" := rand.Intn("+maxVal+")")
				declared[varName] = true
				if verbose {
					fmt.Printf("Detected: random variable %s (max: %s)\n", varName, maxVal)
				}
			}
			// the core stuff! conditionals, loops, basic i/o~
			// this is like the cat bed - essential and comfy
		case strings.HasPrefix(line, "if "):
			stats["conditionals"]++
			condition := strings.TrimPrefix(line, "if ")
			// handle gameState collision in conditions
			if hasJson {
				condition = strings.ReplaceAll(condition, "gameState", "gs_state")
			}
			out = append(out, indent+"if "+condition+" {")
			indentLevel++
			if verbose {
				fmt.Printf("Detected: conditional (total: %d)\n", stats["conditionals"])
			}
		case strings.HasPrefix(line, "peek "):
			stats["conditionals"]++
			rest := strings.TrimPrefix(line, "peek ")
			if strings.HasPrefix(rest, "\"") {
				endQuote := strings.Index(rest[1:], "\"")
				if endQuote != -1 {
					filePath := rest[1 : endQuote]
					pathExpr := expandPath(filePath)
					out = append(out,
						indent+"if _, err := os.Stat("+pathExpr+"); err == nil {",
					)
					indentLevel++
					if verbose {
						fmt.Printf("Detected: peek check for %s (total: %d)\n", filePath, stats["conditionals"])
					}
				}
			}
		case line == "} else {":
			indentLevel--
			out = append(out, strings.Repeat("    ", indentLevel)+"} else {")
			indentLevel++
		case line == "}":
			indentLevel--
			out = append(out, strings.Repeat("    ", indentLevel)+"}")
		case strings.HasPrefix(line, "write "):
			arg := strings.TrimPrefix(line, "write ")
			goCode := transpileStringInterpolation(arg)
			out = append(out, indent+"fmt.Print("+goCode+")")
		case strings.HasPrefix(line, "writeln "):
			arg := strings.TrimPrefix(line, "writeln ")
			goCode := transpileStringInterpolation(arg)
			out = append(out, indent+"fmt.Println("+goCode+")")
		case strings.HasPrefix(line, "read -p "):
			rest := strings.TrimPrefix(line, "read -p ")
			promptEndIdx := strings.LastIndex(rest, "\"")
			if promptEndIdx > 0 && strings.HasPrefix(rest, "\"") {
				prompt := rest[:promptEndIdx+1]
				varName := strings.TrimSpace(rest[promptEndIdx+1:])

				promptCode := transpileStringInterpolation(prompt)
				out = append(out,
					indent+"fmt.Print("+promptCode+")",
					indent+varName+" := \"\"",
					indent+varName+", _ = reader.ReadString('\\n')",
					indent+varName+" = strings.TrimSpace("+varName+")",
				)
				declared[varName] = true
				stats["variables"]++
				if verbose {
					fmt.Printf("Detected: read with prompt, variable: %s\n", varName)
				}
			}
		case strings.HasPrefix(line, "read "):
			varName := strings.TrimPrefix(line, "read ")
			out = append(out, indent+varName+" := \"\"", indent+varName+", _ = reader.ReadString('\\n')", indent+varName+" = strings.TrimSpace("+varName+")")
			declared[varName] = true
			stats["variables"]++
			if verbose {
				fmt.Printf("Detected: read input, variable: %s\n", varName)
			}
			// assigning variables with 'give'! such a polite way to set values~
			// much nicer than just "x = 10", we say "give x=10" like offering a gift!
		case strings.HasPrefix(line, "give "):
			rest := strings.TrimPrefix(line, "give ")
			parts := strings.SplitN(rest, "=", 2)
			if len(parts) == 2 {
				varName := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				
				// handle gameState collision
				if varName == "gameState" && hasJson {
					varName = "gs_state"
				}

				goValue := transpileStringInterpolation(value)
				
				// if it's not a string interpolation, strip dollar signs from plain variables
				if !strings.HasPrefix(value, "\"") {
					goValue = stripDollar(value)
				}

				if !declared[varName] {
					out = append(out, indent+varName+" := "+goValue)
					declared[varName] = true
					stats["variables"]++
					if verbose {
						fmt.Printf("Detected: variable declaration: %s\n", varName)
					}
				} else {
					out = append(out, indent+varName+" = "+goValue)
				}
			}
			// "susu" means bye-bye! time to end the program and take a nap~
		case strings.HasPrefix(line, "susu"):
			out = append(out, indent+"os.Exit(0)")
			if verbose {
				fmt.Println("Detected: susu (program exit)")
			}
		case strings.HasPrefix(line, "run "):
			stats["run"]++
			cmd := strings.TrimPrefix(line, "run ")
			// run something in terminal, most of it!
			escapedCmd := strings.ReplaceAll(cmd, "\"", "\\\"")
			out = append(out,
				indent+"cmd := exec.Command(\"sh\", \"-c\", \""+escapedCmd+"\")",
				indent+"cmd.Stdout = os.Stdout",
				indent+"cmd.Stderr = os.Stderr",
				indent+"cmd.Run()",
			)
			if verbose {
				fmt.Printf("Detected: run %s\n", cmd)
			}
		}
	}
	out = append(out, "}")

	if verbose {
		fmt.Println("\n=== TRANSPILATION SUMMARY ===")
		fmt.Printf("Variables declared: %d\n", stats["variables"])
		fmt.Printf("Game loops: %d\n", stats["gameloops"])
		fmt.Printf("Dialogs: %d\n", stats["dialogs"])
		fmt.Printf("Menus: %d\n", stats["menus"])
		fmt.Printf("Prompts: %d\n", stats["prompts"])
		fmt.Printf("Conditionals: %d\n", stats["conditionals"])
		fmt.Printf("File operations: %d\n", stats["file_ops"])
		fmt.Printf("Damage operations: %d\n", stats["damage"])
		fmt.Printf("Heal operations: %d\n", stats["heal"])
		fmt.Printf("Score operations: %d\n", stats["score"])
		fmt.Printf("Level operations: %d\n", stats["level"])
		fmt.Printf("Run commands: %d\n", stats["run"])
		fmt.Println("=======================")
	}

	return strings.Join(out, "\n")
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}


func stringifyVarList(vars []string) string {
	var parts []string
	for _, v := range vars {
		parts = append(parts, "\""+v+"\"")
	}
	return strings.Join(parts, ", ")
}
