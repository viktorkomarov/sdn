package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokenize(t *testing.T) {
	testCases := []struct {
		expr           string
		expectedTokens []TokenVal
	}{
		{
			expr: "A  & B",
			expectedTokens: []TokenVal{
				{
					Typ: Variable,
					Val: 'A',
				},
				{
					Typ: DoubleOp,
					Val: '&',
				},
				{
					Typ: Variable,
					Val: 'B',
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.expr, func(t *testing.T) {
			actual, err := tokenize(tC.expr)
			require.NoError(t, err, tC.expr)
			require.Equal(t, tC.expectedTokens, actual, tC.expr)
		})
	}
}
