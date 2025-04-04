package backend

type Server struct {
	Port                  int
	publicKey, privateKey []byte
}

func New() (*Server, error) {
	server := &Server{
		Port: 9881,
	}

	go server.start()

	return server, nil
}

func (s *Server) start() {
  
}

type ClientCredentials struct {
	Port         int    `json:"port"`
	PublicKeyB64 string `json:"publicKey"`
	Jwt          string `json:"jwt"` // Authenticates that the client is allowed to talk to RPG-book.
}
