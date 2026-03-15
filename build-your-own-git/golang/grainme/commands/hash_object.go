/*
 * hash_object.go - Compute object hash and optionally store it
 *
 * Usage: mygit hash-object -w <file>
 *
 * Reads a file, wraps it as a blob object, computes its SHA-1,
 * and stores it in .git/objects when -w is given.
 */
package commands

import (
	"fmt"
	"log"
	"os"
)

func HandleHashObject(args []string) {
	if len(args) != 2 {
		log.Fatalf("usage: mygit hash-object -w <filepath>\n")
	}

	flag := args[0]
	if flag != "-w" {
		log.Fatalf("only -w flag is supported\n")
	}

	content, _ := os.ReadFile(args[1])
	sha := writeObject("blob", content)
	fmt.Println(sha)
}
