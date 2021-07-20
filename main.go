package main

import (
	"errors"
	"fmt"
	"math"
)

// E = V | S E | (E) | V D E
// V = one of (A....Z)
// S =  !
// D =  & | '|' | > | - | +

type Expression struct {
	variables map[rune]bool
	executor  Executor
}

func parseDoubleOp(tokens []TokenVal) error {
	if len(tokens) == 0 {
		return ErrEmptyTokens
	}

	for _, val := range []rune{'&', '|', '>', '-', '+'} {
		if val == tokens[0].Val {
			return nil
		}
	}

	return fmt.Errorf("unknown double operation %s", string(tokens[0].Val))
}

var ErrEmptyTokens = errors.New("empty tokens")

func parseExpr(tokens []TokenVal) (Executor, error) {
	if len(tokens) == 0 {
		return nil, ErrEmptyTokens
	}

	switch tokens[0].Typ {
	case Variable:
		if err := parseDoubleOp(tokens[1:]); err != nil {
			if errors.Is(err, ErrEmptyTokens) {
				return Var{Name: tokens[0].Val}, nil
			}

			return nil, err
		}

		second, err := parseExpr(tokens[2:])
		if err != nil {
			return nil, err
		}

		return Expr{
			Left:  Var{Name: tokens[0].Val},
			Right: second,
			Op:    tokens[1],
		}, nil
	case SingleOp:
		exec, err := parseExpr(tokens[1:])
		if err != nil {
			return nil, err
		}
		return ReverseExecutor{exec}, nil
	case OpenBracket:
		exec, err := parseExpr(tokens[1:])
		if err != nil {
			return nil, err
		}

		return exec, nil
	case CloseBracket:
		return nil, nil // ???
	default:
		return nil, fmt.Errorf("unexpected toke %s", string(tokens[0].Val))
	}
}

func Parse(str string) (Expression, error) {
	tokens, err := tokenize(str)
	if err != nil {
		return Expression{}, err
	}

	variables := make(map[rune]bool)
	for _, token := range tokens {
		if token.Typ == Variable {
			variables[token.Val] = true
		}
	}

	expr := Expression{variables: variables}
	for {
		return expr, nil
	}
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
