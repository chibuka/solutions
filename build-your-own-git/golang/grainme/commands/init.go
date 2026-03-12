package commands

import (
	"fmt"
	"log"
	"os"
)

func HandleInit() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Couldn't identify current working directory")
	}

	// owner + group + other
	// 4(r) + 2(w) + 1(x)
	os.MkdirAll(".git/objects", 0755)
	os.MkdirAll(".git/refs", 0755)
	os.WriteFile(".git/HEAD", []byte("ref: refs/heads/main\n"), 0622)

	fmt.Println("Initialized empty Git repository in", currentDir+"/.git")
}
