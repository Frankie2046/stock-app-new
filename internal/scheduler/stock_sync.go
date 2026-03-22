package scheduler

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"project/internal/config"
)

type SymbolSyncer interface {
	SyncSymbol(symbol string) error
}

func StartStockSync(cfg *config.Config, logger *zap.Logger, syncer SymbolSyncer) (*cron.Cron, error) {
	symbols := parseSymbols(cfg.SyncSymbols)
	if len(symbols) == 0 {
		logger.Warn("stock sync scheduler disabled: no symbols configured")
		return nil, nil
	}

	spec := strings.TrimSpace(cfg.SyncCronExpr)
	if spec == "" {
		spec = "*/30 * * * *"
	}

	c := cron.New()
	var running atomic.Bool

	runRound := func(trigger string) {
		if !running.CompareAndSwap(false, true) {
			logger.Warn("stock sync skipped: previous round still running", zap.String("trigger", trigger))
			return
		}
		defer running.Store(false)
		defer func() {
			if r := recover(); r != nil {
				logger.Error(
					"panic recovered in stock sync round",
					zap.String("trigger", trigger),
					zap.Any("panic", r),
				)
			}
		}()

		start := time.Now()
		success := 0
		failed := 0

		for i, symbol := range symbols {
			if err := syncer.SyncSymbol(symbol); err != nil {
				failed++
				logger.Warn("sync symbol failed", zap.String("symbol", symbol), zap.Error(err))
			} else {
				success++
				logger.Info("sync symbol success", zap.String("symbol", symbol))
			}

			// AlphaVantage free plan limit-friendly pacing.
			if i < len(symbols)-1 {
				time.Sleep(12 * time.Second)
			}
		}

		logger.Info(
			"stock sync round finished",
			zap.String("trigger", trigger),
			zap.Int("total", len(symbols)),
			zap.Int("success", success),
			zap.Int("failed", failed),
			zap.Duration("cost", time.Since(start)),
		)
	}

	if cfg.SyncRunOnStart {
		go runRound("startup")
	}

	if _, err := c.AddFunc(spec, func() { runRound("cron") }); err != nil {
		return nil, fmt.Errorf("invalid cron expression %q: %w", spec, err)
	}

	c.Start()
	logger.Info("stock sync scheduler started", zap.String("cron", spec), zap.Strings("symbols", symbols), zap.Bool("run_on_start", cfg.SyncRunOnStart))
	return c, nil
}

func parseSymbols(raw string) []string {
	parts := strings.Split(raw, ",")
	seen := make(map[string]struct{}, len(parts))
	out := make([]string, 0, len(parts))

	for _, p := range parts {
		s := strings.ToUpper(strings.TrimSpace(p))
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
