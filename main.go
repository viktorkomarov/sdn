package main

import (
	"fmt"
	"log"
	"math"
	"strings"
)

type Expression struct {
	variables map[rune]bool
	executor  Executor
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
		return !leftVal && rightVal
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
	return !r.Executor.Execute(variables)
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

func filterRows(rows []Row, f func(r Row) bool) []Row {
	n := 0
	for _, x := range rows {
		if f(x) {
			rows[n] = x
			n++
		}
	}
	return rows[:n]
}

func BuildPDNF(rows []Row) string {
	rows = filterRows(rows, func(r Row) bool { return r.result })

	var builder strings.Builder
	for i, row := range rows {
		for _var, val := range row.values {
			if !val {
				builder.WriteRune('!')
			}

			builder.WriteRune(_var)
		}

		if i != len(rows)-1 {
			builder.WriteString(" || ")
		}
	}

	return builder.String()
}

func main() {
	tokenizer, err := NewTokenizer("A > B")
	if err != nil {
		log.Fatalf("create tokenizer %v", err)
	}

	parser := Parser{tokenizer: tokenizer}
	exec, _, err := parser.parseExpr()
	if err != nil {
		log.Fatalf("parse expr %v", err)
	}

	rows := GenerateExamples(Expression{
		variables: tokenizer.Variables(),
		executor:  exec,
	})

	fmt.Printf("%s\n", BuildPDNF(rows))
}
