package main

import (
	"fmt"
	"strings"
	"unicode"
)

type Token int

const (
	Variable Token = iota
	SingleOp
	DoubleOp
	OpenBracket
	CloseBracket
)

type TokenVal struct {
	Typ Token
	Val rune
}

func tokenize(expr string) ([]TokenVal, error) {
	tokens := make([]TokenVal, 0)

	for _, symbol := range strings.ReplaceAll(expr, " ", "") {
		var token Token
		switch symbol {
		case '(':
			token = OpenBracket
		case ')':
			token = CloseBracket
		case '!':
			token = SingleOp
		case '&', '|', '>', '-', '+':
			token = DoubleOp
		default:
			if !unicode.IsUpper(symbol) || !unicode.IsLetter(symbol) {
				return nil, fmt.Errorf("unknown token %s", string(symbol))
			}

			token = Variable
		}

		tokens = append(tokens, TokenVal{
			Typ: token,
			Val: symbol,
		})
	}

	return tokens, nil
}
