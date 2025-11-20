package main

import "strings"

func Transpile(input string) string {
        lines := strings.Split(input, "\n")
        var out []string
        indentLevel := 1 

        out = append(out,
                "package main",
                "import \"fmt\"",
                "func main() {",
        )

        for _, raw := range lines {
                line := strings.TrimSpace(raw)
                indent := strings.Repeat("    ", indentLevel)

                switch {
                case strings.HasPrefix(line, "if "):
                        condition := strings.TrimPrefix(line, "if ")
                        out = append(out, indent+"if "+condition+" {")
                        indentLevel++

                case line == "} else {":
                        indentLevel--
                        indent = strings.Repeat("    ", indentLevel)
                        out = append(out, indent+"} else {")
                        indentLevel++

                case line == "}":
                        indentLevel--
                        indent = strings.Repeat("    ", indentLevel)
                        out = append(out, indent+"}")
				// I made 2 basic output, write and writeln
                case strings.HasPrefix(line, "write "):
                        content := strings.TrimPrefix(line, "write ")
                        out = append(out, indent+"fmt.Print("+content+")")

                case strings.HasPrefix(line, "writeln "):
                        content := strings.TrimPrefix(line, "writeln ")
                        out = append(out, indent+"fmt.Println("+content+")")
				// Joke syntax (=ω=)
                case strings.HasPrefix(line, "susu "):
                        content := strings.TrimPrefix(line, "susu ")
                        out = append(out, indent+"// "+content)
                }
        }

        out = append(out, "}")
        return strings.Join(out, "\n")
}
