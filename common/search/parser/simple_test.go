package parser_test

// This test only makes sure that things can be parsed. The output is not checked!

import (
	"fmt"
	"testing"

	"github.com/djpiper28/rpg-book/common/search/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The error should be handled by the caller
func blindParse(t *testing.T, input string) error {
	t.Helper()
	output, err := parser.Parse(input)
	t.Logf("output: %+v", output)
	return err
}

func TestBasicQuery1(t *testing.T) {
	err := blindParse(t, "John123")
	require.NoError(t, err)
}

func TestBasicQuery2(t *testing.T) {
	err := blindParse(t, "3c148c94-9b37-4dce-ba50-c6a391af7391")
	require.NoError(t, err)
}

func TestBasicQuery3(t *testing.T) {
	err := blindParse(t, "µtest_123")
	require.NoError(t, err)
}

func TestBasicQuery4(t *testing.T) {
	err := blindParse(t, `"µtest_123"`)
	require.NoError(t, err)
}

func TestBasicQuery5(t *testing.T) {
	err := blindParse(t, `"µtest_123\ntest"`)
	require.NoError(t, err)
}

func TestWhitespaceWrappedQuery(t *testing.T) {
	err := blindParse(t, "    John123  \t\n")
	require.NoError(t, err)
}

func TestQueryThenQuery(t *testing.T) {
	err := blindParse(t, "John Smith")
	require.NoError(t, err)
}

func TestAndOperator(t *testing.T) {
	err := blindParse(t, "John and Smith")
	require.NoError(t, err)
}

func TestOrOperator(t *testing.T) {
	err := blindParse(t, "John or Smith")
	require.NoError(t, err)
}

func TestBrackets(t *testing.T) {
	err := blindParse(t, "(red and green)")
	require.NoError(t, err)
}

func TestNestedBrackets(t *testing.T) {
	err := blindParse(t, "((red and green) or blue)")
	require.NoError(t, err)
}

func TestSetGenerator(t *testing.T) {
	err := blindParse(t, "a:b")
	require.NoError(t, err)
}

func TestSetGeneratorQuotedParts(t *testing.T) {
	err := blindParse(t, `"test":"123"`)
	require.NoError(t, err)
}

func TestSetGenerators(t *testing.T) {
	operators := []string{
		"<", "<=", ">", ">=", "~", "=", ":",
	}

	for _, operator := range operators {
		err := blindParse(t, fmt.Sprintf("a%sb", operator))
		assert.NoError(t, err)
	}
}

func TestNegatedSetGenerator(t *testing.T) {
	err := blindParse(t, "-a:b")
	require.NoError(t, err)
}

func TestBigQuery1(t *testing.T) {
	err := blindParse(t, "red and (green and blue) or test:yes or (uwu or uwu:owo)")
	require.NoError(t, err)
}

func TestBigQuery2(t *testing.T) {
	err := blindParse(t, `lyrics:"what is love?" and (isGood:true or isClassic:true) and lorem:"ipsum dolor sit amet"`)
	require.NoError(t, err)
}

func TestEmptyQuery(t *testing.T) {
	err := blindParse(t, "")
	require.Error(t, err)
}

func TestInvalidQuery(t *testing.T) {
	err := blindParse(t, "test:")
	require.Error(t, err)
}

func TestInvalidQuotation1(t *testing.T) {
	err := blindParse(t, `"testing 123\"`)
	require.Error(t, err)
}

func TestInvalidQuotation2(t *testing.T) {
	err := blindParse(t, `test:"testing 123\"`)
	require.Error(t, err)
}
