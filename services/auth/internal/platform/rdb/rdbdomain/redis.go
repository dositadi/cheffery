package rdbdomain

import (
	"context"
	"fmt"

	"github.com/dositadi/cheffery/services/auth/internal/config"
	"github.com/dositadi/cheffery/services/shared/logger"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
	Metric *ConnMetric
}

func New(logger logger.Logger, cfg config.RedisConfig) *Redis {
	client, metric := newClient(logger, cfg)
	return &Redis{
		Client: client,
		Metric: metric,
	}
}

/* func redisUrl(cfg config.RedisConfig) string {
	u := url.URL{
		Scheme: cfg.Protocol,
		User: url.UserPassword(cfg.Username,cfg.Password),
		Host: cfg.Host,
	}
} */

func newClient(logger logger.Logger, cfg config.RedisConfig) (*redis.Client, *ConnMetric) {
	opt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%v", cfg.Host, cfg.Port),
		Username: cfg.Username,
		Password: cfg.Password,
	}
	scope := "rdbdomain.newClient()"

	opt.MaxActiveConns = 30
	opt.MaxIdleConns = 10

	metric := newConnMetric()
	opt.OnConnect = func(ctx context.Context, cn *redis.Conn) error {
		metric.onCreate()

		logger.PrintInfo("new-redis-conn", "Connection established, running initialization", map[string]string{
			"Context": scope,
		})
		return nil
	}

	client := redis.NewClient(opt)

	return client, metric
}
