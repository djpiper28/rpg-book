package parser_test

// This test only makes sure that things can be parsed. The output is not checked!

import (
	"go/parser"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicQuery(t *testing.T) {
	_, err := parser.ParseExpr("John123")
	require.NoError(t, err)
}

func TestWhitespaceWrappedQuery(t *testing.T) {
	_, err := parser.ParseExpr("    John123  \t\n")
	require.NoError(t, err)
}

func TestQueryThenQuery(t *testing.T) {
	_, err := parser.ParseExpr("John Smith")
	require.NoError(t, err)
}

func TestAndOperator(t *testing.T) {
	_, err := parser.ParseExpr("John and Smith")
	require.NoError(t, err)
}

func TestOrOperator(t *testing.T) {
	_, err := parser.ParseExpr("John or Smith")
	require.NoError(t, err)
}

func TestXorOperator(t *testing.T) {
	_, err := parser.ParseExpr("John xor Smith")
	require.NoError(t, err)
}
