package http

import "testing"

func TestServerConfigValidate_GitHubAppAuth(t *testing.T) {
	t.Run("empty app auth config is valid", func(t *testing.T) {
		cfg := ServerConfig{}
		if err := cfg.Validate(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

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
		err := cfg.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		expected := "only some GitHub App auth settings were set; GITHUB_APP_ID, GITHUB_APP_INSTALLATION_ID, and GITHUB_APP_PRIVATE_KEY are all required"
		if err.Error() != expected {
			t.Fatalf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("scope challenge conflict fails", func(t *testing.T) {
		cfg := ServerConfig{
			GitHubAppID:             1,
			GitHubAppInstallationID: 2,
			GitHubAppPrivateKey:     "pem",
			ScopeChallenge:          true,
		}
		err := cfg.Validate()
		if err == nil {
			t.Fatal("expected error")
		}
		expected := "GitHub App auth mode cannot be combined with --scope-challenge"
		if err.Error() != expected {
			t.Fatalf("expected %q, got %q", expected, err.Error())
		}
	})
}
