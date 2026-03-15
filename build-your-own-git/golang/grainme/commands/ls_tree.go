/*
 * ls_tree.go - List the contents of a tree object
 *
 * Usage: mygit ls-tree --name-only <sha>
 *
 * Reads a tree object and parses its binary entries.
 * Tree entry format: <mode> <name>\0<20-byte raw SHA>
 */
package commands

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

func HandleLsTree(args []string) {
	if len(args) != 2 {
		log.Fatalf("usage: mygit ls-tree --name-only <hash>\n")
	}

	flag := args[0]
	if flag != "--name-only" {
		log.Fatalf("only --name-only flag is supported\n")
	}

	_, content := readObject(args[1])

	// parse binary entries: <mode> <name>\0<20-byte SHA>
	pos := 0
	for pos < len(content) {
		nul := bytes.IndexByte(content[pos:], 0)
		if nul == -1 {
			break
		}

		modeAndName := strings.SplitN(string(content[pos:pos+nul]), " ", 2)
		fmt.Println(modeAndName[1])

		// skip past NUL + 20-byte raw SHA
		pos += nul + 1 + 20
	}
}
