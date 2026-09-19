package config

import "github.com/dositadi/cheffery/services/shared/logger"

type ServiceConfig struct {
	AuthPort uint16
	AuthHost string
	RepoPort uint16
	RepoHost string
}

func loadServiceConfig(logger logger.Logger) *ServiceConfig {
	cfg := &ServiceConfig{
		AuthPort: getIntOrDefault("AUTH_GRPC_PORT", 50051),
		AuthHost: getStringOrDefault("AUTH_GRPC_HOST", "auth"),
		RepoPort: getIntOrDefault("REPO_GRPC_PORT", 50052),
		RepoHost: getStringOrDefault("REPO_GRPC_HOST", "repo"),
	}

	return cfg
}
