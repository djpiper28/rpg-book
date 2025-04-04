package backend_test

import (
	"encoding/pem"
	"testing"

	"github.com/djpiper28/rpg-book/desktop_client/backend"
	"github.com/stretchr/testify/require"
)

func TestStartStopBackend(t *testing.T) {
	server, err := backend.New(backend.RandPort())
	require.NoError(t, err)
	defer server.Stop()
}

func TestServerClientCredentialsProducesValidCert(t *testing.T) {
	port := backend.RandPort()

	server, err := backend.New(port)
	require.NoError(t, err)
	defer server.Stop()

	credentials := server.ClientCredentials()
	require.Equal(t, port, credentials.Port)

	_, rest := pem.Decode([]byte(credentials.PublicKeyB64))
	require.Len(t, rest, 0)
}
