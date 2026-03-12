package commands

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"
)

func HandleCatfile(args []string) {
	// we should handle flags (-p, -e...)
	if len(args) != 2 {
		log.Fatalf("usage: mygit cat-file <flag> <hash>\n")
	}

	flag := args[0]
	// TODO: for now we only support "-p" as a flag
	// check: https://git-scm.com/docs/git-cat-file
	if flag != "-p" {
		log.Fatalf("flags other than -p are not supported yet\n")
	}

	sha := args[1]
	dir, file := sha[:2], sha[2:]

	filePath := fmt.Sprintf(".git/objects/%s/%s", dir, file)

	//
	// ---- in a nutshell:
	// Git never stores raw file bytes.
	// it wraps them in a typed (e.g: "blob", "tree"...), size-prefixed header (e.g: 14, check example below)
	// hashes the whole thing (using SHA-1 - maybe) and zlib-compresses it before touching the disk.
	//
	f, _ := os.Open(filePath)
	r, _ := zlib.NewReader(f)
	data, _ := io.ReadAll(r)

	// next step: strip the header `blob <size>\0`
	// we should find the last position of the header
	// e.g: "blob 14\0Hello World"
	// 				^
	// 			   this
	nul := bytes.IndexByte(data, 0)

	//
	// TODO [optional]: we can validate the size of the content
	// e.g: check if not 14 and error.log "error: object file ... is corrupted"
	//

	// this is equivalent to: fmt.Print(string(data[nul+1:]))
	os.Stdout.Write(data[nul+1:])
}
