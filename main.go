package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var store = make(map[string]string)

func cmd_PING(args []string) string {
	if len(args) == 1 {
		return "+PONG\r\n"
	}

	if len(args) == 2 {
		return encodeBulkString(args[1])
	}

	return "-ERR wrong number of arguments for 'PING' command\r\n"
}

func cmd_ECHO(args []string) string {
	if len(args) == 1 {
		return "-ERR wrong number of arguments for 'ECHO' command\r\n"
	}

	res := strings.Join(args[1:], " ")

	return encodeBulkString(res)
}

func cmd_SET(args []string) string {
	if len(args) != 3 {
		return "-ERR wrong number of arguments for 'SET' command\r\n"
	}

	key := args[1]
	value := args[2]

	store[key] = value

	return "+OK\r\n"
}

func cmd_GET(args []string) string {
	if len(args) != 2 {
		return "-ERR wrong number of arguments for 'GET' command\r\n"
	}

	key := args[1]

	value, exists := store[key]

	if !exists {
		return "$-1\r\n"
	}

	return encodeBulkString(value)
}

func cmd_DOCS() string {
	return "+OK\r\n"
}

func handleCommand(args []string) string {
	if len(args) == 0 {
		return "-ERR empty command\r\n"
	}

	cmd := strings.ToUpper(args[0])

	switch cmd {
	case "PING":
		return cmd_PING(args)

	case "ECHO":
		return cmd_ECHO(args)

	case "SET":
		return cmd_SET(args)

	case "GET":
		return cmd_GET(args)

	case "COMMAND":
		if len(args) > 1 && strings.ToUpper(args[1]) == "DOCS" {
			return cmd_DOCS()
		}
	}

	return fmt.Sprintf("-ERR unknown command '%s'\r\n", cmd)
}

func encodeBulkString(s string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}

func simple_string(s string) string {
	return fmt.Sprintf("+%s\r\n", s)
}

func encodeError(s string) string {
	return fmt.Sprintf("-%s\r\n", s)
}

func integer(i int) string {
	return fmt.Sprintf(":%d\r\n", i)
}

func bulk_string(s string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}

func readRESPArray(r *bufio.Reader) ([]string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSuffix(line, "\r\n")

	if len(line) == 0 || line[0] != '*' {
		return nil, fmt.Errorf("expected RESP array")
	}

	n, err := strconv.Atoi(line[1:])
	if err != nil {
		return nil, err
	}

	args := make([]string, 0, n)

	for i := 0; i < n; i++ {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSuffix(line, "\r\n")

		if len(line) == 0 || line[0] != '$' {
			return nil, fmt.Errorf("expected bulk string")
		}

		length, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, err
		}

		data := make([]byte, length)

		_, err = io.ReadFull(r, data)
		if err != nil {
			return nil, err
		}

		crlf := make([]byte, 2)

		_, err = io.ReadFull(r, crlf)
		if err != nil {
			return nil, err
		}

		args = append(args, string(data))
	}

	return args, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		args, err := readRESPArray(reader)
		if err != nil {
			return
		}

		response := handleCommand(args)

		fmt.Print(response)
	}
}
