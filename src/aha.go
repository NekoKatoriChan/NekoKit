package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: neko file.nk")
		return
	}

	content, _ := os.ReadFile(os.Args[1])
	code := string(content)

	// Biến đơn giản
	vars := map[string]string{}

	lines := strings.Split(code, "\n")
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		
		if strings.HasPrefix(line, "print ") {
			// print "Hello World"
			text := strings.Trim(line[6:], `"`)
			fmt.Println(text)
		
		} else if strings.HasPrefix(line, "if ") {
			// if x > 5
			condition := line[3:]
			if evalCondition(condition, vars) {
				// Thực thi khối lệnh
				i++
				for i < len(lines) && strings.TrimSpace(lines[i]) != "end" {
					execLine(strings.TrimSpace(lines[i]), vars)
					i++
				}
			}
		
		} else if strings.HasPrefix(line, "var ") {
			// var x = 10
			parts := strings.Split(line[4:], "=")
			if len(parts) == 2 {
				vars[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}
}

func evalCondition(cond string, vars map[string]string) bool {
	// Đơn giản: so sánh số
	if strings.Contains(cond, ">") {
		parts := strings.Split(cond, ">")
		return toNumber(parts[0], vars) > toNumber(parts[1], vars)
	}
	return true
}

func toNumber(s string, vars map[string]string) int {
	s = strings.TrimSpace(s)
	if val, ok := vars[s]; ok {
		s = val
	}
	var num int
	fmt.Sscanf(s, "%d", &num)
	return num
}

func execLine(line string, vars map[string]string) {
	if strings.HasPrefix(line, "print ") {
		text := strings.Trim(line[6:], `"`)
		fmt.Println(text)
	}
}
