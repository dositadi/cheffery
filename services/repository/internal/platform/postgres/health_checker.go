package postgres

import (
	"context"
	"time"

	"github.com/dositadi/cheffery/services/shared/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Status string

const (
	UnHealthy = Status("unhealthy")
	Healthy   = Status("healthy")
	Degraded  = Status("degraded")
)

type HealthStatus struct {
	Status    Status
	DBHealth  DBHealth
	Latency   time.Duration
	Timestamp time.Time
}

type DBHealth struct {
	Connected      bool
	PoolHealthy    bool
	TotalConns     int32
	IdleConns      int32
	AcquiredConns  int32
	AcquiredCount  int32
	CreatedCount   int32
	ReleasedCount  int32
	DestroyedCount int32
	Message        string
}

type HealthChecker struct {
	logger  logger.Logger
	pool    *pgxpool.Pool
	metrics *Metrics
	timeout time.Duration
}

func NewHealthChecker(logger logger.Logger, pool *pgxpool.Pool, metrics *Metrics, timeout time.Duration) *HealthChecker {
	return &HealthChecker{
		logger:  logger,
		pool:    pool,
		metrics: metrics,
		timeout: timeout,
	}
}

func (h *HealthChecker) Check(ctx context.Context) *HealthStatus {
	scope := "postgres.Check"
	start := time.Now()

	healthStatus := &HealthStatus{
		Timestamp: start,
	}

	poolStat := h.pool.Stat()

	dbHealth := DBHealth{
		TotalConns:     poolStat.TotalConns(),
		IdleConns:      poolStat.IdleConns(),
		AcquiredConns:  poolStat.AcquiredConns(),
		AcquiredCount:  int32(h.metrics.GetOnAcquire()),
		CreatedCount:   int32(h.metrics.GetOnCreate()),
		ReleasedCount:  int32(h.metrics.GetOnRelease()),
		DestroyedCount: int32(h.metrics.GetOnDestroyed()),
	}

	if poolStat.TotalConns() >= poolStat.MaxConns() && poolStat.IdleConns() == 0 {
		dbHealth.PoolHealthy = false
		dbHealth.Message = "connection pool exhausted"
	} else {
		dbHealth.PoolHealthy = true
	}

	// Test DB connectivity
	var result int

	if err := h.pool.QueryRow(ctx, "SELECT 1").Scan(&result); err != nil {
		dbHealth.Connected = false
		dbHealth.Message = err.Error()
		healthStatus.Status = UnHealthy
		h.logger.PrintError(err, "db:health-check", err.Error(), map[string]string{
			"Context": scope,
		})
	} else {
		dbHealth.Connected = true

		if dbHealth.PoolHealthy {
			healthStatus.Status = Healthy
		} else {
			healthStatus.Status = Degraded
		}
	}

	healthStatus.DBHealth = dbHealth
	healthStatus.Latency = time.Since(start)

	return healthStatus
}
