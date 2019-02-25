package grpc

import (
	"os"
	"strconv"
	"testing"
)

const (
	dummyCert = `-----BEGIN CERTIFICATE-----
MIIDaDCCAlACCQCB64J/4n9WmjANBgkqhkiG9w0BAQsFADB2MQswCQYDVQQGEwJV
UzERMA8GA1UECAwITmV3IFlvcmsxFjAUBgNVBAcMDU5ldyBZb3JrIENpdHkxGTAX
BgNVBAoMEFBhY2tldCBIb3N0LCBJbmMxDDAKBgNVBAsMA1NXRTETMBEGA1UEAwwK
cGFja2V0Lm5ldDAeFw0xOTAyMjQyMDMyNTZaFw0yOTAyMjEyMDMyNTZaMHYxCzAJ
BgNVBAYTAlVTMREwDwYDVQQIDAhOZXcgWW9yazEWMBQGA1UEBwwNTmV3IFlvcmsg
Q2l0eTEZMBcGA1UECgwQUGFja2V0IEhvc3QsIEluYzEMMAoGA1UECwwDU1dFMRMw
EQYDVQQDDApwYWNrZXQubmV0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAwFhH9gh30czIm7TlzTrR6UyCiSAwg5y0Ph68k2aG7udapxnAA6fRsTeUIuPx
tbc4ON0L3n3Wdo/tYuXK/iomLWAuWeafokhvQciHS3eDIklQvp8I85AikfeOd27d
EpohQDKMC09Jy5oflyPMqGUjpxM4SVLlORt7WZazeYXorJ2nk2ALf+Q/CRKctZJk
zzZhyPLCFkjxI+D/GCChmqK3vkVNB3rAj5OOvAnRDKMe/lw7wyOd93ux0NrINBGj
lPMYqAb2s7f9Q1TKpdnVZWboYtsQpwDxbDEZdXAcndFkYxKqwRjeSUOGPPaZpHHD
Mze+fUDv3jCN3ix+UDUNDubuGwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQAGGgI5
ownXuPppQmTjsWY+dyxSppGNdkDAafH0gaGmY3+Cx19YWolk/tQecHS0XJgZzFwP
VDvwbtYLausG3RwGwKSXY+wFUMPt/1g6orvuXNsidf7uLe26VxhfhzM9UKyy4QC0
bGvhaXYDozpVIMk9bGLz34ZKVEQUen91+3HOsCYyxTnIC4uLmZuOswZWoftwnubY
Q67IImj0zNBCmb/OVxhXclsZqIglvqIcH0uBwju6qo9BYOJkD3V6eAlw9Y7s2+LC
3ZrfHPp9VOt4JT+lOp3OhboeEHpKR44PgYi8I2UmuvwXmG832zrQytS1XghpiSsP
IUTdZkHNg8K+VjoQ
-----END CERTIFICATE-----`
	dummyKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAwFhH9gh30czIm7TlzTrR6UyCiSAwg5y0Ph68k2aG7udapxnA
A6fRsTeUIuPxtbc4ON0L3n3Wdo/tYuXK/iomLWAuWeafokhvQciHS3eDIklQvp8I
85AikfeOd27dEpohQDKMC09Jy5oflyPMqGUjpxM4SVLlORt7WZazeYXorJ2nk2AL
f+Q/CRKctZJkzzZhyPLCFkjxI+D/GCChmqK3vkVNB3rAj5OOvAnRDKMe/lw7wyOd
93ux0NrINBGjlPMYqAb2s7f9Q1TKpdnVZWboYtsQpwDxbDEZdXAcndFkYxKqwRje
SUOGPPaZpHHDMze+fUDv3jCN3ix+UDUNDubuGwIDAQABAoIBAQC3r9dRO/cJkfMG
2DQZ0ZGDpoCG6gnKts1fAcV/UwuLfcASEsJP+2WDQ5uh0mQT8NytWVQrb6tvYLYI
m4FHRwNclBzP2DIdLeWqQhIK3SCLjs6grIpE6CJLmcohfut7B1y3zU32wwqreQ2w
Lg0VyDjLJsy7IYItDnS3MvdFF8ADj4fTK3z19tI/7jj6w83mLcpw6B/3zUGNDa9I
Ob9fYXK+BlB7vwL37vJvahfar7mSbDl5dm/oRdd1pJbyUrXU6/02wzsxdFP/ZBcV
ancUA4Dn5zi3tP5HG8LtPN1eDLnAIGDPr69jSyXpsYqoCllZTknEQJfXOmg+HSJe
TkVH7IyZAoGBAOWpdGzzahvblsmlez77NKBZmT/XkIB64aG0kUp0PM1xH2KjTLEQ
0MfNU9l2FUJbMuue3g6buubVNN1JGcH/rSELnVub0D6CQr7cGOubDV1p3wbQUm99
i/oewu4o0yZXQK4b0u4vleUlo737QrNzxDmR5vxMmNTasw8nptlqflqNAoGBANZn
QAv7qBdUPdLvpJVwvpxWTWj/+4DYHvgfGKQxyUfR9+Kix7CDroyLw56WiREk3aan
J1j0VbbAF9J0PpTc28YuPxl8MX7feCNd3H8KJu8Pi89mk2ONlw5C9T30nUlcTgBL
PD9mL2r9QRi00Q+DFG57V0jABU3etYoO/hbzt1VHAoGBALpLzV+bzNUwOY71J5ad
W8E/LSs2h8dQ5rqvqLQGulPEkbsH0GxJwbJyArSCLxiWtiWfx21+MgyRosJmS/is
mBoYO9tV94TdUZtVGvnz2tGN0hbK4jQCWYvZbDKY9z9Aw/z4IRCJlUQ+VicELMU5
AVHZ4s+Cqu7vQRToC1aOJlT5AoGATYfqxiqLv1vsO2IDXzL1Cq2+snCW7yG4GTuN
epqyUbFg9Wit02va6+ICrE99Y2C0cnZRqT453KscMjNtCgHPy5ufn8SkVV/UHt3r
RVlTePFjOm26cK6b6EFYU74oPoYNgteyAq8eCI9qQdfpHbXl5ondp2YgxOb7OOBx
C7W4HzMCgYAER2ACIP8lQLmeCDUKUwyy7nQlwmINyDpL3+PlU69/oIhFC5sr/Lt/
Nzv33lWeTvLBK6dHOga9+PxHB8rCgCykjevLzHwxtjycTboL5PxJpd2EKxwsOgqX
UEbaPdu+CBypOyeoa5tCzDChy8oKWLraPfi9S+BarFlEoIqK2VAHbg==
-----END RSA PRIVATE KEY-----`
	defaultBindValue = "GRPC_BIND_TEST_VALUE"
	defaultPortValue = 50060
)

func setupEnv(envs map[string]string) {
	os.Setenv("GRPC_BIND", defaultBindValue)
	os.Setenv("GRPC_PORT", strconv.Itoa(defaultPortValue))
	os.Setenv("GRPC_SERVER_CERT", dummyCert)
	os.Setenv("GRPC_SERVER_KEY", dummyKey)
	for k, v := range envs {
		os.Setenv(k, v)
	}
}

func TestConfigFromEnvValid(t *testing.T) {
	setupEnv(nil)

	config, err := ConfigFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if config.Bind != defaultBindValue {
		t.Fatalf("expected=%s, got=%s", defaultBindValue, config.Bind)
	}
	if config.Port != defaultPortValue {
		t.Fatalf("expected=%d, got=%d", defaultPortValue, config.Port)
	}

	bindAddr, err := config.BindAddress()
	if err != nil {
		t.Fatal(err)
	}
	want := "GRPC_BIND_TEST_VALUE:50060"
	if bindAddr != want {
		t.Fatalf("error in retrieving BindAddress, want=%s, got=%s", want, bindAddr)
	}

	if config.ServerCredentials == nil {
		t.Fatal("expected non-nil ServerCredentials on config")
	}
}

func TestConfigFromEnvInvalid(t *testing.T) {

	t.Run("BadEnvs", func(t *testing.T) {
		setupEnv(nil)
		bindValue := "GRPC_BIND_TEST_VALUE"
		portValue := "NOT_AN_INT"
		os.Setenv("GRPC_BIND", bindValue)
		os.Setenv("GRPC_PORT", portValue)
		_, err := ConfigFromEnv()
		if err == nil {
			t.Fatalf("expected error parsing port value, got: %v", err)
		}
	})

	t.Run("UnsetEnvs", func(t *testing.T) {
		setupEnv(nil)
		os.Setenv("GRPC_BIND", "")
		os.Setenv("GRPC_PORT", "")

		config, err := ConfigFromEnv()
		if err == nil {
			t.Fatalf("expected error parsing port value, got: %v", err)
		}
		if config != nil {
			t.Fatalf("expected nil config, got: %v", config)
		}
	})

	t.Run("UnsetCert", func(t *testing.T) {
		setupEnv(nil)
		os.Setenv("GRPC_SERVER_CERT", "")
		config, err := ConfigFromEnv()
		if err != nil {
			t.Fatal(err)
		}
		if config.ServerCredentials != nil {
			t.Fatal("expected nil ServerCredentials")
		}
	})

	t.Run("UnsetKey", func(t *testing.T) {
		setupEnv(nil)
		os.Setenv("GRPC_SERVER_KEY", "")
		config, err := ConfigFromEnv()
		if err != nil {
			t.Fatal(err)
		}
		if config.ServerCredentials != nil {
			t.Fatal("expected nil ServerCredentials")
		}
	})

	t.Run("BadKey", func(t *testing.T) {
		setupEnv(nil)
		os.Setenv("GRPC_SERVER_KEY", "abc")
		_, err := ConfigFromEnv()
		if err == nil {
			t.Fatal("expected parse error")
		}
	})

	t.Run("BadCert", func(t *testing.T) {
		setupEnv(nil)
		os.Setenv("GRPC_SERVER_CERT", "abc")
		_, err := ConfigFromEnv()
		if err == nil {
			t.Fatal("expected parse error")
		}
	})
}
