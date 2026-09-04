package config

type ApplicationConfig struct {
	RetryCfg *RetryConfig
}

func LoadAppConfig() *ApplicationConfig {
	cfg := &ApplicationConfig{
		RetryCfg: defaultRetryConfig(),
	}

	return cfg
}
