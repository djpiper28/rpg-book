package backend_test

import (
	"crypto/x509"
	"encoding/base64"
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

	certDer, err := base64.StdEncoding.DecodeString(credentials.CertPem)
	require.NoError(t, err)
	x509.ParseCertificate(certDer)
}
