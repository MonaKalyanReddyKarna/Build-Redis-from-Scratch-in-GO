package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cmd_PING(args []string) string {
	if len(args) == 1 {
		return "+PONG\r\n"
	} else {
		message := args[1]
		return encodeBulkString(message)
	}
}

func cmd_ECHO(args []string) string {
	res := strings.Join(args[1:], " ")
	return encodeBulkString(res)
}

func handleCommand(args []string) string {
	cmd := strings.ToUpper(args[0])

	switch cmd {
	case "PING":
		return cmd_PING(args)
		// TODO: Return "+PONG\r\n" for no args
	case "ECHO":
		return cmd_ECHO(args)
		// TODO: Return bulk string for PING <message>
	}

	return fmt.Sprintf("-ERR unknown command '%s'\r\n", cmd)
}

func encodeBulkString(s string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		args := parseArgs(line)
		response := handleCommand(args)
		fmt.Print(response)
	}
}

func parseArgs(line string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	for _, ch := range line {
		switch {
		case ch == '"' && !inQuotes:
			inQuotes = true
		case ch == '"' && inQuotes:
			inQuotes = false
		case ch == ' ' && !inQuotes:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(ch)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}
