package middlewares_prometheus

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// Define Prometheus metrics
var (
	// Counter for total HTTP requests
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	// Histogram for request duration
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Counter for errors
	httpRequestErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_errors_total",
			Help: "Total number of HTTP request errors",
		},
		[]string{"method", "endpoint", "error_type"},
	)

	// Gauge for active connections
	activeConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)

	// Gauge for application uptime
	appUptime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_uptime_seconds",
			Help: "Application uptime in seconds",
		},
	)

	// Gauge for application info
	appInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_info",
			Help: "Application information",
		},
		[]string{"version", "build_time"},
	)

	// Gin-specific metrics
	ginRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gin_requests_in_flight",
			Help: "Number of requests currently being processed",
		},
	)

	ginRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gin_request_size_bytes",
			Help:    "Size of HTTP requests in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7),
		},
		[]string{"method", "endpoint"},
	)

	ginResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gin_response_size_bytes",
			Help:    "Size of HTTP responses in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7),
		},
		[]string{"method", "endpoint"},
	)

	// Database metrics
	dbConnectionsActive = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Number of active database connections",
		},
	)

	dbConnectionsIdle = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	// Redis metrics
	redisConnectionsActive = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "redis_connections_active",
			Help: "Number of active Redis connections",
		},
	)

	redisOperationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_operations_total",
			Help: "Total number of Redis operations",
		},
		[]string{"operation", "status"},
	)

	// Business metrics
	userRegistrations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_registrations_total",
			Help: "Total number of user registrations",
		},
		[]string{"method"}, // email, social, etc.
	)

	userLogins = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_logins_total",
			Help: "Total number of user logins",
		},
		[]string{"method", "status"}, // success, failed
	)

	commentsCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "comments_created_total",
			Help: "Total number of comments created",
		},
	)

	placesCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "places_created_total",
			Help: "Total number of places created",
		},
	)
)

var (
	startTime time.Time
	isInitialized bool
)

func PrometheusInit() {
	if isInitialized {
		return // Prevent double initialization
	}

	startTime = time.Now()

	// Register metrics with Prometheus
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpRequestErrors)
	prometheus.MustRegister(activeConnections)
	prometheus.MustRegister(appUptime)
	prometheus.MustRegister(appInfo)
	prometheus.MustRegister(ginRequestsInFlight)
	prometheus.MustRegister(ginRequestSize)
	prometheus.MustRegister(ginResponseSize)
	prometheus.MustRegister(dbConnectionsActive)
	prometheus.MustRegister(dbConnectionsIdle)
	prometheus.MustRegister(redisConnectionsActive)
	prometheus.MustRegister(redisOperationsTotal)
	prometheus.MustRegister(userRegistrations)
	prometheus.MustRegister(userLogins)
	prometheus.MustRegister(commentsCreated)
	prometheus.MustRegister(placesCreated)

	// Set application info
	appInfo.WithLabelValues("1.0.0", startTime.Format("2006-01-02T15:04:05Z")).Set(1)

	// Start uptime updater
	go updateUptime()

	isInitialized = true
}

// Gin middleware for Prometheus metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip metrics endpoint
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()

		// Increment active connections and requests in flight
		activeConnections.Inc()
		ginRequestsInFlight.Inc()

		// Record request size
		if c.Request.ContentLength > 0 {
			ginRequestSize.WithLabelValues(c.Request.Method, c.FullPath()).Observe(float64(c.Request.ContentLength))
		}

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get status code
		statusCode := c.Writer.Status()
		statusCodeStr := strconv.Itoa(statusCode)

		// Record metrics
		httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), statusCodeStr).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
		ginResponseSize.WithLabelValues(c.Request.Method, c.FullPath()).Observe(float64(c.Writer.Size()))

		// Record errors for 4xx and 5xx status codes
		if statusCode >= 400 {
			errorType := "client_error"
			if statusCode >= 500 {
				errorType = "server_error"
			}
			httpRequestErrors.WithLabelValues(c.Request.Method, c.FullPath(), errorType).Inc()
		}

		// Decrement counters
		activeConnections.Dec()
		ginRequestsInFlight.Dec()
	}
}

// Update uptime metric periodically
func updateUptime() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			appUptime.Set(time.Since(startTime).Seconds())
		}
	}
}

// Business metric helpers - call these from your handlers
func IncrementUserRegistration(method string) {
	userRegistrations.WithLabelValues(method).Inc()
}

func IncrementUserLogin(method, status string) {
	userLogins.WithLabelValues(method, status).Inc()
}

func IncrementCommentCreated() {
	commentsCreated.Inc()
}

func IncrementPlaceCreated() {
	placesCreated.Inc()
}

func RecordRedisOperation(operation, status string) {
	redisOperationsTotal.WithLabelValues(operation, status).Inc()
}

// Database metrics updater - call this periodically or in a health check
func UpdateDatabaseMetrics(activeConns, idleConns int) {
	dbConnectionsActive.Set(float64(activeConns))
	dbConnectionsIdle.Set(float64(idleConns))
}

func UpdateRedisMetrics(activeConns int) {
	redisConnectionsActive.Set(float64(activeConns))
}

// Health check handler with metrics
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(startTime).Seconds(),
	})
}