package parser_test

import (
	"testing"

	"github.com/djpiper28/rpg-book/common/search/parser"
	"github.com/stretchr/testify/require"
)

func TestUnescapeBaseCase(t *testing.T) {
	const baseCaseText = "This is a test 123"
	output, err := parser.UnescapeText(baseCaseText)
	require.NoError(t, err)
	require.Equal(t, baseCaseText, output)
}

func TestUnescapeUnknownChar(t *testing.T) {
	_, err := parser.UnescapeText("\\d")
	require.Error(t, err)
}

func TestUnescape1(t *testing.T) {
	output, err := parser.UnescapeText(`this is a \"test\" of the \'system\'`)
	require.NoError(t, err)
	require.Equal(t, `this is a "test" of the 'system'`, output)
}

func TestUnescape2(t *testing.T) {
	output, err := parser.UnescapeText(`This is a test \\`)
	require.NoError(t, err)
	require.Equal(t, `This is a test \`, output)
}

func TestUnescape3(t *testing.T) {
	output, err := parser.UnescapeText(`this is a \"test\" of the \'system\'\r\n\tuwu\\`)
	require.NoError(t, err)
	require.Equal(t, `this is a "test" of the 'system'`+"\r\n\tuwu\\", output)
}

func TestUnescape4(t *testing.T) {
	output, err := parser.UnescapeText(`this\nis\na\ntest`)
	require.NoError(t, err)
	require.Equal(t, `this
is
a
test`, output)
}

func TestUnescapeError(t *testing.T) {
	_, err := parser.UnescapeText("this is a test \\")
	require.Error(t, err)
}
