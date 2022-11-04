package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/domain"
)

type ConfigGetter interface {
	HostCache() string
	PortCache() int
	UsernameCache() string
	PasswordCache() string
}

type ReportCache struct {
	week  *redis.Client
	month *redis.Client
	year  *redis.Client
}

func New(config ConfigGetter) *ReportCache {
	return &ReportCache{
		week: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.HostCache(), config.PortCache()),
			Username: config.UsernameCache(),
			Password: config.PasswordCache(),
			DB:       0,
		}),
		month: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.HostCache(), config.PortCache()),
			Username: config.UsernameCache(),
			Password: config.PasswordCache(),
			DB:       1,
		}),
		year: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.HostCache(), config.PortCache()),
			Username: config.UsernameCache(),
			Password: config.PasswordCache(),
			DB:       2,
		}),
	}
}

// WEEK

func (r *ReportCache) GetWeekReport(ctx context.Context, key string) *domain.Report {
	return nil
}

func (r *ReportCache) SetWeekReport(ctx context.Context, key string, value *domain.Report) error {
	return nil
}

func (r *ReportCache) RemoveWeekReport(ctx context.Context, key string) error {
	return nil
}

// MONTH

func (r *ReportCache) GetMonthReport(ctx context.Context, key string) *domain.Report {
	return nil
}

func (r *ReportCache) SetMonthReport(ctx context.Context, key string, value *domain.Report) error {
	return nil
}

func (r *ReportCache) RemoveMonthReport(ctx context.Context, key string) error {
	return nil
}

// YEAR

func (r *ReportCache) GetYearReport(ctx context.Context, key string) *domain.Report {
	return nil
}

func (r *ReportCache) SetYearReport(ctx context.Context, key string, value *domain.Report) error {
	return nil
}

func (r *ReportCache) RemoveYearReport(ctx context.Context, key string) error {
	return nil
}

// ALL

func (r *ReportCache) RemoveFromAll(ctx context.Context, key string) error
