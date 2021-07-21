package main

import (
	"fmt"
	"math"
)

type Expression struct {
	variables map[rune]bool
	executor  Executor
}

func parseDoubleOp(tokens []TokenVal) ([]TokenVal, error) {
	if len(tokens) == 0 {
		return nil, ErrEndOfTokens
	}

	for _, val := range []rune{'&', '|', '>', '-', '+'} {
		if val == tokens[0].Val {
			return tokens[1:], nil
		}
	}

	return nil, fmt.Errorf("unknown double operation %s", string(tokens[0].Val))
}

type Executor interface {
	Execute(variables map[rune]bool) bool
}

type Var struct {
	Name rune
}

func (v Var) Execute(variables map[rune]bool) bool {
	return variables[v.Name]
}

type Expr struct {
	Left  Executor
	Op    TokenVal
	Right Executor
}

func (e Expr) Execute(variables map[rune]bool) bool {
	leftVal := e.Left.Execute(variables)
	rightVal := e.Right.Execute(variables)

	switch e.Op.Val {
	case '&':
		return leftVal && rightVal
	case '|':
		return leftVal || rightVal
	case '>':
		return leftVal || !rightVal
	case '-':
		return (leftVal && rightVal) || (!rightVal && !leftVal)
	case '+':
		return !((leftVal && rightVal) || (!rightVal && !leftVal))
	}

	return false
}

type ReverseExecutor struct {
	Executor
}

func (r ReverseExecutor) Execute(variables map[rune]bool) bool {
	return !r.Execute(variables)
}

type Row struct {
	values map[rune]bool
	result bool
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
