package internal

type Expr interface {
	printExpr()
}

type Literal struct {
	value any
}

func (l *Literal) prinExpr() {
	// TODO
}

func Parse(s Scanner) {
	// TODO
}
