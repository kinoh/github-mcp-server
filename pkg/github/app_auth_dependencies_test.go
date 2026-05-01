package github

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"
)

func TestRequestDepsRejectsPartialAppAuthConfig(t *testing.T) {
	deps := &RequestDeps{
		app: RequestDepsAppAuthConfig{AppID: 1},
	}

	_, err := deps.GetClient(context.Background())
	if err == nil {
		t.Fatal("expected partial app auth config to fail")
	}
	expected := "only some GitHub App auth settings were set; AppID, InstallationID, and PrivateKeyPEM are all required"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

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
	if strings.HasSuffix(err.Error(), ":") {
		t.Fatalf("expected wrapped error to include original details, got %q", err.Error())
	}
}

func TestRequestDepsGetOrCreateAppTransportCachesByBaseRESTURL(t *testing.T) {
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
	second, err := deps.getOrCreateAppTransport("https://api.github.com/")
	if err != nil {
		t.Fatalf("expected cached transport, got %v", err)
	}

	if first != second {
		t.Fatal("expected matching base URL to return the cached transport")
	}
	if second.BaseURL != "https://api.github.com/" {
		t.Fatalf("expected cached base URL, got %q", second.BaseURL)
	}

	third, err := deps.getOrCreateAppTransport("https://api.example.com/")
	if err != nil {
		t.Fatalf("expected separate transport, got %v", err)
	}
	if first == third {
		t.Fatal("expected different base URL to use a separate transport")
	}
	if third.BaseURL != "https://api.example.com/" {
		t.Fatalf("expected separate base URL, got %q", third.BaseURL)
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
