/*
 * main.go - mygit: a minimal git implementation
 *
 * Plumbing commands: init, cat-file, hash-object, ls-tree, write-tree, commit-tree
 */
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/95/testers/git/go/commands"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: mygit <command> [<args>...]\n")
	}

	// this can be: "init", "commit", "push"...
	subCommand := os.Args[1]
	args := os.Args[2:]

	switch subCommand {
	case "init":
		commands.HandleInit()
	case "cat-file":
		commands.HandleCatfile(args)
	case "hash-object":
		commands.HandleHashObject(args)
	case "ls-tree":
		commands.HandleLsTree(args)
	case "write-tree":
		commands.HandleWriteTree()
	case "commit-tree":
		commands.HandleCommitTree(args)
	default:
		fmt.Println("sub-command not supported")
	}
}
