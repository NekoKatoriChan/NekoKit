package main

import (
	"fmt"
	"strings"
)

// this is the transpiler for making silly game scripts into proper go code!
// made with love and fish treats~ each bug fix was also done with fish treats!

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

		// if we have variables, use fmt.Sprintf; otherwise just return the escaped string
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

	// we now track two separate output slices - mainOut for main(), blockOut for named funcs~
	// think of it like two separate food bowls. do NOT mix the kibble!! nya!!
	var mainOut []string
	var blockOut []string
	inNamedBlock := false // are we currently writing inside a named block func?

	// two indent levels too~ one for main, one for whatever block we're in
	// like two separate yarn balls that must never tangle~
	mainIndent := 1
	blockIndent := 1
	// save/restore mainIndent when diving into a named block
	savedMainIndent := 1

	// emitLine sends code to the right output slice - very important, meow!
	emitLine := func(s string) {
		if inNamedBlock {
			blockOut = append(blockOut, s)
		} else {
			mainOut = append(mainOut, s)
		}
	}

	// getIndent returns the correct indentation based on where we are~
	getIndent := func() string {
		if inNamedBlock {
			return strings.Repeat("    ", blockIndent)
		}
		return strings.Repeat("    ", mainIndent)
	}

	// bumpIndent increments the active indent level~
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
	hasJson := false

	// a variable declared in Boss() has no idea what lives in main(), and that's correct!
	// like how the cat in the bedroom doesn't know what snacks are in the kitchen~
	mainDeclared := make(map[string]bool)
	blockDeclared := make(map[string]bool)

	// getDeclared returns the declared map for whoever is currently speaking~
	getDeclared := func() map[string]bool {
		if inNamedBlock {
			return blockDeclared
		}
		return mainDeclared
	}

	// no more "cmd redeclared in this block" disasters~ each run gets its own tiny bowl~
	runCmdCount := 0

	gameBlocks := make(map[string]bool)     // tracks named blocks like Game1~
	callonceBlocks := make(map[string]bool) // tracks blocks that can only be called once!

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
			// check if it's a named block (not gameloop, no spaces in name)
			if blockName != "gameloop" && !strings.Contains(blockName, " ") {
				gameBlocks[blockName] = true
			}
		}
	}

	// second pass: find callonce blocks~ these little ones can only run once, how exclusive!
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if strings.HasPrefix(line, "callonce ") {
			blockName := strings.TrimPrefix(line, "callonce ")
			callonceBlocks[blockName] = true
		}
	}

	// FIX BUG 7: callonce needs gameState machinery even without saveall/loadall!
	// we track this separately so we don't accidentally import encoding/json just for callonce~
	// it's like... needing a cat door even if you don't have a full cat hotel, meow!
	hasCallonce := len(callonceBlocks) > 0

	// third pass: scan for what imports we actually need! no unnecessary baggage~
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
		"// GameState holds all your precious game variables! guard them with your life~",
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
	// add callonce tracking fields~ so we know which blocks already ran, like a stamp card!
	for blockName := range callonceBlocks {
		out = append(out, "    "+blockName+"_called bool `json:\""+blockName+"_called\"`")
	}
	out = append(out, "}")
	out = append(out, "")
	out = append(out, "func main() {")

	if hasRandom {
		out = append(out, "    rand.Seed(time.Now().UnixNano())")
	}
	if hasRead {
		out = append(out, "    reader := bufio.NewReader(os.Stdin)")
	}

	// FIX BUG 7 (continued): gameState is needed for BOTH hasJson AND hasCallonce!
	// we warm up the gameState cat bed regardless, as long as someone needs it~
	if hasJson || hasCallonce {
		out = append(out, "    gameState := &GameState{}")
		// initialize all callonce flags to false~ ready to be called for the first time!
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
		if hasCallonce {
			fmt.Printf("         - gameState: yes (callonce detected, even without saveall)\n")
		}
	}

	// here comes the fun part! parsing all the silly game commands~
	// paw through each command type and turn it into go code, m3ow.
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
			if verbose {
				fmt.Printf("Detected: gameloop (total: %d)\n", stats["gameloops"])
			}

		case line == "gameloop end":
			// FIX: decrement FIRST, then use the new level for the closing brace!
			// previously it used the old (too-deep) indent before decrementing~ oopsie!
			bumpIndent(-1)
			if inNamedBlock {
				emitLine(strings.Repeat("    ", blockIndent) + "}")
			} else {
				emitLine(strings.Repeat("    ", mainIndent) + "}")
			}

		// they're not inside main() anymore~ they have their own house now, nya!
		case strings.HasSuffix(line, " start"):
			blockName := strings.TrimSuffix(line, " start")
			if gameBlocks[blockName] {
				// save our place in main so we can come back after the block is done~
				savedMainIndent = mainIndent
				inNamedBlock = true
				blockIndent = 1
				// Boss's 'choice' variable has NOTHING to do with main's 'choice'~
				// they live in completely separate apartments, meow!
				blockDeclared = make(map[string]bool)
				blockOut = append(blockOut, "")
				blockOut = append(blockOut, "// game block '"+blockName+"'~ enter at your own risk! nya~")
				blockOut = append(blockOut, "func "+blockName+"() {")
				if verbose {
					fmt.Printf("Detected: game block start %s\n", blockName)
				}
			}

		case strings.HasSuffix(line, " end"):
			blockName := strings.TrimSuffix(line, " end")
			if gameBlocks[blockName] {
				blockIndent--
				blockOut = append(blockOut, "}")
				blockOut = append(blockOut, "")
				// return to main scope and restore the indent we saved before~
				inNamedBlock = false
				mainIndent = savedMainIndent
				blockDeclared = make(map[string]bool) // clean the block's toy box on exit~
				if verbose {
					fmt.Printf("Detected: game block end %s\n", blockName)
				}
			}

		// call a named game block! like calling a friend over to play~
		case strings.HasPrefix(line, "call "):
			blockName := strings.TrimPrefix(line, "call ")
			if gameBlocks[blockName] {
				emitLine(indent + blockName + "()")
				if verbose {
					fmt.Printf("Detected: call to game block %s\n", blockName)
				}
			}

		// the accidental screen-clear that was leaking out has been caught and scolded~
		// 'clear' is its own separate case below where it belongs! very organized~
		case strings.HasPrefix(line, "callonce "):
			blockName := strings.TrimPrefix(line, "callonce ")
			if gameBlocks[blockName] && callonceBlocks[blockName] {
				emitLine(indent + "if !gameState." + blockName + "_called {")
				emitLine(indent + "    " + blockName + "()")
				emitLine(indent + "    gameState." + blockName + "_called = true")
				emitLine(indent + "} else {")
				emitLine(indent + "    // nya! already called this block once, no second helpings~")
				emitLine(indent + "}")
				if verbose {
					fmt.Printf("Detected: callonce to game block %s\n", blockName)
				}
			}
			// NOTE: the screen clear that used to sneak in here has been removed!
			// it was living in the wrong case like a cat in the wrong box~ very rude!

		// FIX BUG 2 (the other half): clear finally gets its own home!
		// it was homeless before, now it lives right here where it belongs~ meow!
		case line == "clear":
			emitLine(indent + `fmt.Print("\033[2J\033[H")`)

		case strings.HasPrefix(line, "border "):
			boxType := strings.TrimPrefix(line, "border ")
			switch boxType {
			case "top":
				emitLine(indent + `fmt.Println(strings.Repeat("═", 50))`)
			case "mid":
				emitLine(indent + `fmt.Println(strings.Repeat("─", 50))`)
			case "bot":
				emitLine(indent + `fmt.Println(strings.Repeat("═", 50))`)
			}

		case strings.HasPrefix(line, "dialog "):
			stats["dialogs"]++
			dialogText := strings.TrimPrefix(line, "dialog ")
			goCode := transpileStringInterpolation(dialogText)
			emitLine(indent + `fmt.Println("════════════════════════════════════════════════════")`)
			emitLine(indent + `fmt.Println(" " + ` + goCode + ")")
			emitLine(indent + `fmt.Println("════════════════════════════════════════════════════")`)
			if verbose {
				fmt.Printf("Detected: dialog (total: %d)\n", stats["dialogs"])
			}

		case strings.HasPrefix(line, "menu "):
			stats["menus"]++
			menuContent := strings.TrimPrefix(line, "menu ")
			options := strings.Split(menuContent, ",")
			emitLine(indent + `fmt.Println("[MENU]")`)
			for i, opt := range options {
				emitLine(indent + `fmt.Println("  (` + fmt.Sprintf("%d", i+1) + `) ` + strings.TrimSpace(opt) + `")`)
			}
			if verbose {
				fmt.Printf("Detected: menu (total: %d)\n", stats["menus"])
			}

		case strings.HasPrefix(line, "prompt "):
			stats["prompts"]++
			promptText := strings.TrimPrefix(line, "prompt ")
			goCode := transpileStringInterpolation(promptText)
			emitLine(indent + `fmt.Print(" > " + ` + goCode + ")")
			if verbose {
				fmt.Printf("Detected: prompt (total: %d)\n", stats["prompts"])
			}

		case strings.HasPrefix(line, "inventory "):
			itemName := strings.TrimPrefix(line, "inventory ")
			parts := strings.Fields(itemName)
			// also fixed a latent panic here! we check len >= 3 before accessing parts[2]~
			// previously it would crash like a cat knocking a glass off a table~ on purpose!
			if len(parts) >= 3 {
				varName := parts[0]
				action := parts[1]
				item := parts[2]
				if action == "add" {
					emitLine(indent + varName + ` += " [" + ` + item + ` + "]"`)
				} else if action == "remove" {
					emitLine(indent + varName + ` = strings.ReplaceAll(` + varName + `, " [" + ` + item + ` + "]", "")`)
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
				emitLine(indent + "fmt.Printf(" + formatStr + ", " + statValue + ")")
			}

		// loading and saving files! like hiding your toys and finding them later~
		// loading the same file twice in a long script will no longer explode, meow!
		case strings.HasPrefix(line, "load "):
			stats["file_ops"]++
			filePath := strings.TrimPrefix(line, "load ")
			pathExpr := expandPath(filePath)
			varName := extractFilename(filePath)
			dataVar := varName + "Data"
			// check if we already declared the data variable~ careful careful~
			if !declared[dataVar] {
				emitLine(indent + dataVar + ", _ := ioutil.ReadFile(" + pathExpr + ")")
				declared[dataVar] = true
			} else {
				emitLine(indent + dataVar + ", _ = ioutil.ReadFile(" + pathExpr + ")")
			}
			// and the string variable too, don't forget~
			if !declared[varName] {
				emitLine(indent + varName + " := string(" + dataVar + ")")
				declared[varName] = true
			} else {
				emitLine(indent + varName + " = string(" + dataVar + ")")
			}
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
				emitLine(indent + "ioutil.WriteFile(" + pathExpr + ", []byte(" + varName + "), 0644)")
				if verbose {
					fmt.Printf("Detected: save to %s (total file ops: %d)\n", filePath, stats["file_ops"])
				}
			}

		// BULK SAVE~ saving ALL the things at once! very efficient kitty~
		case strings.HasPrefix(line, "saveall "):
			stats["file_ops"]++
			filePath := strings.TrimPrefix(line, "saveall ")
			pathExpr := expandPath(filePath)
			emitLine(indent + "for _, varName := range []string{" + stringifyVarList(stateVars) + "} {")
			emitLine(indent + "    switch varName {")
			for _, varName := range stateVars {
				capitalizedName := strings.ToUpper(varName[:1]) + varName[1:]
				emitLine(indent + "    case \"" + varName + "\":")
				emitLine(indent + "        gameState." + capitalizedName + " = " + varName)
			}

			// check if jsonData was already declared, reuse if so~
			jsonDeclare := "="
			if !declared["jsonData"] {
				jsonDeclare = ":="
				declared["jsonData"] = true
			}

			emitLine(indent + "    }")
			emitLine(indent + "}")
			emitLine(indent + "jsonData, _ " + jsonDeclare + " json.MarshalIndent(gameState, \"\", \"  \") // reusing jsonData, like a good recycling kitty~")
			emitLine(indent + "ioutil.WriteFile(" + pathExpr + ", jsonData, 0644)")
			if verbose {
				fmt.Printf("Detected: saveall to %s (total file ops: %d)\n", filePath, stats["file_ops"])
			}

		// BULK LOAD~ restoring ALL the things! welcome back, precious variables~
		case strings.HasPrefix(line, "loadall "):
			stats["file_ops"]++
			filePath := strings.TrimPrefix(line, "loadall ")
			pathExpr := expandPath(filePath)

			jsonDeclare := "="
			if !declared["jsonData"] {
				jsonDeclare = ":="
				declared["jsonData"] = true
			}

			emitLine(indent + "jsonData, _ " + jsonDeclare + " ioutil.ReadFile(" + pathExpr + ")")
			emitLine(indent + "json.Unmarshal(jsonData, gameState)")
			emitLine(indent + "// Restore all variables from state~")
			for _, varName := range stateVars {
				capitalizedName := strings.ToUpper(varName[:1]) + varName[1:]
				emitLine(indent + "if gameState." + capitalizedName + " != nil {")
				emitLine(indent + "    switch v := gameState." + capitalizedName + ".(type) {")
				emitLine(indent + "    case float64:")
				emitLine(indent + "        " + varName + " = int(v)")
				emitLine(indent + "    // strings from json stay as interface{}, meow!")
				emitLine(indent + "    }")
				emitLine(indent + "}")
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
				emitLine(indent + varName + " -= " + damageAmount)
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
				emitLine(indent + varName + " += " + healAmount)
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
					emitLine(indent + varName + " := " + scoreAmount)
					declared[varName] = true
				} else {
					emitLine(indent + varName + " += " + scoreAmount)
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
					emitLine(indent + varName + " := " + levelAmount)
					declared[varName] = true
				} else {
					emitLine(indent + varName + " = " + levelAmount)
				}
				if verbose {
					fmt.Printf("Detected: level operation on %s\n", varName)
				}
			}

		case strings.HasPrefix(line, "reset "):
			varName := strings.TrimPrefix(line, "reset ")
			emitLine(indent + varName + " = 0")

		// re-randomizing the same variable in a long loop won't explode anymore~
		// it now reassigns like a well-behaved kitty instead of crashing~
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
			emitLine(indent + "if " + condition + " {")
			bumpIndent(1)
			if verbose {
				fmt.Printf("Detected: conditional (total: %d)\n", stats["conditionals"])
			}

		// now correctly uses endQuote+1 to include the full filename~ nya!!
		case strings.HasPrefix(line, "peek "):
			stats["conditionals"]++
			rest := strings.TrimPrefix(line, "peek ")
			if strings.HasPrefix(rest, "\"") {
				endQuote := strings.Index(rest[1:], "\"")
				if endQuote != -1 {
					// FIX: endQuote is relative to rest[1:], so real index is endQuote+1
					// rest[1 : endQuote+1] correctly captures the full filename~
					filePath := rest[1 : endQuote+1]
					pathExpr := expandPath(filePath)
					emitLine(indent + "if _, err := os.Stat(" + pathExpr + "); err == nil {")
					bumpIndent(1)
					if verbose {
						fmt.Printf("Detected: peek check for %s (total: %d)\n", filePath, stats["conditionals"])
					}
				}
			}

		case line == "} else {":
			bumpIndent(-1)
			if inNamedBlock {
				emitLine(strings.Repeat("    ", blockIndent) + "} else {")
			} else {
				emitLine(strings.Repeat("    ", mainIndent) + "} else {")
			}
			bumpIndent(1)

		case line == "}":
			bumpIndent(-1)
			if inNamedBlock {
				emitLine(strings.Repeat("    ", blockIndent) + "}")
			} else {
				emitLine(strings.Repeat("    ", mainIndent) + "}")
			}

		case strings.HasPrefix(line, "write "):
			arg := strings.TrimPrefix(line, "write ")
			goCode := transpileStringInterpolation(arg)
			emitLine(indent + "fmt.Print(" + goCode + ")")

		case strings.HasPrefix(line, "writeln "):
			arg := strings.TrimPrefix(line, "writeln ")
			goCode := transpileStringInterpolation(arg)
			emitLine(indent + "fmt.Println(" + goCode + ")")

		// a variable read twice in a long script will reuse, not redeclare~
		// like asking for the same toy twice - we already have it, no need to buy another!
		case strings.HasPrefix(line, "read -p "):
			rest := strings.TrimPrefix(line, "read -p ")
			promptEndIdx := strings.LastIndex(rest, "\"")
			if promptEndIdx > 0 && strings.HasPrefix(rest, "\"") {
				prompt := rest[:promptEndIdx+1]
				varName := strings.TrimSpace(rest[promptEndIdx+1:])
				promptCode := transpileStringInterpolation(prompt)

				// only declare with := on first use~ subsequent reads just reassign!
				if !declared[varName] {
					emitLine(indent + varName + ` := ""`)
					declared[varName] = true
				}
				emitLine(indent + "fmt.Print(" + promptCode + ")")
				emitLine(indent + varName + `, _ = reader.ReadString('\n')`)
				emitLine(indent + varName + " = strings.TrimSpace(" + varName + ")")
				stats["variables"]++
				if verbose {
					fmt.Printf("Detected: read with prompt, variable: %s\n", varName)
				}
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
			if verbose {
				fmt.Printf("Detected: read input, variable: %s\n", varName)
			}

		// assigning variables with 'give'! such a polite way to set values~
		// much nicer than just "x = 10", we say "give x=10" like offering a gift! meow~
		case strings.HasPrefix(line, "give "):
			rest := strings.TrimPrefix(line, "give ")
			parts := strings.SplitN(rest, "=", 2)
			if len(parts) == 2 {
				varName := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				// handle gameState collision~
				if varName == "gameState" && hasJson {
					varName = "gs_state"
				}

				goValue := transpileStringInterpolation(value)

				// if it's not a string interpolation, strip dollar signs from plain variables
				if !strings.HasPrefix(value, "\"") {
					goValue = stripDollar(value)
				}

				if !declared[varName] {
					emitLine(indent + varName + " := " + goValue)
					declared[varName] = true
					stats["variables"]++
					if verbose {
						fmt.Printf("Detected: variable declaration: %s\n", varName)
					}
				} else {
					emitLine(indent + varName + " = " + goValue)
				}
			}

		// "susu" means bye-bye! time to end the program and take a nap~
		case strings.HasPrefix(line, "susu"):
			emitLine(indent + "os.Exit(0)")
			if verbose {
				fmt.Println("Detected: susu (program exit)")
			}

		// like giving every shell command its own little collar tag, meow!
		case strings.HasPrefix(line, "run "):
			stats["run"]++
			cmd := strings.TrimPrefix(line, "run ")
			escapedCmd := strings.ReplaceAll(cmd, "\"", "\\\"")
			cmdVar := fmt.Sprintf("_cmd%d", runCmdCount)
			runCmdCount++
			emitLine(indent + cmdVar + ` := exec.Command("sh", "-c", "` + escapedCmd + `")`)
			emitLine(indent + cmdVar + ".Stdout = os.Stdout")
			emitLine(indent + cmdVar + ".Stderr = os.Stderr")
			emitLine(indent + cmdVar + ".Run()")
			if verbose {
				fmt.Printf("Detected: run %s (as %s)\n", cmd, cmdVar)
			}
		}
	}

	// main body goes first, then we close main() with its }, then ALL named blocks after!
	// this is the correct Go file structure, nya!! like a well-organized litter box!
	out = append(out, mainOut...)
	out = append(out, "}") // properly close main() before the named block functions!
	out = append(out, blockOut...)

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
