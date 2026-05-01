package http

import "fmt"

func (cfg ServerConfig) IsGitHubAppAuthEnabled() bool {
	return cfg.GitHubAppID != 0 && cfg.GitHubAppInstallationID != 0 && cfg.GitHubAppPrivateKey != ""
}

func (cfg ServerConfig) Validate() error {
	anySet := cfg.GitHubAppID != 0 || cfg.GitHubAppInstallationID != 0 || cfg.GitHubAppPrivateKey != ""
	if anySet {
		if !cfg.IsGitHubAppAuthEnabled() {
			return fmt.Errorf("only some GitHub App auth settings were set; GITHUB_APP_ID, GITHUB_APP_INSTALLATION_ID, and GITHUB_APP_PRIVATE_KEY are all required")
		}
		if cfg.ScopeChallenge {
			return fmt.Errorf("GitHub App auth mode cannot be combined with --scope-challenge")
		}
	}
	return nil
}
