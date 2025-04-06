package backend_test

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"

	testutils "github.com/djpiper28/rpg-book/common/test_utils"
	"github.com/djpiper28/rpg-book/desktop_client/backend"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_common"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_settings"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func TestStartStopBackend(t *testing.T) {
	t.Parallel()
	server, err := backend.New(backend.RandPort())
	require.NoError(t, err)
	defer server.Stop()
}

func TestServerClientCredentialsProducesValidCert(t *testing.T) {
	t.Parallel()
	port := backend.RandPort()

	server, err := backend.New(port)
	require.NoError(t, err)
	defer server.Stop()

	credentials := server.ClientCredentials()
	require.Equal(t, port, credentials.Port)

	block, rest := pem.Decode([]byte(credentials.CertPem))
	require.Len(t, rest, 0)
	require.NotNil(t, block)
}

func TestServerNoTlsGet(t *testing.T) {
	t.Parallel()
	port := backend.RandPort()

	server, err := backend.New(port)
	require.NoError(t, err)
	defer server.Stop()

	grpcClient, err := grpc.NewClient(fmt.Sprintf("http://localhost:%d", server.ClientCredentials().Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer grpcClient.Close()

	client := pb_settings.NewSettingsSvcClient(grpcClient)
	ctx, cancel := testutils.NewTestContext()
	defer cancel()

	_, err = client.GetSettings(ctx, &pb_common.Empty{})
	require.Error(t, err)
}

func TestServerTlsGet(t *testing.T) {
	t.Parallel()
	port := backend.RandPort()

	server, err := backend.New(port)
	require.NoError(t, err)
	defer server.Stop()

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM([]byte(server.ClientCredentials().CertPem))
	require.True(t, ok)

	grpcClient, err := grpc.NewClient(
		fmt.Sprintf("https://127.0.0.1:%d", server.ClientCredentials().Port),
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(certPool, "")))
	require.NoError(t, err)
	defer grpcClient.Close()

	client := pb_settings.NewSettingsSvcClient(grpcClient)
	ctx, cancel := testutils.NewTestContext()
	defer cancel()

	settings, err := client.GetSettings(ctx, &pb_common.Empty{})
	require.NoError(t, err)
	require.NotNil(t, settings)
}
