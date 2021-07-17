package main

import (
	"errors"
	"math"
)

// E =  F | H | S H | S F
// H =  F D F | (H)
// F = one of (A....Z)
// S =  !
// D =  & | '|' | > | - | +

type Expression struct {
	variables map[rune]bool
	executor  Executor
}

var ErrEmptyTokens = errors.New("empty tokens")

func parseExpr(tokens []TokenVal) (Expr, error) {
	if len(tokens) == 0 {
		return Expr{}, ErrEmptyTokens
	}
}

func Parse(str string) (Expression, error) {
	tokens := tokenize(str)
	variables := make(map[rune]bool)
	for _, token := range tokens {
		if token.Typ == Variable {
			variables[token.Val] = true
		}
	}

	parseExpr(tokens)
}

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
