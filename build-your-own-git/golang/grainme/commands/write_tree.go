/*
 * this commands recursively snapshots the current working directory
 * into the object store, creating tree objects for each directory,
 * and prints the 40-chars SHA of the root tree.
 *
 */
package commands

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// take a dir path, writes a tree object to .git/objects and return its SHA.
func writeTreeForDir(dir string) string {
	//
	// os.ReadDir already returns them sorted, so you're good here.
	// otherwise, we needed to do it manually.
	//
	entries, _ := os.ReadDir(dir)

	buffer := make([]byte, 0)
	for _, entry := range entries {
		entryPath := filepath.Join(dir, entry.Name())
		sha := ""
		mode := "100644"

		if entry.Name() == ".git" {
			continue
		}

		if entry.Type().IsRegular() {
			sha = writeBlob(entryPath)
		}
		if entry.Type().IsDir() {
			// mode for dirs 40000
			mode = "40000"
			sha = writeTreeForDir(entryPath)
		}

		hexToRaw, _ := hex.DecodeString(sha)

		// The format is <mode> <name>\0<raw-sha>
		buffer = append(buffer, []byte(mode+" "+entry.Name()+"\x00")...)
		buffer = append(buffer, hexToRaw...)
	}

	return writeTree(buffer)
}

func HandleWriteTree() {
	// Walk the current directory recursively (skip .git/)
	cwd, _ := os.Getwd()
	sha := writeTreeForDir(cwd)
	fmt.Println(sha)
}
