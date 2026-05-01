package http

import "fmt"

func (cfg ServerConfig) IsGitHubAppAuthEnabled() bool {
	return cfg.GitHubAppID != 0 || cfg.GitHubAppInstallationID != 0 || cfg.GitHubAppPrivateKey != ""
}

func (cfg ServerConfig) Validate() error {
	enabled := cfg.IsGitHubAppAuthEnabled()
	if enabled {
		if cfg.GitHubAppID == 0 || cfg.GitHubAppInstallationID == 0 || cfg.GitHubAppPrivateKey == "" {
			return fmt.Errorf("GITHUB_APP_ID, GITHUB_APP_INSTALLATION_ID, and GITHUB_APP_PRIVATE_KEY must all be set when GitHub App auth mode is enabled")
		}
		if cfg.ScopeChallenge {
			return fmt.Errorf("GitHub App auth mode cannot be combined with --scope-challenge")
		}
	}
	return nil
}
