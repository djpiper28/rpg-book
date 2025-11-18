package normalisation_test

import (
	"testing"

	"github.com/djpiper28/rpg-book/common/normalisation"
	"github.com/stretchr/testify/require"
)

func TestBaseCase(t *testing.T) {
	const text = "johnny went to the market"
	require.Equal(t, text, normalisation.Normalise(text))
}

func TestIsLowerCase(t *testing.T) {
	require.Equal(t, "test-123", normalisation.Normalise("TeSt-123"))
}

func TestSpecialIsLowerCase(t *testing.T) {
	require.Equal(t, "francais", normalisation.Normalise("FRANÇAIS"))
}

func TestNonPrintable(t *testing.T) {
	require.Equal(t, "test-oooo", normalisation.Normalise("TEST-ÒÓÕÖ"))
}
