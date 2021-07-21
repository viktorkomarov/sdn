package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var ErrEndOfTokens = errors.New("end of tokens")

type Token int

const (
	Variable Token = iota
	SingleOp
	DoubleOp
	OpenBracket
	CloseBracket
	Empty
)

type TokenVal struct {
	Typ Token
	Val rune
}

type Tokenizer struct {
	tokens []TokenVal
	cur    int
}

func NewTokenizer(expr string) (Tokenizer, error) {
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
				return Tokenizer{}, fmt.Errorf("unknown token %s", string(symbol))
			}

			token = Variable
		}

		tokens = append(tokens, TokenVal{
			Typ: token,
			Val: symbol,
		})
	}

	return Tokenizer{
		tokens: tokens,
	}, nil
}

func (t Tokenizer) Next() (TokenVal, error) {
	if t.cur == len(t.tokens) {
		return TokenVal{}, ErrEndOfTokens
	}

	t.cur += 1
	return t.tokens[t.cur-1], nil
}
