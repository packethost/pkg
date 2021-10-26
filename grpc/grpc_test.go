// Copyright 2019 - 2020, Packethost, Inc and contributors
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/packethost/pkg/internal/testenv"
	"github.com/packethost/pkg/log"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const svc = "github.com/packethost/pkg/grpc"

func TestDefaultsAndOrdering(t *testing.T) {
	l := log.Test(t, svc)
	assert := require.New(t)

	index := 0

	var pre1I int
	pre1 := func(*Server) {
		index++
		pre1I = index
	}
	var pre2I int
	pre2 := func(*Server) {
		index++
		pre2I = index
	}
	var srvI int
	srv := func(*Server) {
		index++
		srvI = index
	}
	var reg1I int
	reg1 := func(*grpc.Server) {
		index++
		reg1I = index
	}
	var reg2I int
	reg2 := func(*grpc.Server) {
		index++
		reg2I = index
	}

	s, err := NewServer(l, srv, pre1, Register(reg1), pre2, Register(reg2))
	assert.NoError(err)
	assert.NotNil(s)
	assert.Equal(s.port, 8080)

	// ensure order
	assert.True(pre1I > 0)
	assert.True(pre2I > pre1I)
	assert.True(srvI > pre2I)
	assert.True(reg1I > srvI)
	assert.True(reg2I > reg1I)
}

func TestPort(t *testing.T) {
	defer testenv.Clear().Restore()

	l := log.Test(t, svc)
	assert := require.New(t)

	s, err := NewServer(l, defSrv)
	assert.NoError(err)
	assert.NotNil(s)
	assert.Equal(s.port, 8080)
	assert.Equal(s.Port(), 8080)

	os.Setenv("GRPC_PORT", "4242")
	defer os.Unsetenv("GRPC_PORT")

	s, err = NewServer(l, defSrv)
	assert.NoError(err)
	assert.NotNil(s)
	assert.Equal(s.port, 4242)
	assert.Equal(s.Port(), 4242)

	s, err = NewServer(l, defSrv, Port(2424))
	assert.NoError(err)
	assert.NotNil(s)
	assert.Equal(s.port, 2424)
	assert.Equal(s.Port(), 2424)

	os.Setenv("GRPC_PORT", "0")
	s, err = NewServer(l, defSrv)
	assert.Error(err)
	assert.Nil(s)

	os.Setenv("GRPC_PORT", "-1")
	s, err = NewServer(l, defSrv)
	assert.Error(err)
	assert.Nil(s)

	s, err = NewServer(l, defSrv, Port(0))
	assert.Error(err)
	assert.Nil(s)

	s, err = NewServer(l, defSrv, Port(-1))
	assert.Error(err)
	assert.Nil(s)
}

func genCert(t *testing.T) (string, string) {
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		DNSNames:  []string{"localhost"},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 180),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		t.Fatal(err)
	}
	out := &bytes.Buffer{}
	err = pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		t.Fatal(err)
	}
	cert := out.String()
	out.Reset()

	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	err = pem.Encode(out, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
	if err != nil {
		t.Fatal(err)
	}
	key := out.String()

	return cert, key
}

func writeTLSFiles(t *testing.T, cert, key string) (string, string) {
	f, err := ioutil.TempFile("", "pkg-grpc-testing-cert-*.pem")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	c := f.Name()
	if err = ioutil.WriteFile(f.Name(), []byte(cert), 0); err != nil {
		t.Fatal(err)
	}

	f, err = ioutil.TempFile("", "pkg-grpc-testing-key-*.pem")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	k := f.Name()
	if err = ioutil.WriteFile(f.Name(), []byte(key), 0); err != nil {
		t.Fatal(err)
	}

	return c, k
}

func serve(t *testing.T, s *Server, test func()) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		wg.Done()
		if err := s.Serve(); err != nil {
			panic(err)
		}
		wg.Done()
	}()

	wg.Wait()
	wg.Add(1)
	test()

	s.server.Stop()
	wg.Wait()
}

func connectGRPC(t *testing.T, port int, cert string) error {
	address := fmt.Sprintf("localhost:%d", port)

	creds := grpc.WithInsecure()
	if cert != "" {
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM([]byte(cert)) {
			t.Fatal("failed to add cert to pool")
		}
		creds = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(cp, ""))
	}
	conn, err := grpc.Dial(address, creds)
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	_, err = c.SayHello(context.Background(), &pb.HelloRequest{Name: t.Name()})
	if err != nil {
		return err
	}
	return nil
}

func TestX509(t *testing.T) {
	defer testenv.Clear().Restore()

	t.Run("insecure", func(t *testing.T) {
		l := log.Test(t, svc)
		assert := require.New(t)

		s, err := NewServer(l, defSrv)
		assert.NoError(err)
		assert.NotNil(s)
		serve(t, s, func() {
			assert.NoError(connectGRPC(t, s.port, ""))
		})
	})

	certE, keyE := genCert(t)
	os.Setenv("GRPC_CERT", certE)
	os.Setenv("GRPC_KEY", keyE)
	defer os.Unsetenv("GRPC_KEY")
	defer os.Unsetenv("GRPC_CERT")
	t.Run("env-certs", func(t *testing.T) {
		l := log.Test(t, svc)
		assert := require.New(t)

		s, err := NewServer(l, defSrv)
		assert.NoError(err)
		assert.NotNil(s)
		serve(t, s, func() {
			assert.Error(connectGRPC(t, s.port, ""))
			assert.NoError(connectGRPC(t, s.port, certE))
		})
	})

	certKP, keyKP := genCert(t)
	t.Run("X509KeyPair", func(t *testing.T) {
		l := log.Test(t, svc)
		assert := require.New(t)

		s, err := NewServer(l, defSrv, X509KeyPair(certKP, keyKP))
		assert.NoError(err)
		assert.NotNil(s)
		serve(t, s, func() {
			assert.Error(connectGRPC(t, s.port, ""))
			assert.Error(connectGRPC(t, s.port, certE))
			assert.NoError(connectGRPC(t, s.port, certKP))
		})
	})

	certLKP, keyLKP := genCert(t)
	certLKPF, keyLKPF := writeTLSFiles(t, certLKP, keyLKP)
	t.Run("LoadX509KeyPair", func(t *testing.T) {
		l := log.Test(t, svc)
		assert := require.New(t)

		s, err := NewServer(l, defSrv, LoadX509KeyPair(certLKPF, keyLKPF))
		assert.NoError(err)
		assert.NotNil(s)
		serve(t, s, func() {
			assert.Error(connectGRPC(t, s.port, ""))
			assert.Error(connectGRPC(t, s.port, certE))
			assert.Error(connectGRPC(t, s.port, certKP))
			assert.NoError(connectGRPC(t, s.port, certLKP))
		})
	})

	t.Run("assert-fail", func(t *testing.T) {
		l := log.Test(t, svc)
		assert := require.New(t)

		s, err := NewServer(l, defSrv, X509KeyPair(certKP, keyKP), LoadX509KeyPair(certLKP, keyLKP))
		assert.Error(err)
		assert.Nil(s)
	})
}
