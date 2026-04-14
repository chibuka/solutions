package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chibuka/build-your-own-interpreter/golang/grainme/internal"
)

func main() {
	command, args := os.Args[1], os.Args[2:]
	filePath := args[0]
	loxSource := internal.GetSource(filePath)
	s := internal.Scanner{
		Source:  loxSource,
		Current: 0,
		Line:    1,
		Errors:  false,
	}

	switch command {
	case "tokenize":
		tokens := s.Tokenize()
		for _, token := range tokens {
			// output format: <token_type> <lexeme> <literal>
			fmt.Printf("%s %s %s\n", token.Type, token.Lexeme, token.Literal)
		}
		fmt.Println("EOF  null")

		if s.Errors {
			os.Exit(65)
		}
	case "parse":
		internal.Parse(s)
	default:
		log.Fatalf("fatal error")
	}
}
