package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// Counter 示例
var httpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of HTTP requests",
	},
	[]string{"method", "path"},
)

// Gauge 示例
var systemLoad = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "system_load",
		Help: "Current system load",
	},
)

// Histogram 示例
var requestDuration = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Histogram of the duration of HTTP requests",
		Buckets: prometheus.DefBuckets,
	},
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(systemLoad)
	prometheus.MustRegister(requestDuration)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpRequestsTotal.With(prometheus.Labels{"method": r.Method, "path": r.URL.Path}).Inc()
		// 更新 Gauge
		//systemLoad.Set(getSystemLoad())
		timer := prometheus.NewTimer(requestDuration)
		defer timer.ObserveDuration()
		w.Write([]byte("Hello, world!"))
	})

	// 暴露 Prometheus 指标端点
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

}
