/*
 * cat_file.go - Read and display git objects
 *
 * Usage: mygit cat-file -p <sha>
 *
 * Reads an object from .git/objects, decompresses it,
 * strips the header, and prints the raw content.
 */
package commands

import (
	"log"
	"os"
)

func HandleCatfile(args []string) {
	if len(args) != 2 {
		log.Fatalf("usage: mygit cat-file -p <hash>\n")
	}

	flag := args[0]
	if flag != "-p" {
		log.Fatalf("only -p flag is supported\n")
	}

	_, content := readObject(args[1])
	os.Stdout.Write(content)
}
