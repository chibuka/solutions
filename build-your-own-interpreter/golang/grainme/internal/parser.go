/*
 * notes for me.
 *
 * for a new language i'll need to decide on precedence.
 * 1 + 2 .. 3 + 4 => (1 + 2) .. (3 + 4) or 1 + (2 .. 3) + 4
 *
 */
package internal

import (
	"errors"
	"log"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParse(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() (Expr, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.comparison()
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(LESS, LESS_EQUAL, GREATER, GREATER_EQUAL, EQUAL_EQUAL, BANG_EQUAL) {
		op := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = Binary{op, expr, right}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(PLUS, MINUS) {
		op := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = Binary{op, expr, right}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(STAR, SLASH) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = Binary{op, expr, right}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return Unary{op, right}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(NIL) {
		return Literal{nil}, nil
	}
	if p.match(TRUE) {
		return Literal{true}, nil
	}
	if p.match(FALSE) {
		return Literal{false}, nil
	}
	if p.match(NUMBER, STRING) {
		return Literal{p.previous().Literal}, nil
	}
	if p.match(LEFT_PAREN) {
		innerExpr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if p.check(RIGHT_PAREN) {
			p.advance()
		} else {
			return nil, errors.New("expect ')' after expression")
		}
		return Grouping{innerExpr}, nil
	}
	log.Fatalf("expect primary, got: %v", p.peek())

	return nil, errors.New("parsing failed")
}

func (p *Parser) match(ttypes ...TokenType) bool {
	for _, t := range ttypes {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) advance() Token {
	// we're not checking EOF because so far
	// all advance() calls are under check() condition
	p.current++
	return p.previous()
}

func (p Parser) check(ttype TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == ttype
}

func (p Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p Parser) peek() Token {
	return p.tokens[p.current]
}

func (p Parser) previous() Token {
	if p.current <= 0 {
		log.Fatalf("can't get token at index %d", p.current)
	}
	return p.tokens[p.current-1]
}
