package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "respond_now_server_api_requests_total",
			Help: "Total number of requests by status code.",
		},
		[]string{"status", "path"},
	)
	ErrorRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "respond_now_server_api_error_requests_total",
			Help: "Total number of error requests by status code.",
		},
		[]string{"status", "path"},
	)
	Error4xxRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "respond_now_server_api_error_4xx_requests_total",
			Help: "Total number of 4xx error requests by status code.",
		},
		[]string{"path"},
	)
	Error5xxRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "respond_now_server_api_error_5xx_requests_total",
			Help: "Total number of 5xx error requests by status code.",
		},
		[]string{"path"},
	)
	ResponseTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "respond_now_server_api_response_time_milli_seconds",
			Help:    "Response time of API requests in milli seconds.",
			Buckets: []float64{1, 10, 20, 50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000, 2000, 3000, 4000, 5000, 10000, 20000}, // 1ms to 20s
		},
		[]string{"path"},
	)
	ResponseTimeInSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "respond_now_server_api_response_time_seconds",
			Help:    "Response time of API requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func Init() {
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(ErrorRequests)
	prometheus.MustRegister(ResponseTime)
	prometheus.MustRegister(ResponseTimeInSeconds)
	prometheus.MustRegister(Error4xxRequests)
	prometheus.MustRegister(Error5xxRequests)
}
