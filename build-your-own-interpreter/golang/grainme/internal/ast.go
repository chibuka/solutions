package internal

import (
	"fmt"
)

type Expr interface {
	String() string
}

type Literal struct {
	// true, false, nil, a number, or a string
	value any
}

type Grouping struct {
	inner Expr
}

// TODO: not sure about this.
type Unary struct {
	operator Token
	operand  Expr
}

// TODO: not sure about this.
type Binary struct {
	operator Token
	left     Expr
	right    Expr
}

func (l Literal) String() string {
	switch v := l.value.(type) {
	case int, float64:
		vStr := fmt.Sprint(v)
		return formatNumber(vStr)
	case nil:
		return "nil"
	default:
		return fmt.Sprint(v)
	}
}

func (g Grouping) String() string {
	return "(group " + g.inner.String() + ")"
}

func (u Unary) String() string {
	return "(" + u.operator.Lexeme + " " + u.operand.String() + ")"
}

func (b Binary) String() string {
	return "(" + b.operator.Lexeme + " " + b.left.String() + " " + b.right.String() + ")"
}
