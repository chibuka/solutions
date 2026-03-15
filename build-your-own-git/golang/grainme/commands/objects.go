/*
 * objects.go - Git object database I/O
 *
 * Every git object (blob, tree, commit) is stored the same way:
 *
 *     <type> <size>\0<content>  →  SHA-1  →  zlib  →  .git/objects/ab/cdef...
 *
 * This file owns all reads and writes to .git/objects.
 * Command files should never touch the object store directly.
 */
package commands

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// writeObject stores a git object and returns its 40-char hex SHA.
func writeObject(objType string, content []byte) string {
	header := fmt.Sprintf("%s %d\x00", objType, len(content))
	store := append([]byte(header), content...)

	shaBytes := sha1.Sum(store)
	sha := hex.EncodeToString(shaBytes[:])

	gitObjectPath := fmt.Sprintf(".git/objects/%s/%s", sha[:2], sha[2:])

	_ = os.MkdirAll(filepath.Dir(gitObjectPath), 0755)
	f, _ := os.Create(gitObjectPath)

	w := zlib.NewWriter(f)
	_, _ = w.Write(store)
	_ = w.Close()
	_ = f.Close()

	return sha
}

// readObject reads a git object and returns its type and content.
func readObject(sha string) (string, []byte) {
	gitObjectPath := fmt.Sprintf(".git/objects/%s/%s", sha[:2], sha[2:])

	f, _ := os.Open(gitObjectPath)
	defer f.Close()

	r, _ := zlib.NewReader(f)
	defer r.Close()

	data, _ := io.ReadAll(r)

	// split at the NUL byte: "<type> <size>\0<content>"
	nul := bytes.IndexByte(data, 0)
	header := string(data[:nul])
	content := data[nul+1:]

	// extract type from header (e.g. "blob 14" → "blob")
	objType := header[:bytes.IndexByte([]byte(header), ' ')]

	return objType, content
}
