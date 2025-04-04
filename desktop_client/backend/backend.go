package backend

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"net"
	"time"
)

type Server struct {
	ctx    context.Context
	Port   int
	cert   []byte
	cancel func()
}

func New() (*Server, error) {
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
		IPAddresses:        []net.IP{net.IPv4(127, 0, 0, 1)},
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
		Port:   9881,
		ctx:    ctx,
		cancel: cancel,
		cert:   cert,
	}

	go server.start()

	return server, nil
}

func (s *Server) start() {

}

func (s *Server) Stop() {
	s.cancel()
}

type ClientCredentials struct {
	Port         int    `json:"port"`
	PublicKeyB64 string `json:"publicKey"`
}
