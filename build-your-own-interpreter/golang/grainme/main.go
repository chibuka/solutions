package main

import (
	"fmt"
	"log"
	"os"
)

func peek(fileContent string, idx int) (byte, bool) {
	if idx+1 < len(fileContent) {
		return fileContent[idx+1], true
	}
	return 0, false
}

func main() {
	command, args := os.Args[1], os.Args[2:]
	tokens := map[string]string{
		"(":  "LEFT_PAREN",
		")":  "RIGHT_PAREN",
		"}":  "RIGHT_BRACE",
		"{":  "LEFT_BRACE",
		"*":  "STAR",
		".":  "DOT",
		",":  "COMMA",
		"+":  "PLUS",
		"-":  "MINUS",
		"/":  "SLASH",
		";":  "SEMICOLON",
		"=":  "EQUAL",
		"==": "EQUAL_EQUAL",
		"!":  "BANG",
		"!=": "BANG_EQUAL",
		">":  "GREATER",
		"<":  "LESS",
		">=": "GREATER_EQUAL",
		"<=": "LESS_EQUAL",
	}

	switch command {
	case "tokenize":
		filePath := args[0]
		contentBytes, _ := os.ReadFile(filePath)
		fileContent := string(contentBytes)
		tokenNotFound := false
		lineNum := 1

		for idx := 0; idx < len(fileContent); idx++ {
			currChar := fileContent[idx]

			// handling spaces, newline, tabs
			if currChar == ' ' || currChar == '\n' || currChar == '\t' {
				if currChar == '\n' {
					lineNum++
				}
				continue
			}

			nextChar, ok := peek(fileContent, idx)
			if ok {
				// check if we're into a comment //
				if currChar == '/' && nextChar == '/' {
					// skip next '/'
					idx++
					// skip all chars after // until \n
					for idx < len(fileContent) && fileContent[idx] != '\n' {
						idx++
					}
					lineNum++
					continue
				}

				currNext := string(currChar) + string(nextChar)
				val, ok := tokens[string(currNext)]
				if ok {
					fmt.Printf("%s %s null\n", val, string(currNext))
					idx++
					continue
				}
			}

			val, ok := tokens[string(currChar)]
			if !ok {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", lineNum, string(currChar))
				tokenNotFound = true
				continue
			}
			fmt.Printf("%s %s null\n", val, string(currChar))
		}
		fmt.Printf("EOF  null\n")
		if tokenNotFound {
			os.Exit(65)
		}
	default:
		log.Fatalf("fatal error")
	}
}
