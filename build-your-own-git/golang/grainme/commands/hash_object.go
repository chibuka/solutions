package commands

import (
	"fmt"
	"log"
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

	//
	// TODO: we should check the object type?
	// right now we're assuming Blob, which is false!
	//
	hash := writeBlob(filePath)

	fmt.Println(hash)
}
