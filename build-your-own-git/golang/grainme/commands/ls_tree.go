package commands

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func HandleLsTree(args []string) {
	// we should handle flags (-w...)
	if len(args) != 2 {
		log.Fatalf("usage: mygit ls-tree <flag> <hash>\n")
	}

	flag := args[0]
	// TODO: for now we only support "--name-only" as a flag
	// ls-tree + --name-only reads a tree object and prints the name
	// of the entries in the same order they appear in the object.
	if flag != "--name-only" {
		log.Fatalf("flags other than --name-only are not supported yet\n")
	}

	sha := args[1]
	dir := sha[:2]
	filePath := sha[2:]
	gitObjectPath := fmt.Sprintf(".git/objects/%s/%s", dir, filePath)

	f, _ := os.Open(gitObjectPath)
	defer f.Close()
	r, _ := zlib.NewReader(f)
	defer r.Close()
	data, _ := io.ReadAll(r)

	//
	// ls-tree output example:
	//
	// tree 66\0
	// 100644 a.txt\0<binary SHA>
	// 100644 b.txt\0<binary SHA>
	//               ^
	// 			     should be converted to hex for proper output
	//

	nulPos := bytes.IndexByte(data, 0)
	typeAndLength := string(data[:nulPos])
	// e.g: tree 66
	_ = typeAndLength

	// now we should parse the entries
	pos := nulPos + 1
	for pos < len(data) {
		nextNulPos := bytes.IndexByte(data[pos:], 0)
		if nextNulPos == -1 {
			break
		}

		fileModeAndName := strings.Split(string(data[pos:pos+nextNulPos]), " ")
		_, fileName := fileModeAndName[0], fileModeAndName[1]

		fmt.Println(fileName)

		// +1 for the NUL byte
		// +20 for the <binary SHA>
		pos += nextNulPos + 1 + 20
	}
}
