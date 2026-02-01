package main

import "strings"

func Transpile(input string) string {
    lines := strings.Split(input, "\n")
    var out []string
    indentLevel := 1 

    out = append(out, "package main", "import (", "    \"fmt\"", "    \"os\"", ")", "func main() {")

    for _, raw := range lines {
        line := strings.TrimSpace(raw)
        indent := strings.Repeat("    ", indentLevel)

        switch {
        //basix comands, blin
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
            out = append(out, indent+"fmt.Print("+strings.TrimPrefix(line, "write ")+")")
        case strings.HasPrefix(line, "writeln "):
            out = append(out, indent+"fmt.Println("+strings.TrimPrefix(line, "writeln ")+")")
            //kys, a.k.a kill your system's process
        case strings.HasPrefix(line, "susu"):
            out = append(out, indent+"os.Exit(0)")
        }
    }
    out = append(out, "}")
    return strings.Join(out, "\n")
}
