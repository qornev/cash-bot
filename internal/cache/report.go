package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
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

func (r *ReportCache) GetWeekReport(ctx context.Context, key int64) (string, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get week report from cache")
	defer span.Finish()

	var report string
	err := r.week.Get(ctx, fmt.Sprint(key)).Scan(&report)
	if err != nil {
		return "", false
	}

	return report, true
}

func (r *ReportCache) SetWeekReport(ctx context.Context, key int64, value string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "save week report in cache")
	defer span.Finish()

	err := r.week.Set(ctx, fmt.Sprint(key), value, time.Hour*24*7).Err()
	if err != nil {
		return err
	}
	return nil
}

// MONTH

func (r *ReportCache) GetMonthReport(ctx context.Context, key int64) (string, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get month report from cache")
	defer span.Finish()

	var report string
	err := r.month.Get(ctx, fmt.Sprint(key)).Scan(&report)
	if err != nil {
		return "", false
	}

	return report, true
}

func (r *ReportCache) SetMonthReport(ctx context.Context, key int64, value string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "save month report in cache")
	defer span.Finish()

	untilDate := time.Now().AddDate(0, 1, 0)
	untilDate = time.Date(untilDate.Year(), untilDate.Month(), 1, 0, 0, 0, 0, untilDate.Location())

	err := r.month.Set(ctx, fmt.Sprint(key), value, time.Until(untilDate)).Err()
	if err != nil {
		return err
	}
	return nil
}

// YEAR

func (r *ReportCache) GetYearReport(ctx context.Context, key int64) (string, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get year report from cache")
	defer span.Finish()

	var report string
	err := r.year.Get(ctx, fmt.Sprint(key)).Scan(&report)
	if err != nil {
		return "", false
	}

	return report, true
}

func (r *ReportCache) SetYearReport(ctx context.Context, key int64, value string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "save year report in cache")
	defer span.Finish()

	untilDate := time.Now().AddDate(1, 0, 0)
	untilDate = time.Date(untilDate.Year(), 1, 1, 0, 0, 0, 0, untilDate.Location())

	err := r.year.Set(ctx, fmt.Sprint(key), value, time.Until(untilDate)).Err()
	if err != nil {
		return err
	}
	return nil
}

// ALL

func (r *ReportCache) RemoveFromAll(ctx context.Context, keys []int64) error {
	var strKeys []string
	for _, key := range keys {
		strKeys = append(strKeys, fmt.Sprint(key))
	}
	if err := r.week.Del(ctx, strKeys...).Err(); err != nil {
		return err
	}
	if err := r.month.Del(ctx, strKeys...).Err(); err != nil {
		return err
	}
	if err := r.year.Del(ctx, strKeys...).Err(); err != nil {
		return err
	}
	return nil
}
