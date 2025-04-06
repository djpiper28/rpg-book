package backend

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	mathrand "math/rand/v2"
	"net"
	"time"

	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_settings"
	settingssvc "github.com/djpiper28/rpg-book/desktop_client/backend/svc/settings_svc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	ctx      context.Context
	Port     int
	cert     []byte
	certKeys *rsa.PrivateKey
	cancel   func()
	server   *grpc.Server
}

func RandPort() int {
	return 9000 + mathrand.IntN(1000)
}

var localhost = net.IPv4(127, 0, 0, 1)

func New(port int) (*Server, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot create RSA key for gRPC"), err)
	}

	template := &x509.Certificate{
		Issuer: pkix.Name{
			Organization:       []string{"djpiper28"},
			OrganizationalUnit: []string{"RPG Book"},
			Province:           []string{"Reading"},
			Locality:           []string{},
			Country:            []string{"United Kingdom"},
		},
		PublicKeyAlgorithm: x509.RSA,
		SignatureAlgorithm: x509.SHA512WithRSA,
		NotBefore:          time.Now(),
		NotAfter:           time.Now().Add(time.Hour * 24 * 365),
		KeyUsage:           x509.KeyUsage(x509.ExtKeyUsageAny),
		EmailAddresses:     []string{"djpiper28@gmail.com"},
		IsCA:               true,
		IPAddresses:        []net.IP{localhost},
	}
	cert, err := x509.CreateCertificate(rand.Reader,
		template, template,
		key.Public(),
		key)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot create x509 certificate for gRPC"), err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	server := &Server{
		Port:     port,
		ctx:      ctx,
		cancel:   cancel,
		cert:     cert,
		certKeys: key,
	}

	err = server.start()
	return server, err
}

func (s *Server) loadTLSCredentials() (credentials.TransportCredentials, error) {
	keyPem := x509.MarshalPKCS1PrivateKey(s.certKeys)
	certPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: s.cert,
		},
	))

	serverCert, err := tls.X509KeyPair([]byte(certPem), []byte(keyPem))
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func (s *Server) start() error {
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{IP: localhost, Port: s.Port})
	if err != nil {
		return err
	}

	s.server = grpc.NewServer()
	s.server.RegisterService(&pb_settings.SettingsSvc_ServiceDesc, settingssvc.New())
	go s.server.Serve(listener)
	return nil
}

func (s *Server) Stop() {
	s.cancel()
	s.server.GracefulStop()
}

type ClientCredentials struct {
	Port       int    `json:"port"`
	CertPem string `json:"cert"`
}

func (s *Server) ClientCredentials() *ClientCredentials {
	return &ClientCredentials{
		Port:       s.Port,
		CertPem: base64.StdEncoding.EncodeToString(s.cert),
	}
}
