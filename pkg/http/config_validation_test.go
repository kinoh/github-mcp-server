package http

import "testing"

func TestServerConfigValidate_GitHubAppAuth(t *testing.T) {
	t.Run("valid app auth config", func(t *testing.T) {
		cfg := ServerConfig{
			GitHubAppID:             1,
			GitHubAppInstallationID: 2,
			GitHubAppPrivateKey:     "-----BEGIN RSA PRIVATE KEY-----\nabc\n-----END RSA PRIVATE KEY-----",
		}
		if err := cfg.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("partial config fails", func(t *testing.T) {
		cfg := ServerConfig{GitHubAppID: 1}
		if err := cfg.Validate(); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("scope challenge conflict fails", func(t *testing.T) {
		cfg := ServerConfig{
			GitHubAppID:             1,
			GitHubAppInstallationID: 2,
			GitHubAppPrivateKey:     "pem",
			ScopeChallenge:          true,
		}
		if err := cfg.Validate(); err == nil {
			t.Fatal("expected error")
		}
	})
}
