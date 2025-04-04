package backend_test

import (
	"testing"

	"github.com/djpiper28/rpg-book/desktop_client/backend"
	"github.com/stretchr/testify/require"
)

func TestStartStopBackend(t *testing.T) {
	server, err := backend.New()
	require.NoError(t, err)
	server.Stop()
}
