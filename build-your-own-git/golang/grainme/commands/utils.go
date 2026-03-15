package commands

import (
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

func writeBlob(filePath string) string {
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

	_ = os.MkdirAll(filepath.Dir(gitObjectPath), 0755)
	f, _ := os.Create(gitObjectPath)

	// compresses the file
	w := zlib.NewWriter(f)
	w.Write([]byte(store))
	w.Close()

	return sha
}

func writeTree(buffer []byte) string {
	// raw file bytes
	bufferSize := len(buffer)

	// \x00 is the NUL character (one byte)
	header := fmt.Sprintf("tree %d\x00", bufferSize)
	store := header + string(buffer)

	shaBytes := sha1.Sum([]byte(store))
	sha := hex.EncodeToString(shaBytes[:])

	dir := sha[:2]
	filename := sha[2:]

	gitObjectPath := fmt.Sprintf(".git/objects/%s/%s", dir, filename)

	_ = os.MkdirAll(filepath.Dir(gitObjectPath), 0755)
	f, _ := os.Create(gitObjectPath)

	// compresses the file
	w := zlib.NewWriter(f)
	w.Write([]byte(store))
	w.Close()

	return sha
}
