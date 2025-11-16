package parser_test

// This test only makes sure that things can be parsed. The output is not checked!

import (
	"fmt"
	"testing"

	"github.com/djpiper28/rpg-book/common/search/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parse(input string) error {
	return parser.Parse(input)
}

func TestBasicQuery(t *testing.T) {
	err := parse("John123")
	require.NoError(t, err)
}

func TestWhitespaceWrappedQuery(t *testing.T) {
	err := parse("    John123  \t\n")
	require.NoError(t, err)
}

func TestQueryThenQuery(t *testing.T) {
	err := parse("John Smith")
	require.NoError(t, err)
}

func TestAndOperator(t *testing.T) {
	err := parse("John and Smith")
	require.NoError(t, err)
}

func TestOrOperator(t *testing.T) {
	err := parse("John or Smith")
	require.NoError(t, err)
}

func TestXorOperator(t *testing.T) {
	err := parse("John xor Smith")
	require.NoError(t, err)
}

func TestBrackets(t *testing.T) {
	err := parse("(red and green)")
	require.NoError(t, err)
}

func TestNestedBrackets(t *testing.T) {
	err := parse("((red and green) or blue)")
	require.NoError(t, err)
}

func TestSetGenerator(t *testing.T) {
	err := parse("a:b")
	require.NoError(t, err)
}

func TestSetGeneratorQuotedParts(t *testing.T) {
	err := parse(`"test":"123"`)
	require.NoError(t, err)
}

func TestSetGenerators(t *testing.T) {
	operators := []string{
		"<", "<=", ">", ">=", "~", "=", ":",
	}

	for _, operator := range operators {
		err := parse(fmt.Sprintf("a%sb", operator))
		assert.NoError(t, err)
	}
}

func TestNegatedSetGenerator(t *testing.T) {
	err := parse("-a:b")
	require.NoError(t, err)
}

func TestBigQuery(t *testing.T) {
	err := parse("red and (green and blue) or test:yes xor (uwu or uwu:owo)")
	require.NoError(t, err)
}
