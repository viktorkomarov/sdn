package main

import (
	"errors"
	"fmt"
)

type Parser struct {
	tokenizer Tokenizer
}

// E = V | S E | (E) | V D E
// V = one of (A....Z)
// S =  !
// D =  & | '|' | > | - | +

func (p Parser) parseDoubleOp() (TokenVal, error) {
	token, err := p.tokenizer.Next()
	if err != nil {
		return TokenVal{Typ: Empty}, ErrEndOfTokens
	}

	for _, val := range []rune{'&', '|', '>', '-', '+'} {
		if val == token.Val {
			return token, nil
		}
	}

	return TokenVal{}, fmt.Errorf("unknown double operation %s", string(token.Val))
}

func (p Parser) parseExpr() (Executor, Token, error) {
	token, err := p.tokenizer.Next()
	if err != nil {
		return nil, Empty, ErrEndOfTokens
	}

	switch token.Typ {
	case Variable:
		op, err := p.parseDoubleOp()
		if err != nil {
			if errors.Is(err, ErrEndOfTokens) {
				return Var{Name: token.Val}, token.Typ, nil
			}

			return nil, Empty, err
		}

		second, lastToken, err := p.parseExpr()
		if err != nil {
			return nil, Empty, err
		}

		return Expr{
			Left:  Var{Name: token.Val},
			Right: second,
			Op:    op,
		}, lastToken, nil
	case SingleOp:
		exec, _, err := p.parseExpr()
		if err != nil {
			return nil, Empty, err
		}

		return ReverseExecutor{exec}, token.Typ, nil
	case OpenBracket:
		exec, lastToken, err := p.parseExpr()
		if err != nil {
			return nil, Empty, err
		}

		if lastToken != CloseBracket {
			return nil, lastToken, errors.New("unexpected end of expr )")
		}
		return exec, lastToken, nil
	case CloseBracket:
		return nil, CloseBracket, nil
	default:
		return nil, Empty, fmt.Errorf("unexpected toke %s", string(token.Val))
	}

}
