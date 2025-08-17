package authentication

type Alg string

type Config struct {
	Iss        string
	Alg        Alg
	SigningKey string
	JwksURL    string
}

// Supported Algorithms.
const (
	AlgRS256 = "RS256"
	AlgHS256 = "HS256"
)
