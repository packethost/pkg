package grpc

import (
	"crypto/tls"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials"
)

// Config represents the configuration fields required for running a gRPC service
type Config struct {
	Bind              string
	Port              int
	ServerCredentials credentials.TransportCredentials
}

// ConfigFromEnv will produce the config for a gRPC service from the standard environment variables used across Packet
func ConfigFromEnv() (*Config, error) {
	port := os.Getenv("GRPC_PORT")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse GRPC_PORT into an int")
	}

	bind := os.Getenv("GRPC_BIND")

	certStr := os.Getenv("GRPC_SERVER_CERT")
	keyStr := os.Getenv("GRPC_SERVER_KEY")

	var cred credentials.TransportCredentials
	if certStr != "" && keyStr != "" {
		cert, err := tls.X509KeyPair([]byte(certStr), []byte(keyStr))
		if err != nil {
			return nil, errors.Wrap(err, "could not load certificate")
		}
		cred = credentials.NewServerTLSFromCert(&cert)
	}

	return &Config{
		Port:              portInt,
		Bind:              bind,
		ServerCredentials: cred,
	}, nil
}

// BindAddress constructs a bind address from the Port and BindHost in the Config
func (c *Config) BindAddress() (string, error) {
	if c == nil {
		return "", errors.New("nil config")
	}
	return c.Bind + ":" + strconv.Itoa(c.Port), nil
}
