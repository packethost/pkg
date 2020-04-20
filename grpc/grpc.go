// Copyright 2019 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"crypto/tls"
	"fmt"
	"net"
	"strconv"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/packethost/pkg/env"
	"github.com/packethost/pkg/log"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Server is used to hold configured info and ultimately the grpc server
type Server struct {
	err       error
	port      int
	server    *grpc.Server
	cert      tls.Certificate
	options   []grpc.ServerOption
	streamers []grpc.StreamServerInterceptor
	unariers  []grpc.UnaryServerInterceptor
	registry  []func(*grpc.Server)
}

// The ServiceRegister type is used as a callback once the underlying grpc server is setup to register the main service.
type ServiceRegister func(*Server)

// The Option type describes functions that operate on Server during NewServer.
// It is a convenience type to make it easier for callers to build up the slice of options apart from the call to NewServer.
type Option func(*Server)

// NewServer creates a new grpc server.
// By default the server will be an insecure server listening on port 8080 with logging and prometheus interceptors setup.
//
// The server's port is configured via the GRPC_PORT env variable, but can be overriden by the Port helper func.
// A tls server is setup if keys are provided in either the environment variables GRPC_CERT/GRPC_KEY, or using the X509KeyPair or LoadX509KeyPair helper funcs.
// Logging is always setup using the provided log.Logger.
// Prometheus is always setup using the default prom interceptors and Register func.
//
// req is called after the server has been setup.
// This is where your service is gets registered with grpc, equivalent to pb.RegisterMyServiceServer(s, &myServiceImpl{}).
//
// After your service has been registered any callbacks that were setup with Register will be called to finish up registration.
func NewServer(l log.Logger, reg ServiceRegister, options ...Option) (*Server, error) {
	s := &Server{}

	logStream, logUnary := l.GRPCLoggers()
	s.streamers = append(s.streamers, logStream, grpc_prometheus.StreamServerInterceptor)
	s.unariers = append(s.unariers, logUnary, grpc_prometheus.UnaryServerInterceptor)
	s.registry = append(s.registry, grpc_prometheus.Register)

	for _, opt := range options {
		opt(s)
		if s.err != nil {
			return nil, s.err
		}
	}

	if err := maybeSetPortFromEnv(s); err != nil {
		return nil, err
	}

	if err := maybeSetTLSFromEnv(s); err != nil {
		return nil, err
	}

	s.options = append(s.options,
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(s.streamers...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(s.unariers...)),
	)

	s.server = grpc.NewServer(s.options...)

	reg(s)

	for _, r := range s.registry {
		r(s.server)
	}

	return s, nil
}

// Port returns the port the server is listening on
func (s *Server) Port() int {
	return s.port
}

// Server returns the grpc server
func (s *Server) Server() *grpc.Server {
	return s.server
}

// Serve starts the grpc server
func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return errors.Wrap(err, "listen")
	}
	defer lis.Close()

	return errors.Wrap(s.server.Serve(lis), "serve")
}

// ServerOption will add the opt param to the underlying grpc.NewServer() call.
func ServerOption(opt grpc.ServerOption) Option {
	return func(s *Server) {
		s.options = append(s.options, opt)
	}
}

// maybeSetPortFromEnv will pick up the port from the environment, but only if the user hasn't specified the port via `Port`
func maybeSetPortFromEnv(s *Server) error {
	if s.port != 0 {
		return nil
	}

	port, err := strconv.Atoi(env.Get("GRPC_PORT", "8080"))
	if err != nil {
		return errors.Wrap(err, "parse grpc port from env")
	}
	if port < 1 {
		return errors.New("port must be > 1")
	}

	s.port = port
	return nil
}

// Port will set the port the server will bind to, Port must be > 0
func Port(port int) Option {
	return func(s *Server) {
		if port < 1 {
			s.err = errors.New("port must be > 1")
		}

		s.port = port
	}
}

// maybeSetTLSFromEnv will pick up the tls info from the environment, but only if the user hasn't specified the info via `X509KeyPair` or `LoadX509KeyPair`
func maybeSetTLSFromEnv(s *Server) error {
	var creds credentials.TransportCredentials

	if s.cert.PrivateKey == nil {
		cert := env.Get("GRPC_CERT")
		key := env.Get("GRPC_KEY")

		if cert != "" && key != "" {
			kp, err := tls.X509KeyPair([]byte(cert), []byte(key))
			if err != nil {
				return errors.Wrap(err, "parse tls files")
			}
			s.cert = kp
			creds = credentials.NewServerTLSFromCert(&kp)
		}
	} else {
		creds = credentials.NewServerTLSFromCert(&s.cert)
	}

	if creds != nil {
		s.options = append(s.options, grpc.Creds(creds))
	}
	return nil
}

// X509KeyPair will setup server as a secure server using the provided cert and key
// This function overrides GRPC_CERT and GRPC_KEY environment variables.
// NewServer will return an error if both X509KeyPair and LoadX509KeyPair are used.
func X509KeyPair(certPEMBlock, keyPEMBlock string) Option {
	return func(s *Server) {
		if s.cert.PrivateKey != nil {
			s.err = errors.New("certificate is already set")
			return
		}

		var err error
		s.cert, err = tls.X509KeyPair([]byte(certPEMBlock), []byte(keyPEMBlock))
		if err != nil {
			s.err = errors.Wrap(err, "parse x509 key pair")
			return
		}
	}
}

// LoadX509KeyPair will setup server as a secure server by reading the cert and key from the provided file locations.
// This function overrides GRPC_CERT and GRPC_KEY environment variables.
// NewServer will return an error if both X509KeyPair and LoadX509KeyPair are used.
func LoadX509KeyPair(certFile, keyFile string) Option {
	return func(s *Server) {
		if s.cert.PrivateKey != nil {
			s.err = errors.New("certificate is already set")
			return
		}

		var err error
		s.cert, err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			s.err = errors.Wrap(err, "load x509 key pair")
			return
		}
	}
}

// StreamInterceptor adds the argument to the list of interceptors in a grpc_middleware.Chain
// Logging and Prometheus interceptors are always included in the set
func StreamInterceptor(si grpc.StreamServerInterceptor) Option {
	return func(s *Server) {
		s.streamers = append(s.streamers, si)
	}
}

// UnaryInterceptor adds the argument to the list of interceptors in a grpc_middleware.Chain
// Logging and Prometheus interceptors are always included in the set
func UnaryInterceptor(ui grpc.UnaryServerInterceptor) Option {
	return func(s *Server) {
		s.unariers = append(s.unariers, ui)
	}
}

// Register will call the callback func after the main grpc service has been setup.
// The Prometheus register is always included in the set
func Register(r func(*grpc.Server)) Option {
	return func(s *Server) {
		s.registry = append(s.registry, r)
	}
}
