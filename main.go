package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"strconv"
	"flag"

	"numconv/converter"
)

// runInteractive prompts the user in a loop so they don't have to
// remember flags. Type "quit" or Ctrl+D to exit.
func runInteractive() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("numconv interactive mode. Enter values like: 0x25B9D2 -> 2")
	fmt.Println("Formats accepted for the value: 0x.. (hex), 0b.. (binary), 0o.. (octal), or plain decimal.")
	fmt.Println(`Type "quit" to exit.`)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == "quit" || line == "exit" {
			return
		}

		value, targetBase, err := parseInteractiveLine(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, "  error:", err)
			continue
		}

		result, err := converter.Convert(value, 0, targetBase)
		if err != nil {
			fmt.Fprintln(os.Stderr, "  error:", err)
			continue
		}
		fmt.Println("  " + withBasePrefix(result, targetBase))
	}
}


func main() {
}
