package internal

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

var reservedWords = map[string]TokenType{
	"and":    AND,
	"or":     OR,
	"class":  CLASS,
	"if":     IF,
	"else":   ELSE,
	"false":  FALSE,
	"true":   TRUE,
	"for":    FOR,
	"fun":    FUN,
	"nil":    NIL,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"var":    VAR,
	"while":  WHILE,
}

type Scanner struct {
	Source  string
	Current int
	Line    int
	Errors  bool
}

func (s *Scanner) isAtEnd() bool {
	return s.Current >= len(s.Source)
}

func (s *Scanner) advance() byte {
	if s.isAtEnd() {
		return 0
	}
	curr := s.Source[s.Current]
	s.Current++
	return curr
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.Source[s.Current]
}

func (s *Scanner) peekNext() byte {
	if s.Current+1 >= len(s.Source) {
		return 0
	}
	return s.Source[s.Current+1]
}

func (s *Scanner) match(c byte) bool {
	if s.Current >= len(s.Source) {
		return false
	}
	if s.Source[s.Current] == c {
		s.Current++
		return true
	}
	return false
}

func (s *Scanner) scanString() (Token, bool) {
	var b strings.Builder
	for {
		if s.isAtEnd() {
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", s.Line)
			s.Errors = true
			return Token{}, false
		}
		currByte := s.peek()
		if currByte == '"' {
			s.Current++ // consume closing "
			break
		}
		b.WriteByte(currByte)
		s.Current++
	}

	return Token{
		Type:    STRING,
		Lexeme:  fmt.Sprintf("\"%s\"", b.String()),
		Literal: b.String(),
	}, true
}

func (s *Scanner) scanIdentifier(first byte) Token {
	var b strings.Builder
	b.WriteByte(first)
	for {
		currByte := s.peek()
		if (!unicode.IsLetter(rune(currByte)) && !unicode.IsDigit(rune(currByte)) && currByte != '_') || s.isAtEnd() {
			break
		}
		b.WriteByte(currByte)
		s.Current++
	}

	val, ok := reservedWords[b.String()]
	if ok {
		return s.makeToken(val, b.String(), "null")
	}

	return s.makeToken(IDENTIFIER, b.String(), "null")
}

func (s *Scanner) scanNumber(first byte) Token {
	var b strings.Builder
	b.WriteByte(first)
	for {
		currByte := s.peek()
		if !unicode.IsDigit(rune(currByte)) && currByte != '.' {
			break
		}
		b.WriteByte(currByte)
		s.Current++
	}

	return Token{
		Type:    NUMBER,
		Lexeme:  b.String(),
		Literal: formatNumber(b.String()),
	}
}

func (s *Scanner) scanComment() {
	for s.peek() != '\n' && !s.isAtEnd() {
		s.Current++
	}
}

func (s *Scanner) makeToken(t TokenType, lexeme, literal string) Token {
	return Token{Type: t, Lexeme: lexeme, Literal: literal}
}

func (s *Scanner) Tokenize() []Token {
	tokens := make([]Token, 0)

	for !s.isAtEnd() {
		lexeme := rune(s.advance())
		switch lexeme {
		case ')':
			tokens = append(tokens, s.makeToken(RIGHT_PAREN, ")", "null"))
		case '(':
			tokens = append(tokens, s.makeToken(LEFT_PAREN, "(", "null"))
		case '{':
			tokens = append(tokens, s.makeToken(LEFT_BRACE, "{", "null"))
		case '}':
			tokens = append(tokens, s.makeToken(RIGHT_BRACE, "}", "null"))
		case '*':
			tokens = append(tokens, s.makeToken(STAR, "*", "null"))
		case '.':
			tokens = append(tokens, s.makeToken(DOT, ".", "null"))
		case ',':
			tokens = append(tokens, s.makeToken(COMMA, ",", "null"))
		case ';':
			tokens = append(tokens, s.makeToken(SEMICOLON, ";", "null"))
		case '+':
			tokens = append(tokens, s.makeToken(PLUS, "+", "null"))
		case '-':
			tokens = append(tokens, s.makeToken(MINUS, "-", "null"))
		case '/':
			if s.match('/') {
				s.scanComment()
			} else {
				tokens = append(tokens, s.makeToken(SLASH, "/", "null"))
			}
		case '=':
			if s.match('=') {
				tokens = append(tokens, s.makeToken(EQUAL_EQUAL, "==", "null"))
			} else {
				tokens = append(tokens, s.makeToken(EQUAL, "=", "null"))
			}
		case '!':
			if s.match('=') {
				tokens = append(tokens, s.makeToken(BANG_EQUAL, "!=", "null"))
			} else {
				tokens = append(tokens, s.makeToken(BANG, "!", "null"))
			}
		case '<':
			if s.match('=') {
				tokens = append(tokens, s.makeToken(LESS_EQUAL, "<=", "null"))
			} else {
				tokens = append(tokens, s.makeToken(LESS, "<", "null"))
			}
		case '>':
			if s.match('=') {
				tokens = append(tokens, s.makeToken(GREATER_EQUAL, ">=", "null"))
			} else {
				tokens = append(tokens, s.makeToken(GREATER, ">", "null"))
			}
		case '"':
			token, ok := s.scanString()
			if !ok {
				break
				// return tokens
			}
			tokens = append(tokens, token)
		case '\n':
			s.Line++
		case '\r', ' ', '\t':
		default:
			if unicode.IsDigit(lexeme) {
				token := s.scanNumber(byte(lexeme))
				tokens = append(tokens, token)
			} else if unicode.IsLetter(lexeme) || lexeme == '_' {
				token := s.scanIdentifier(byte(lexeme))
				tokens = append(tokens, token)
			} else {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", s.Line, string(lexeme))
				s.Errors = true
			}
		}
	}
	tokens = append(tokens, s.makeToken(EOF, "", "null"))
	return tokens
}

func formatNumber(s string) string {
	if !strings.Contains(s, ".") {
		return s + ".0"
	}
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
	}
	if strings.HasSuffix(s, ".") {
		s += "0"
	}
	return s
}

func GetSource(filePath string) string {
	contentByte, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("error reading file: ", err)
	}
	return string(contentByte)
}
