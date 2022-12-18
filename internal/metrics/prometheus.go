package metrics

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Prometheus struct {
	config *PrometheusConfig

	textUpdatesReceivedHistogram *prometheus.HistogramVec
	totalStartOperationHistogram *prometheus.HistogramVec
	totalSendOperationHistogram  *prometheus.HistogramVec
}

func NewPrometheusCollector(config *PrometheusConfig, reg prometheus.Registerer) *Prometheus {
	updates := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "stats",
		Name:      "total_updates",
		Help:      "Shows how many text updates has been received",
		Buckets:   prometheus.DefBuckets,
	}, []string{"value"})
	starts := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "stats",
		Name:      "total_start",
		Help:      "Show how many times users executed /start command",
	}, []string{"value"})
	sends := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "stats",
		Name:      "total_send",
		Help:      "Show how many times users executed /send command",
	}, []string{"value"})

	p := &Prometheus{
		textUpdatesReceivedHistogram: updates,
		totalStartOperationHistogram: starts,
		totalSendOperationHistogram:  sends,
	}

	reg.MustRegister(
		p.textUpdatesReceivedHistogram,
		p.totalStartOperationHistogram,
		p.totalSendOperationHistogram,
	)

	return p
}

func (p *Prometheus) ObserveTextUpdatesDuration(ctx context.Context, value int, t time.Duration) {
	p.textUpdatesReceivedHistogram.WithLabelValues(strconv.Itoa(value)).Observe(t.Seconds())
}

func (p *Prometheus) ObserveTotalStartOperationDuration(ctx context.Context, value int, t time.Duration) {
	p.totalStartOperationHistogram.WithLabelValues(strconv.Itoa(value)).Observe(t.Seconds())
}

func (p *Prometheus) ObserveTotalSendOperationDuration(ctx context.Context, value int, t time.Duration) {
	p.totalSendOperationHistogram.WithLabelValues(strconv.Itoa(value)).Observe(t.Seconds())
}
