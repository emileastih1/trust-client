package metrics

import (
	"bulletin-board-api/internal/constants"
	"fmt"
	"os"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	MetricBulletinBoardGRPCResponse = "grpc.response"
	MetricBulletinBoardGRPCLatency  = "grpc.latency"
)

type StatsD struct {
	client statsd.ClientInterface
	logger *zap.Logger
}

func NewStatsD(logger *zap.Logger) (*StatsD, error) {
	statsdHost := os.Getenv("DD_AGENT_HOST")
	if statsdHost == "" {
		statsdHost = "127.0.0.1"
	}
	logger.Sugar().Infow("Metrics StatsD initialization", "host", statsdHost)
	client, err := statsd.New(
		fmt.Sprintf("%s:8125", statsdHost),
		statsd.WithNamespace(constants.BulletinBoardServiceName),
		statsd.WithTags([]string{
			fmt.Sprintf("environment:%s", viper.GetString("stage")),
		}),
		statsd.WithoutTelemetry(),
	)
	if err != nil {
		return nil, fmt.Errorf("datadog statsd initialization error: %w", err)
	}

	return &StatsD{
		client: client,
		logger: logger,
	}, nil
}

// Gauge measures the value of a metric at a particular time.
func (d *StatsD) Gauge(name string, value float64, options ...string) {
	if err := d.client.Gauge(name, value, makeTagsSlice(getTags(options...)), 1); err != nil {
		d.logger.Sugar().Warnw("Failed to update a gauge metric", "name", name, "value", value, "err", err)
	}
}

// Count tracks how many times something happened per second.
func (d *StatsD) Count(name string, value int64, options ...string) {
	if err := d.client.Count(name, value, makeTagsSlice(getTags(options...)), 1); err != nil {
		d.logger.Sugar().Warnw("Failed to update a count metic", "name", name, "value", value, "err", err)
	}
}

// Incr is just Count of 1.
func (d *StatsD) Incr(name string, options ...string) {
	if err := d.client.Incr(name, makeTagsSlice(getTags(options...)), 1); err != nil {
		d.logger.Sugar().Warnw("Failed to update a incr metric", "name", name, "err", err)
	}
}

// Histogram tracks the statistical distribution of a set of values on each host.
func (d *StatsD) Histogram(name string, value int64, options ...string) {
	if err := d.client.Histogram(name, float64(value), makeTagsSlice(getTags(options...)), 1); err != nil {
		d.logger.Sugar().Warnw("Failed to update a histogram metric", "name", name, "value", value, "err", err)
	}
}

// Timing sends timing information, it is an alias for TimeInMilliseconds.
func (d *StatsD) Timing(name string, value time.Duration, options ...string) {
	if err := d.client.Timing(name, value, makeTagsSlice(getTags(options...)), 1); err != nil {
		d.logger.Sugar().Warnw("Failed to update a timing metric", "name", name, "value", value, "err", err)
	}
}

func NewTest() *StatsD {
	return &StatsD{
		client: &statsd.NoOpClient{},
		logger: zap.NewNop(),
	}
}

type Timer struct {
	statsd    *StatsD
	startTime time.Time
	name      string
	tags      []string
}

func NewTimer(c *StatsD, name string, tags ...string) *Timer {
	return &Timer{
		statsd:    c,
		name:      name,
		startTime: time.Now(),
		tags:      tags,
	}
}

func (t *Timer) AppendTags(tags ...string) {
	t.tags = append(t.tags, tags...)
}

func (t *Timer) Emit() {
	t.statsd.Timing(t.name, time.Since(t.startTime), t.tags...)
}

func getTags(options ...string) map[string]string {
	tags := make(map[string]string)

	if len(options)%2 == 0 {
		for i := 1; i < len(options); i += 2 {
			tags[options[i-1]] = options[i]
		}
	}

	return tags
}

func makeTagsSlice(tagsMap map[string]string) []string {
	result := make([]string, 0, len(tagsMap))
	for k, v := range tagsMap {
		result = append(result, fmt.Sprintf("%s:%s", k, v))
	}
	return result
}
