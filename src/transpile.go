package main

import "strings"

func Transpile(input string) string {
	lines := strings.Split(input, "\n")
	var out []string

	out = append(out, "package main")
	out = append(out, "import \"fmt\"")
	out = append(out, "func main() {")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "write ") {
			content := strings.TrimPrefix(line, "write ")
			out = append(out, "    fmt.Print("+content+")")
		}
	}

	out = append(out, "}")
	return strings.Join(out, "\n")
}
