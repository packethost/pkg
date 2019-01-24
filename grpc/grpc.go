package grpc

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

// Config represents the configuration fields required for running a gRPC service
type Config struct {
	Port     int
	BindHost string
}

// ConfigFromEnv will produce the config for a gRPC service from the standard environment variables used across Packet
func ConfigFromEnv() (*Config, error) {
	port := os.Getenv("GRPC_PORT")
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse GRPC_PORT into an int")
	}

	bindHost := os.Getenv("GRPC_BIND_HOST")

	return &Config{
		Port:     portInt,
		BindHost: bindHost,
	}, nil
}

// BindAddress constructs a bind address from the Port and BindHost in the Config
func (c *Config) BindAddress() (string, error) {
	return c.BindHost + ":" + strconv.Itoa(c.Port), nil
}
