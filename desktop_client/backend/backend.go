package backend

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	mathrand "math/rand/v2"
	"net"
	"net/url"
	"time"

	"github.com/charmbracelet/log"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
	"github.com/djpiper28/rpg-book/desktop_client/backend/database"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_system"
	systemsvc "github.com/djpiper28/rpg-book/desktop_client/backend/svc/system_svc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	ctx       context.Context
	Port      int
	cert      []byte
	certKeys  *rsa.PrivateKey
	cancel    func()
	server    *grpc.Server
	primaryDb *database.Db
}

func RandPort() int {
	return 9000 + mathrand.IntN(1000)
}

var localhost = net.IPv4(127, 0, 0, 1)

const primaryDb = "rpg-book-primary" + database.DbExtension

func New(port int) (*Server, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot create RSA key for gRPC"), err)
	}

	hostedUrl, err := url.Parse(fmt.Sprintf("https://127.0.0.1:%d", port))
	if err != nil {
		return nil, errors.Join(errors.New("Cannot create URL for cert"), err)
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
		URIs:               []*url.URL{hostedUrl},
		Version:            1,
	}
	cert, err := x509.CreateCertificate(rand.Reader,
		template, template,
		key.Public(),
		key)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot create x509 certificate for gRPC"), err)
	}

	db, err := database.New(primaryDb)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Cannot create primaryDb (%s)", primaryDb), err)
	}

	err = Migrate(db)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot perform default migrations"), err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	server := &Server{
		Port:      port,
		ctx:       ctx,
		cancel:    cancel,
		cert:      cert,
		certKeys:  key,
		primaryDb: db,
	}

	err = server.start()
	return server, err
}

func (s *Server) loadTLSCredentials() (credentials.TransportCredentials, error) {
	keyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(s.certKeys),
		},
	)
	certPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: s.cert,
		},
	)

	serverCert, err := tls.X509KeyPair([]byte(certPem), []byte(keyPem))
	if err != nil {
		return nil, errors.Join(errors.New("Cannot parse generated x509 cert"), err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

const (
	k              = 1024
	m              = k * k
	g              = m * k
	maxMessageSize = g
	maxHeaderSize  = 3 * m
)

func (s *Server) start() error {
	listener, err := net.ListenTCP("tcp4", &net.TCPAddr{IP: localhost, Port: s.Port})
	if err != nil {
		return errors.Join(errors.New("Cannot bind to port"), err)
	}

	creds, err := s.loadTLSCredentials()
	if err != nil {
		return errors.Join(errors.New("Cannot load TLS certificates"), err)
	}

	s.server = grpc.NewServer(
		grpc.Creds(creds),
		grpc.MaxSendMsgSize(maxMessageSize),
		grpc.MaxRecvMsgSize(maxMessageSize),
		grpc.MaxHeaderListSize(maxHeaderSize),
	)
	s.server.RegisterService(&pb_system.SystemSvc_ServiceDesc, systemsvc.New(s.primaryDb))
	go func() {
		defer listener.Close()
		err := s.server.Serve(listener)
		if err != nil {
			log.Info("Server died", loggertags.TagError, err)
		}
	}()
	return nil
}

func (s *Server) Stop() {
	defer s.server.GracefulStop()
	defer s.primaryDb.Close()
	s.cancel()
}

type ClientCredentials struct {
	Port    int    `json:"port"`
	CertPem string `json:"cert"`
}

func (s *Server) ClientCredentials() *ClientCredentials {
	return &ClientCredentials{
		Port: s.Port,
		CertPem: string(pem.EncodeToMemory(
			&pem.Block{
				Type:  "CERTIFICATE",
				Bytes: s.cert,
			},
		)),
	}
}
