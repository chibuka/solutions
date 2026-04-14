package internal

type TokenType int

const (
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	STAR
	DOT
	COMMA
	SEMICOLON
	PLUS
	MINUS
	SLASH
	EQUAL
	EQUAL_EQUAL
	BANG
	BANG_EQUAL
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL
	IDENTIFIER
	STRING
	NUMBER
	AND
	CLASS
	ELSE
	FALSE
	FOR
	FUN
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
	EOF
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal string
	Line    int
}

func (t TokenType) String() string {
	switch t {
	case LEFT_PAREN:
		return "LEFT_PAREN"
	case RIGHT_PAREN:
		return "RIGHT_PAREN"
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case STAR:
		return "STAR"
	case DOT:
		return "DOT"
	case COMMA:
		return "COMMA"
	case SEMICOLON:
		return "SEMICOLON"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case SLASH:
		return "SLASH"
	case EQUAL:
		return "EQUAL"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case BANG:
		return "BANG"
	case BANG_EQUAL:
		return "BANG_EQUAL"
	case LESS:
		return "LESS"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case GREATER:
		return "GREATER"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case AND:
		return "AND"
	case CLASS:
		return "CLASS"
	case ELSE:
		return "ELSE"
	case FALSE:
		return "FALSE"
	case FOR:
		return "FOR"
	case FUN:
		return "FUN"
	case IF:
		return "IF"
	case NIL:
		return "NIL"
	case OR:
		return "OR"
	case PRINT:
		return "PRINT"
	case RETURN:
		return "RETURN"
	case SUPER:
		return "SUPER"
	case THIS:
		return "THIS"
	case TRUE:
		return "TRUE"
	case VAR:
		return "VAR"
	case WHILE:
		return "WHILE"
	case EOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}
