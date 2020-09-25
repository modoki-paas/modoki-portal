package config

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
)

type GitHub struct {
	URL          string `yaml:"url" json:"url" config:"GITHUB_URL"`
	ClientID     string `yaml:"client_id" json:"client_id" config:"GITHUB_CLIENT_ID"`
	ClientSecret string `yaml:"client_secret" json:"client_secret" config:"GITHUB_CLIENT_SECRET"`
	AppName      string `yaml:"app_name" json:"app_name" config:"GITHUB_APP_NAME"`
	RedirectURL  string `yaml:"redirect_url" json:"redirect_url" config:"GITHUB_REDIRECT_URL"`
}

type Config struct {
	ClusterAddress     string   `yaml:"cluster_address" json:"cluster_address" config:"K8S_CLUSTER_ADDRESS"`
	CAData             string   `yaml:"ca_data" json:"ca_data" config:"K8S_CERTIFICATE_AUTHORITY_DATA"`
	ClientID           string   `yaml:"client_id" json:"client_id" config:"OIDC_CLIENT_ID"`
	ClientSecret       string   `yaml:"client_secret" json:"client_secret" config:"OIDC_CLIENT_SECRET"`
	Scopes             []string `yaml:"scopes" json:"scopes" config:"OIDC_SCOPES"`
	RedirectURL        string   `yaml:"redirect_url" json:"redirect_url" config:"OIDC_REDIRECT_URL"`
	IssuerURL          string   `yaml:"issuer_url" json:"issuer_url" config:"OIDC_ISSUER_URL"`
	SessionStoreSecret string   `yaml:"session_store_secret" json:"session_store_secret" config:"SESSION_STORE_SECRET"`
	StaticFilesDir     string   `yaml:"static_files_dir" json:"static_files_dir" config:"STATIC_FILE_SERVING"`
	ReverseProxy       string   `yaml:"reverse_proxy" json:"reverse_proxy" config:"REVERSE_PROXY"`
	Local              bool     `yaml:"local" json:"local" config:"LOCAL"`

	GitHub GitHub `yaml:"github" json:"github"`
}

func ReadConfig() (*Config, error) {
	loader := confita.NewLoader(
		env.NewBackend(),
		file.NewOptionalBackend("proxy.yaml"),
		file.NewOptionalBackend("proxy.json"),
	)

	cfg := &Config{
		GitHub: GitHub{
			URL: "https://github.com",
		},
		SessionStoreSecret: "very-secure-secret",
	}

	if err := loader.Load(context.Background(), cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
