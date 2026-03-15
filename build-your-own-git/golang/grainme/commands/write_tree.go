/*
 * write_tree.go - Snapshot the working directory into tree objects
 *
 * Usage: mygit write-tree
 *
 * Recursively walks the current directory (skipping .git/),
 * creates blob objects for files and tree objects for directories,
 * and prints the root tree SHA.
 *
 * Tree entry binary format: <mode> <name>\0<20-byte raw SHA>
 * Mode 100644 = regular file, 40000 = directory.
 */
package commands

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// writeTreeForDir recursively writes tree objects for a directory.
func writeTreeForDir(dir string) string {
	// os.ReadDir returns entries sorted — git requires sorted tree entries
	entries, _ := os.ReadDir(dir)

	var buf []byte
	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}

		entryPath := filepath.Join(dir, entry.Name())
		var sha string
		var mode string

		if entry.Type().IsDir() {
			mode = "40000"
			sha = writeTreeForDir(entryPath)
		} else {
			mode = "100644"
			content, _ := os.ReadFile(entryPath)
			sha = writeObject("blob", content)
		}

		rawSha, _ := hex.DecodeString(sha)

		// <mode> <name>\0<20-byte raw SHA>
		buf = append(buf, []byte(mode+" "+entry.Name()+"\x00")...)
		buf = append(buf, rawSha...)
	}

	return writeObject("tree", buf)
}

func HandleWriteTree() {
	cwd, _ := os.Getwd()
	fmt.Println(writeTreeForDir(cwd))
}
