package parser_test

import (
	"testing"

	"github.com/djpiper28/rpg-book/common/search/parser"
	"github.com/stretchr/testify/require"
)

func TestImpliedAnd(t *testing.T) {
	ast, err := parser.Parse("a b")
	require.NoError(t, err)
	require.Equal(t, &parser.Node{
		Type:           parser.NodeType_BinaryOperator,
		BinaryOperator: parser.BinaryOperator_And,
		Left: &parser.Node{
			Type: parser.NodeType_Basic,
			BasicQuery: parser.TextQuery{
				Value: "b",
			},
		},
		Right: &parser.Node{
			Type: parser.NodeType_Basic,
			BasicQuery: parser.TextQuery{
				Value: "a",
			},
		},
	}, ast)
}

func TestAnd(t *testing.T) {
	ast, err := parser.Parse("a and b")
	require.NoError(t, err)
	require.Equal(t, &parser.Node{
		Type:           parser.NodeType_BinaryOperator,
		BinaryOperator: parser.BinaryOperator_And,
		Left: &parser.Node{
			Type: parser.NodeType_Basic,
			BasicQuery: parser.TextQuery{
				Value: "b",
			},
		},
		Right: &parser.Node{
			Type: parser.NodeType_Basic,
			BasicQuery: parser.TextQuery{
				Value: "a",
			},
		},
	}, ast)
}

func TestOr(t *testing.T) {
	ast, err := parser.Parse("a or b")
	require.NoError(t, err)
	require.Equal(t, &parser.Node{
		Type:           parser.NodeType_BinaryOperator,
		BinaryOperator: parser.BinaryOperator_Or,
		Left: &parser.Node{
			Type: parser.NodeType_Basic,
			BasicQuery: parser.TextQuery{
				Value: "b",
			},
		},
		Right: &parser.Node{
			Type: parser.NodeType_Basic,
			BasicQuery: parser.TextQuery{
				Value: "a",
			},
		},
	}, ast)
}

func TestXor(t *testing.T) {
	ast, err := parser.Parse("a xor b")
	require.NoError(t, err)
	require.Equal(t, &parser.Node{
		Type:           parser.NodeType_BinaryOperator,
		BinaryOperator: parser.BinaryOperator_Xor,
		Left: &parser.Node{
			Type: parser.NodeType_Basic,
			BasicQuery: parser.TextQuery{
				Value: "b",
			},
		},
		Right: &parser.Node{
			Type: parser.NodeType_Basic,
			BasicQuery: parser.TextQuery{
				Value: "a",
			},
		},
	}, ast)
}

func TestMultipleOperators(t *testing.T) {
	ast, err := parser.Parse("a xor b and c")
	require.NoError(t, err)
	require.Equal(t, &parser.Node{
		Type:           parser.NodeType_BinaryOperator,
		BinaryOperator: parser.BinaryOperator_Xor,
		Left: &parser.Node{
			Type:           parser.NodeType_BinaryOperator,
			BinaryOperator: parser.BinaryOperator_And,
			Left: &parser.Node{
				Type: parser.NodeType_Basic,
				BasicQuery: parser.TextQuery{
					Value: "c",
				},
			},
			Right: &parser.Node{
				Type: parser.NodeType_Basic,
				BasicQuery: parser.TextQuery{
					Value: "b",
				},
			},
		},
		Right: &parser.Node{
			Type: parser.NodeType_Basic,
			BasicQuery: parser.TextQuery{
				Value: "a",
			},
		},
	}, ast)
}

func TestMultipleOperatorsBrackets1(t *testing.T) {
	ast, err := parser.Parse("a xor (b and c)")
	require.NoError(t, err)
	require.Equal(t, &parser.Node{
		Type:           parser.NodeType_BinaryOperator,
		BinaryOperator: parser.BinaryOperator_Xor,
		Left: &parser.Node{
			Type:           parser.NodeType_BinaryOperator,
			BinaryOperator: parser.BinaryOperator_And,
			Left: &parser.Node{
				Type: parser.NodeType_Basic,
				BasicQuery: parser.TextQuery{
					Value: "c",
				},
			},
			Right: &parser.Node{
				Type: parser.NodeType_Basic,
				BasicQuery: parser.TextQuery{
					Value: "b",
				},
			},
		},
		Right: &parser.Node{
			Type: parser.NodeType_Basic,
			BasicQuery: parser.TextQuery{
				Value: "a",
			},
		},
	}, ast)
}

func TestMultipleOperatorsBrackets2(t *testing.T) {
	ast, err := parser.Parse("(a xor (b and c))")
	require.NoError(t, err)
	require.Equal(t, &parser.Node{
		Type:           parser.NodeType_BinaryOperator,
		BinaryOperator: parser.BinaryOperator_Xor,
		Left: &parser.Node{
			Type:           parser.NodeType_BinaryOperator,
			BinaryOperator: parser.BinaryOperator_And,
			Left: &parser.Node{
				Type: parser.NodeType_Basic,
				BasicQuery: parser.TextQuery{
					Value: "c",
				},
			},
			Right: &parser.Node{
				Type: parser.NodeType_Basic,
				BasicQuery: parser.TextQuery{
					Value: "b",
				},
			},
		},
		Right: &parser.Node{
			Type: parser.NodeType_Basic,
			BasicQuery: parser.TextQuery{
				Value: "a",
			},
		},
	}, ast)
}
