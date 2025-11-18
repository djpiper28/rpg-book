package parser_test

import (
	"testing"

	"github.com/djpiper28/rpg-book/common/search/parser"
	"github.com/stretchr/testify/require"
)

func TestBasicQuery(t *testing.T) {
  const testString = "Testing-123"
	ast, err := parser.Parse(testString)
	require.NoError(t, err)
	require.Equal(t, parser.NodeType_Basic, ast.Type)
	require.Equal(t, testString, ast.BasicQuery.Value)
}

func TestBasicQueryQuotation(t *testing.T) {
	ast, err := parser.Parse(`"testing 123\ntest"`)
	require.NoError(t, err)
	require.Equal(t, parser.NodeType_Basic, ast.Type)
	require.Equal(t, `testing 123
test`, ast.BasicQuery.Value)
}
