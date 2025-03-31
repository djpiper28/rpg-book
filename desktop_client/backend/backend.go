package backend

type Server struct {
	Port                  int
	publicKey, privateKey []byte
}

func New() (*Server, error) {
	server := &Server{
		Port: 9881,
	}

	return server, nil
}

type ClientCredentials struct {
	Port         int    `json:"port"`
	PublicKeyB64 string `json:"publicKey"`
	Jwt          string `json:"jwt"` // Authenticates that the client is allowed to talk to RPG-book.
}
