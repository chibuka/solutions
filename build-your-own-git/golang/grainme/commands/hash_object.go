package commands

import (
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func HandleHashObject(args []string) {
	// we should handle flags (-w...)
	if len(args) != 2 {
		log.Fatalf("usage: mygit hash-object <flag> <filepath>\n")
	}

	flag := args[0]
	// TODO: for now we only support "-w" as a flag
	if flag != "-w" {
		log.Fatalf("flags other than -w are not supported yet\n")
	}

	filePath := args[1]

	// raw file bytes
	content, _ := os.ReadFile(filePath)
	contentSize := len(content)

	// \x00 is the NUL character (one byte)
	header := fmt.Sprintf("blob %d\x00", contentSize)
	store := header + string(content)

	shaBytes := sha1.Sum([]byte(store))
	sha := hex.EncodeToString(shaBytes[:])
	//
	// previous line is equivalent to this:
	// sha := fmt.Sprintf("%x", shaBytes)
	//

	dir := sha[:2]
	filename := sha[2:]

	gitObjectPath := fmt.Sprintf(".git/objects/%s/%s", dir, filename)

	//
	// TODO: i don't know if this needed, because os.Create panics otherwise.
	//
	_ = os.MkdirAll(filepath.Dir(gitObjectPath), 0755)
	f, _ := os.Create(gitObjectPath)

	// compresses the file
	w := zlib.NewWriter(f)
	w.Write([]byte(store))
	w.Close()

	fmt.Println(sha)
}
