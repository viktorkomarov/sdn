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
		return TokenVal{}, err
	}

	for _, val := range []rune{'&', '|', '>', '-', '+'} {
		if val == token.Val {
			return token, nil
		}
	}

	return TokenVal{}, fmt.Errorf("unknown double operation %s", string(token.Val))
}

func (p Parser) parseExpr() (Executor, error) {
	token, err := p.tokenizer.Next()
	if err != nil {
		return nil, err
	}

	switch token.Typ {
	case Variable:
		op, err := p.parseDoubleOp()
		if err != nil {
			if errors.Is(err, ErrEndOfTokens) {
				return Var{Name: token.Val}, nil
			}

			return nil, err
		}

		second, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		return Expr{
			Left:  Var{Name: token.Val},
			Right: second,
			Op:    op,
		}, nil
	case SingleOp:
		exec, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		return ReverseExecutor{exec}, nil
	case OpenBracket:
		exec, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		return exec, nil
	case CloseBracket:
		return nil, nil // ???
	default:
		return nil, fmt.Errorf("unexpected toke %s", string(token.Val))
	}

}
