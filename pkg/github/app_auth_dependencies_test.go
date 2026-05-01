package github

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"
)

func TestRequestDepsGetOrCreateAppTransportWrapsInitializationError(t *testing.T) {
	deps := &RequestDeps{
		app: RequestDepsAppAuthConfig{
			AppID:          1,
			InstallationID: 2,
			PrivateKeyPEM:  "not a pem",
		},
	}

	_, err := deps.getOrCreateAppTransport("https://api.github.com/")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "failed to initialize GitHub App transport:") {
		t.Fatalf("expected wrapped initialization error, got %q", err.Error())
	}
	if err.Error() == "failed to initialize GitHub App transport" {
		t.Fatal("expected original error details to be preserved")
	}
}

func TestRequestDepsGetOrCreateAppTransportCachesInitialBaseRESTURL(t *testing.T) {
	deps := &RequestDeps{
		app: RequestDepsAppAuthConfig{
			AppID:          1,
			InstallationID: 2,
			PrivateKeyPEM:  testRSAPrivateKeyPEM(t),
		},
	}

	first, err := deps.getOrCreateAppTransport("https://api.github.com/")
	if err != nil {
		t.Fatalf("expected first transport, got %v", err)
	}
	second, err := deps.getOrCreateAppTransport("https://api.example.com/")
	if err != nil {
		t.Fatalf("expected cached transport, got %v", err)
	}

	if first != second {
		t.Fatal("expected sync.Once to return the cached transport")
	}
	if second.BaseURL != "https://api.github.com/" {
		t.Fatalf("expected cached base URL, got %q", second.BaseURL)
	}
}

func testRSAPrivateKeyPEM(t *testing.T) string {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return string(pem.EncodeToMemory(block))
}
