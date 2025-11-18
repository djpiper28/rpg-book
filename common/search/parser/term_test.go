package parser_test

import (
	"testing"

	"github.com/djpiper28/rpg-book/common/normalisation"
	"github.com/djpiper28/rpg-book/common/search/parser"
	"github.com/stretchr/testify/require"
)

func TestBasicQueryNonSpecial(t *testing.T) {
	testString := normalisation.Normalise("testing-123")
	ast, err := parser.Parse(testString)
	require.NoError(t, err)
	require.Equal(t, parser.NodeType_Basic, ast.Type)
	require.Equal(t, testString, ast.BasicQuery.Value)
}

func TestBasicQuerySpecial(t *testing.T) {
	testString := normalisation.Normalise("Testing-123µ")
	ast, err := parser.Parse(testString)
	require.NoError(t, err)
	require.Equal(t, parser.NodeType_Basic, ast.Type)
	require.Equal(t, testString, ast.BasicQuery.Value)
}

func TestBasicQueryQuotation(t *testing.T) {
	ast, err := parser.Parse(`"testing 123\ntest"`)
	require.NoError(t, err)
	require.Equal(t, parser.NodeType_Basic, ast.Type)
	require.Equal(t, "testing 123test", ast.BasicQuery.Value)
}

func TestBasicQueryQuotationSpecial(t *testing.T) {
	ast, err := parser.Parse(`"ÒÓÕÖTEsting 123\ntest"`)
	require.NoError(t, err)
	require.Equal(t, parser.NodeType_Basic, ast.Type)
	require.Equal(t, "ooootesting 123test", ast.BasicQuery.Value)
}

func TestMixedSetGenerator(t *testing.T) {
	ast, err := parser.Parse(`"test":value-123`)
	require.NoError(t, err)
	require.Equal(t, parser.NodeType_SetGenerator, ast.Type)
	require.Equal(t, "test", ast.SetGenerator.Key)
	require.Equal(t, "value-123", ast.SetGenerator.Value)
	require.False(t, ast.SetGenerator.Negated)
	require.Equal(t, parser.GeneratorOperator_Includes, ast.SetGenerator.GeneratorOperator)
}

func TestNegatedMixedSetGenerator(t *testing.T) {
	ast, err := parser.Parse(`-power>="5.3"`)
	require.NoError(t, err)
	require.Equal(t, parser.NodeType_SetGenerator, ast.Type)
	require.Equal(t, "power", ast.SetGenerator.Key)
	require.Equal(t, "5.3", ast.SetGenerator.Value)
	require.True(t, ast.SetGenerator.Negated)
	require.Equal(t, parser.GeneratorOperator_GreaterThanEquals, ast.SetGenerator.GeneratorOperator)
}

func TestBracketedExpression(t *testing.T) {
	ast, err := parser.Parse(`("test":value-123)`)
	require.NoError(t, err)
	require.Equal(t, parser.NodeType_SetGenerator, ast.Type)
	require.Equal(t, "test", ast.SetGenerator.Key)
	require.Equal(t, "value-123", ast.SetGenerator.Value)
	require.False(t, ast.SetGenerator.Negated)
	require.Equal(t, parser.GeneratorOperator_Includes, ast.SetGenerator.GeneratorOperator)
}
