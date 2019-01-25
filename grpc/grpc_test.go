package grpc

import (
	"os"
	"strconv"
	"testing"
)

func TestConfigFromEnvValid(t *testing.T) {
	bindValue := "GRPC_BIND_TEST_VALUE"
	portValue := 50060
	os.Setenv("GRPC_BIND", bindValue)
	os.Setenv("GRPC_PORT", strconv.Itoa(portValue))
	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if config.Bind != bindValue {
		t.Fatalf("expected=%s, got=%s", bindValue, config.Bind)
	}
	if config.Port != portValue {
		t.Fatalf("expected=%d, got=%d", portValue, config.Port)
	}

	bindAddr, err := config.BindAddress()
	if err != nil {
		t.Fatal(err)
	}
	want := "GRPC_BIND_TEST_VALUE:50060"
	if bindAddr != want {
		t.Fatalf("error in retrieving BindAddress, want=%s, got=%s", want, bindAddr)
	}
}

func TestConfigFromEnvInvalid(t *testing.T) {
	bindValue := "GRPC_BIND_TEST_VALUE"
	portValue := "NOT_AN_INT"
	os.Setenv("GRPC_BIND", bindValue)
	os.Setenv("GRPC_PORT", portValue)
	_, err := ConfigFromEnv()
	if err == nil {
		t.Fatalf("expected error parsing port value, got: %v", err)
	}
}
