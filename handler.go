package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HTTPHandlerOpts describes a set of configurable options of http metrics
type HTTPHandlerOpts struct {
	DurationBuckets     []float64
	RequestSizeBuckets  []float64
	ResponseSizeBuckets []float64
}

// HTTPMetric describes a metric used to instrument an HTTP handler.
//
// HTTPMetric values must be created using the HTTP metric methods on
// [Namespace], such as [Namespace.NewDefaultHttpMetrics] or
// [Namespace.NewHttpMetricsWithOpts].
type HTTPMetric struct {
	collector prometheus.Collector
	wrap      func(http.Handler) http.Handler
}

var _ prometheus.Collector = (*HTTPMetric)(nil)

// Describe implements [prometheus.Collector].
func (m *HTTPMetric) Describe(ch chan<- *prometheus.Desc) {
	m.collector.Describe(ch)
}

// Collect implements [prometheus.Collector].
func (m *HTTPMetric) Collect(ch chan<- prometheus.Metric) {
	m.collector.Collect(ch)
}

var (
	defaultDurationBuckets     = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 25, 60}
	defaultRequestSizeBuckets  = prometheus.ExponentialBuckets(1024, 2, 22) // 1K to 4G
	defaultResponseSizeBuckets = defaultRequestSizeBuckets
)

// Handler returns the global http.Handler that provides the prometheus
// metrics format on GET requests. This handler is no longer instrumented.
func Handler() http.Handler {
	return promhttp.Handler()
}

func InstrumentHandler(metrics []*HTTPMetric, handler http.Handler) http.HandlerFunc {
	return instrumentHandler(metrics, handler)
}

func InstrumentHandlerFunc(metrics []*HTTPMetric, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return instrumentHandler(metrics, handlerFunc)
}

func instrumentHandler(metrics []*HTTPMetric, handler http.Handler) http.HandlerFunc {
	for _, metric := range metrics {
		handler = metric.wrap(handler)
	}
	return handler.ServeHTTP
}
