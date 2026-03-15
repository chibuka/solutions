/*
 * commit_tree.go - Create a commit object
 *
 * Usage: mygit commit-tree <tree> [-p <parent>] -m <message>
 *
 * Commit object content (plain text, unlike tree's binary format):
 *
 *     tree <40-char hex SHA>
 *     parent <40-char hex SHA>       ← omitted for initial commit
 *     author <name> <email> <time>
 *     committer <name> <email> <time>
 *
 *     <commit message>
 */
package commands

import (
	"fmt"
	"log"
)

func HandleCommitTree(args []string) {
	if len(args) < 3 {
		log.Fatalf("usage: mygit commit-tree <tree> [-p <parent>] -m <message>\n")
	}

	treeSha := args[0]
	parentSha := ""
	message := ""

	// parse flags dynamically — handles both with and without parent
	i := 1
	for i < len(args) {
		switch args[i] {
		case "-p":
			parentSha = args[i+1]
			i += 2
		case "-m":
			message = args[i+1]
			i += 2
		default:
			i++
		}
	}

	// build commit content line by line
	var content []byte

	content = append(content, []byte(fmt.Sprintf("tree %s\n", treeSha))...)

	// no parent line for the initial commit
	if parentSha != "" {
		content = append(content, []byte(fmt.Sprintf("parent %s\n", parentSha))...)
	}

	// hardcoded author/committer — good enough for the challenge
	content = append(content, []byte("author Grainme <grainme@example.com> 1234567890 +0000\n")...)
	content = append(content, []byte("committer Grainme <grainme@example.com> 1234567890 +0000\n")...)
	content = append(content, []byte(fmt.Sprintf("\n%s\n", message))...)

	sha := writeObject("commit", content)
	fmt.Println(sha)
}
