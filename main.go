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

// parseInteractiveLine expects "<value> -> <base>" or "<value> <base>".
func parseInteractiveLine(line string) (value string, base int, err error) {
	line = strings.ReplaceAll(line, "->", " ")
	fields := strings.Fields(line)
	if len(fields) != 2 {
		return "", 0, fmt.Errorf(`expected format "<value> -> <base>", e.g. "0x25B9D2 -> 2"`)
	}
	base, err = strconv.Atoi(fields[1])
	if err != nil {
		return "", 0, fmt.Errorf("target base must be a number: %w", err)
	}
	return fields[0], base, nil
}

// withBasePrefix adds a conventional prefix for common bases so output
func withBasePrefix(s string, base int) string {
	switch base {
	case converter.Binary:
		return "0b" + s
	case converter.Octal:
		return "0o" + s
	case converter.Hexadecimal:
		return "0x" + s
	default:
		return s
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `numconv - convert numbers between bases 2-36

Usage:
  numconv [--from BASE] [--to BASE] <value>   (flags must come before value)
  numconv                                     # interactive mode

Flags:`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, `
Examples:
  numconv --to 2 0x25B9D2
  numconv --from 2 --to 16 1010111001001001
  numconv --to 16 255`)
}

func main() {
	from := flag.Int("from", 0, "source base (2-36). 0 = auto-detect from 0x/0b/0o prefix, else decimal")
	to := flag.Int("to", 10, "target base (2-36)")
	flag.Usage = printUsage
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		runInteractive()
		return
	}

	result, err := converter.Convert(args[0], *from, *to)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println(withBasePrefix(result, *to))
}


