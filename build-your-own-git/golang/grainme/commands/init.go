/*
 * init.go - Initialize a new git repository
 *
 * Creates the .git directory structure:
 *     .git/HEAD       → points to refs/heads/main
 *     .git/objects/   → the object database
 *     .git/refs/      → branch and tag pointers
 */
package commands

import (
	"fmt"
	"os"
)

func HandleInit() {
	cwd, _ := os.Getwd()

	_ = os.MkdirAll(".git/objects", 0755)
	_ = os.MkdirAll(".git/refs", 0755)
	_ = os.WriteFile(".git/HEAD", []byte("ref: refs/heads/main\n"), 0644)

	fmt.Println("Initialized empty Git repository in", cwd+"/.git")
}
