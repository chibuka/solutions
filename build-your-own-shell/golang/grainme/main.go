package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("$ ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}

		lineSplit := strings.Split(strings.TrimRight(line, "\n"), " ")
		command, args := lineSplit[0], lineSplit[1:]

		commands_supported := []string{"exit", "echo", "type"}

		switch command {
		case "exit":
			if len(args) == 0 {
				os.Exit(0)
			}
			exitCode, err := strconv.Atoi(strings.Join(args, ""))
			if err != nil {
				log.Fatalf("error: %v\n", err)
			}
			os.Exit(exitCode)
		case "echo":
			fmt.Println(strings.Join(args, " "))
		case "type":
			arg := strings.Join(args, " ")
			if slices.Contains(commands_supported, arg) {
				fmt.Printf("%s is a shell builtin\n", arg)
			} else if path, found := findInPath(arg); found {
				fmt.Printf("%s is %s\n", arg, path)
			} else {
				fmt.Printf("%s: not found\n", arg)
			}
		default:
			fmt.Printf("%s: command not found\n", command)
		}

	}
}

func findInPath(cmd string) (string, bool) {
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		path := filepath.Join(dir, cmd)
		if info, err := os.Stat(path); err == nil && !info.IsDir() && info.Mode().Perm()&0111 != 0 {
			return path, true
		}
	}
	return "", false
}
