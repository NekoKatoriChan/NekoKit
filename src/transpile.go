package main

import "strings"



// Help

func transpileStringInterpolation(arg string) string {
    if strings.HasPrefix(arg, "\"") && strings.HasSuffix(arg, "\"") && strings.Contains(arg, "$") {
        content := arg[1 : len(arg)-1]
        
        var vars []string
        var formatStr strings.Builder
        formatStr.WriteString("\"")
        
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
        
        result := "fmt.Sprintf(" + formatStr.String()
        for _, v := range vars {
            result += ", " + v
        }
        result += ")"
        return result
    }
    return arg
}

// eh, work as a variable, example "writeln "hello, $s" 
func stripDollar(s string) string {
    if strings.HasPrefix(s, "$") {
        return s[1:]
    }
    return s
}

func isAlphaNumericOrUnderscore(b byte) bool {
    return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

func expandPath(path string) string {
    if strings.HasPrefix(path, "~/") {
        return "os.ExpandEnv(\"$HOME" + path[1:] + "\")"
    }
    return "\"" + path + "\""
}

func extractFilename(path string) string {
    parts := strings.Split(path, "/")
    filename := parts[len(parts)-1]
    dotIdx := strings.LastIndex(filename, ".")
    if dotIdx > 0 {
        return filename[:dotIdx]
    }
    return filename
}

// Transpiler core
func Transpile(input string) string {
    lines := strings.Split(input, "\n")
    var out []string
    indentLevel := 1 
    hasRead := false
    hasRandom := false
    hasFileOps := false
    declared := make(map[string]bool)

    for _, raw := range lines {
        line := strings.TrimSpace(raw)
        if strings.HasPrefix(line, "read ") {
            hasRead = true
        }
        if strings.HasPrefix(line, "random ") {
            hasRandom = true
        }
        if strings.HasPrefix(line, "load ") || strings.HasPrefix(line, "save ") {
            hasFileOps = true
        }
    }
		// basic import thing
    out = append(out,
        "package main",
        "import (",
        "    \"fmt\"",
        "    \"os\"",
    )

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

    out = append(out,
        ")",
        "func main() {",
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

    for _, raw := range lines {
        line := strings.TrimSpace(raw)
        if line == "" {
            continue
        }
        indent := strings.Repeat("    ", indentLevel)
		// The game loops.
        switch {
        case line == "gameloop start":
            out = append(out, indent+"for {")
            indentLevel++
        case line == "gameloop end":
            indentLevel--
            out = append(out, indent+"}")
        case line == "clear":
            out = append(out, indent+"fmt.Print(\"\\033[2J\\033[H\")")
        case strings.HasPrefix(line, "border "):
            boxType := strings.TrimPrefix(line, "border ")
            switch boxType {
            case "simple":
                out = append(out, indent+"fmt.Println(\"┌─────────────────────┐\")")
                out = append(out, indent+"fmt.Println(\"│                     │\")")
                out = append(out, indent+"fmt.Println(\"└─────────────────────┘\")")
            case "double":
                out = append(out, indent+"fmt.Println(\"╔═════════════════════╗\")")
                out = append(out, indent+"fmt.Println(\"║                     ║\")")
                out = append(out, indent+"fmt.Println(\"╚═════════════════════╝\")")
            case "thick":
                out = append(out, indent+"fmt.Println(\"▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓\")")
                out = append(out, indent+"fmt.Println(\"▓                   ▓\")")
                out = append(out, indent+"fmt.Println(\"▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓\")")
            }
           // Dialog, for some reason. preset gamey thing
        case strings.HasPrefix(line, "dialog "):
            dialogText := strings.TrimPrefix(line, "dialog ")
            goCode := transpileStringInterpolation(dialogText)
            out = append(out, indent+"fmt.Println(\">>> \" + "+goCode+")")
        case strings.HasPrefix(line, "menu "):
            parts := strings.SplitN(strings.TrimPrefix(line, "menu "), " ", 2)
            if len(parts) >= 1 {
                out = append(out, indent+"fmt.Println(\"[1] Option 1\")", indent+"fmt.Println(\"[2] Option 2\")", indent+"fmt.Println(\"[3] Option 3\")")
            }
        case strings.HasPrefix(line, "prompt "):
            promptText := strings.TrimPrefix(line, "prompt ")
            goCode := transpileStringInterpolation(promptText)
            out = append(out, indent+"fmt.Print(\" > \" + "+goCode+")")
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
            // This is load/save file
        case strings.HasPrefix(line, "load "):
            filePath := strings.TrimPrefix(line, "load ")
            pathExpr := expandPath(filePath)
            varName := extractFilename(filePath)
            out = append(out, 
                indent+varName+"Data, _ := ioutil.ReadFile("+pathExpr+")",
                indent+varName+" := string("+varName+"Data)",
            )
            declared[varName] = true
        case strings.HasPrefix(line, "save "):
            parts := strings.SplitN(strings.TrimPrefix(line, "save "), " ", 2)
            if len(parts) == 2 {
                varName := parts[0]
                filePath := parts[1]
                pathExpr := expandPath(filePath)
                out = append(out, indent+"ioutil.WriteFile("+pathExpr+", []byte("+varName+"), 0644)")
            }
        case strings.HasPrefix(line, "damage "):
            parts := strings.Fields(strings.TrimPrefix(line, "damage "))
            if len(parts) >= 2 {
                varName := parts[0]
                damageAmount := stripDollar(parts[1])
                out = append(out, indent+varName+" -= "+damageAmount)
            }
        case strings.HasPrefix(line, "heal "):
            parts := strings.Fields(strings.TrimPrefix(line, "heal "))
            if len(parts) >= 2 {
                varName := parts[0]
                healAmount := stripDollar(parts[1])
                out = append(out, indent+varName+" += "+healAmount)
            }
        case strings.HasPrefix(line, "score "):
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
            }
        case strings.HasPrefix(line, "level "):
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
            }
            // Core
        case strings.HasPrefix(line, "if "):
            out = append(out, indent+"if "+strings.TrimPrefix(line, "if ")+" {")
            indentLevel++
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
            }
        case strings.HasPrefix(line, "read "):
            varName := strings.TrimPrefix(line, "read ")
            out = append(out, indent+varName+" := \"\"", indent+varName+", _ = reader.ReadString('\\n')", indent+varName+" = strings.TrimSpace("+varName+")")
            declared[varName] = true
            // Exit 0, kill the program.
        case strings.HasPrefix(line, "susu"):
            out = append(out, indent+"os.Exit(0)")
        }
    }
    out = append(out, "}")
    return strings.Join(out, "\n")
}
