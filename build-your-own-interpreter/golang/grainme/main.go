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
		if s.Errors {
			os.Exit(65)
		}
	case "parse":
		tokens := s.Tokenize()
		if s.Errors {
			os.Exit(65)
		}

		p := internal.NewParse(tokens)
		expr, err := p.Parse()
		if err != nil {
			log.Fatalf("parsing went wrong %v", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", expr.String())
	default:
		log.Fatalf("fatal error")
	}
}
