package main

import (
	"fmt"
	"math"
	"strings"
	"unicode"
)

// E =  F | H | S H | S F
// H =  F D F | (H)
// F = one of (A....Z)
// S =  !
// D =  & | '|' | > | - | +

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

type Expression struct {
	variables map[rune]bool
	executor  Executor
}

// func Parse(expr string) (Expression, error) {
// 	tokens, err := tokenize(expr)
// 	if err != nil {
// 		return Expression{}, err
// 	}

// 	variables := make(map[rune]bool)
// 	for _, token := range tokens {
// 		if token.Typ == Variable {
// 			variables[token.Val] = true
// 		}
// 	}

// 	return Expression{
// 		variables: variables,
// 		tokens:    tokens,
// 	}, nil
// }

type Row struct {
	values map[rune]bool
	result bool
}

type Executor interface {
	Execute(variables map[rune]bool) bool
}

type Var struct {
	Name     rune
	OpBefore TokenVal
}

func (v Var) Execute(variables map[rune]bool) bool {
	if v.OpBefore.Typ == SingleOp && v.OpBefore.Val == '!' {
		return !variables[v.Name]
	}

	return variables[v.Name]
}

type Expr struct {
	OpBefore TokenVal
	Left     Executor
	Op       TokenVal
	Right    Executor
}

func (e Expr) Execute(variables map[rune]bool) bool {
	leftVal := e.Left.Execute(variables)
	rightVal := e.Right.Execute(variables)

	safety := true
	if e.OpBefore.Typ == SingleOp && e.OpBefore.Val == '!' {
		safety = false
	}

	val := false
	switch e.Op.Val {
	case '&':
		val = leftVal && rightVal
	case '|':
		val = leftVal || rightVal
	case '>':
		val = leftVal || !rightVal
	case '-':
		val = (leftVal && rightVal) || (!rightVal && !leftVal)
	case '+':
		val = !((leftVal && rightVal) || (!rightVal && !leftVal))
	}

	return safety && val
}

func GenerateExamples(expr Expression) []Row {
	count := int(math.Pow(float64(len(expr.variables)), 2.0))
	rows := make([]Row, 0, count)
	variables := make([]rune, 0, len(expr.variables))
	for k := range expr.variables {
		variables = append(variables, k)
	}

	for mask := 0; mask < count; mask++ {
		values := make(map[rune]bool)

		for i, r := range variables {
			if ((1 << i) & mask) == 1 {
				values[r] = true
			} else {
				values[r] = false
			}
		}

		rows = append(rows, Row{
			result: expr.executor.Execute(values),
			values: values,
		})
	}

	return rows
}
