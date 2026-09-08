package config

import "time"

type RetryConfig struct {
	MinWait    time.Duration
	MaxWait    time.Duration
	MaxAttempt int
}

func defaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MinWait:    time.Duration(200 * time.Millisecond),
		MaxWait:    time.Duration(2 * time.Second),
		MaxAttempt: 5,
	}
}
