package grpc

import (
	"os"
	"strconv"
	"testing"
)

func TestConfigFromEnv(t *testing.T) {
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
}
