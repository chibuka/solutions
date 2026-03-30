/*
 * This is a very hard task to do in one stage.
 */
package commands

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// packfile object type codes (3-bit field)
const (
	packCommit   = 1
	packTree     = 2
	packBlob     = 3
	packTag      = 4
	packOfsDelta = 6
	packRefDelta = 7
)

func HandleClone(args []string) {
	if len(args) != 2 {
		log.Fatalf("usage: mygit clone <url> <directory>\n")
	}
	url := strings.TrimRight(args[0], "/")
	dir := args[1]

	// discover refs — find HEAD's SHA and branch name
	headSha, branchName := discoverRefs(url)

	// create target directory and init .git
	if err := os.Mkdir(dir, 0755); err != nil {
		log.Fatalf("mkdir %s: %s\n", dir, err)
	}
	if err := os.Chdir(dir); err != nil {
		log.Fatalf("chdir %s: %s\n", dir, err)
	}
	os.MkdirAll(".git/objects", 0755)
	os.MkdirAll(".git/refs/heads", 0755)

	// fetch packfile and unpack all objects into .git/objects
	fetchPack(url, headSha)

	// create branch ref and HEAD
	refPath := fmt.Sprintf(".git/refs/heads/%s", branchName)
	os.WriteFile(refPath, []byte(headSha+"\n"), 0644)
	os.WriteFile(".git/HEAD", []byte(fmt.Sprintf("ref: refs/heads/%s\n", branchName)), 0644)

	// checkout HEAD commit
	checkoutCommit(headSha)
}

// readPktLine reads one pkt-line from r. Returns nil on flush (0000).
func readPktLine(r io.Reader) []byte {
	lenHex := make([]byte, 4)
	if _, err := io.ReadFull(r, lenHex); err != nil {
		return nil
	}
	n, err := strconv.ParseUint(string(lenHex), 16, 32)
	if err != nil || n <= 4 {
		return nil
	}
	data := make([]byte, n-4)
	io.ReadFull(r, data)
	return data
}

func discoverRefs(url string) (headSha, branchName string) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(fmt.Sprintf("%s/info/refs?service=git-upload-pack", url))
	if err != nil {
		log.Fatalf("ref discovery: %s\n", err)
	}
	defer resp.Body.Close()

	// skip service announcement packets until flush
	for readPktLine(resp.Body) != nil {
	}

	// first ref line: "<sha> HEAD\0<capabilities>\n"
	first := readPktLine(resp.Body)
	if first == nil {
		log.Fatal("no refs advertised\n")
	}
	line := strings.SplitN(string(first), "\x00", 2)[0]
	parts := strings.SplitN(strings.TrimRight(line, "\n"), " ", 2)
	headSha = parts[0]

	// remaining lines — find the ref whose SHA matches HEAD
	headRef := ""
	for {
		pkt := readPktLine(resp.Body)
		if pkt == nil {
			break
		}
		s := strings.TrimRight(string(pkt), "\n")
		p := strings.SplitN(s, " ", 2)
		if len(p) == 2 && p[0] == headSha {
			headRef = p[1]
		}
	}
	if headRef == "" {
		log.Fatal("no branch matches HEAD\n")
	}
	branchName = headRef[strings.LastIndex(headRef, "/")+1:]
	return
}

func fetchPack(url, wantSha string) {
	// pkt-line request body:
	//   "0032want <40-char sha>\n"   (length 50 = 0x32)
	//   "0000"                       (flush)
	//   "0009done\n"                 (length 9)
	body := fmt.Sprintf("0032want %s\n00000009done\n", wantSha)

	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequest("POST",
		fmt.Sprintf("%s/git-upload-pack", url),
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/x-git-upload-pack-request")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("upload-pack: %s\n", err)
	}
	defer resp.Body.Close()

	// skip the server's NAK pkt-line
	readPktLine(resp.Body)

	unpackPackfile(resp.Body)
}

func unpackPackfile(r io.Reader) {
	packData, _ := io.ReadAll(r)

	if string(packData[:4]) != "PACK" {
		log.Fatal("bad packfile signature\n")
	}

	// verify checksum (last 20 bytes = SHA-1 of everything before it)
	checksum := packData[len(packData)-20:]
	packData = packData[:len(packData)-20]
	h := sha1.Sum(packData)
	if !bytes.Equal(checksum, h[:]) {
		log.Fatal("packfile checksum mismatch\n")
	}

	buf := bytes.NewBuffer(packData)
	readNBytes(8, buf) // skip "PACK" (4) + version (4)

	count := binary.BigEndian.Uint32(readNBytes(4, buf))

	for i := uint32(0); i < count; i++ {
		objType, size := readPackTypeAndSize(buf)

		switch objType {
		case packCommit:
			writeObject("commit", zlibRead(size, buf))
		case packTree:
			writeObject("tree", zlibRead(size, buf))
		case packBlob:
			writeObject("blob", zlibRead(size, buf))
		case packTag:
			zlibRead(size, buf) // read but discard
		case packOfsDelta:
			log.Fatal("OFS_DELTA not supported\n")
		case packRefDelta:
			baseHash := fmt.Sprintf("%x", readNBytes(20, buf))
			delta := zlibRead(size, buf)
			resolveRefDelta(baseHash, delta)
		}
	}
}

// readNBytes reads exactly n bytes from r.
func readNBytes(n int, r io.Reader) []byte {
	data := make([]byte, n)
	io.ReadFull(r, data)
	return data
}

// readPackTypeAndSize decodes the variable-length type+size header.
//
//	first byte:       [MSB] [type:3] [size:4]
//	subsequent bytes: [MSB] [size:7]
func readPackTypeAndSize(buf *bytes.Buffer) (objType int, size uint64) {
	b, _ := buf.ReadByte()
	objType = int((b >> 4) & 0x07)
	size = uint64(b & 0x0F)

	if b&0x80 == 0 {
		return
	}
	shift := uint(4)
	for {
		b, _ = buf.ReadByte()
		size |= uint64(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}
	return
}

// zlibRead decompresses exactly `size` bytes from a zlib stream in r.
// bytes.Buffer implements io.ByteReader, so Go's flate reads byte-by-byte
// and leaves the buffer positioned right after the compressed data.
func zlibRead(size uint64, r io.Reader) []byte {
	zr, _ := zlib.NewReader(r)
	data := make([]byte, size)
	io.ReadFull(zr, data)
	zr.Close()
	return data
}

func resolveRefDelta(baseSha string, delta []byte) {
	baseType, baseContent := readObject(baseSha)

	buf := bytes.NewBuffer(delta)
	srcLen := readVarInt(buf)
	tgtLen := readVarInt(buf)

	if uint64(len(baseContent)) != srcLen {
		log.Fatalf("delta src len mismatch: %d vs %d\n", len(baseContent), srcLen)
	}

	var result []byte
	for buf.Len() > 0 {
		cmd, _ := buf.ReadByte()

		if cmd&0x80 == 0 {
			// insert: next cmd bytes go straight into result
			result = append(result, readNBytes(int(cmd&0x7F), buf)...)
		} else {
			// copy from base object
			var offset, size uint32
			for i := uint(0); i < 4; i++ {
				if cmd&(1<<i) != 0 {
					b, _ := buf.ReadByte()
					offset |= uint32(b) << (8 * i)
				}
			}
			for i := uint(0); i < 3; i++ {
				if cmd&(0x10<<i) != 0 {
					b, _ := buf.ReadByte()
					size |= uint32(b) << (8 * i)
				}
			}
			if size == 0 {
				size = 0x10000
			}
			result = append(result, baseContent[offset:offset+size]...)
		}
	}

	if uint64(len(result)) != tgtLen {
		log.Fatalf("delta target len mismatch: %d vs %d\n", len(result), tgtLen)
	}

	writeObject(baseType, result)
}

// readVarInt reads a variable-length integer (used in delta size headers).
func readVarInt(buf *bytes.Buffer) uint64 {
	var val uint64
	var shift uint
	for {
		b, _ := buf.ReadByte()
		val |= uint64(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}
	return val
}

// checkoutCommit reads a commit object, extracts the tree SHA, and checks it out.
func checkoutCommit(sha string) {
	_, content := readObject(sha)

	// first line: "tree <40-char sha>"
	treeSha := strings.SplitN(string(content), "\n", 2)[0]
	treeSha = strings.TrimPrefix(treeSha, "tree ")

	checkoutTree(".", treeSha)
}

// checkoutTree recursively reconstructs files from a tree object.
func checkoutTree(base, sha string) {
	_, content := readObject(sha)

	// tree entries: <mode> <name>\0<20-byte raw sha>
	pos := 0
	for pos < len(content) {
		nul := bytes.IndexByte(content[pos:], 0)
		if nul == -1 {
			break
		}

		parts := strings.SplitN(string(content[pos:pos+nul]), " ", 2)
		mode, name := parts[0], parts[1]

		rawSha := content[pos+nul+1 : pos+nul+21]
		entrySha := hex.EncodeToString(rawSha)
		pos += nul + 21

		entryPath := filepath.Join(base, name)

		if mode == "40000" {
			os.MkdirAll(entryPath, 0755)
			checkoutTree(entryPath, entrySha)
		} else {
			_, blob := readObject(entrySha)
			os.WriteFile(entryPath, blob, 0644)
		}
	}
}
