package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/logger"
	"time"
)

const headerRequestId = "system-request-id"

func TracingRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := logger.SetRequestId(c.Request.Context(), c.Request.Header.Get(headerRequestId))
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(headerRequestId, logger.GetRequestId(ctx))
	}
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "weather_loader_http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path", "method"},
)

var requestsDurationSecondsSum = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "weather_loader_http_request_duration_seconds_sum",
	Help: "Sum of seconds spent on all requests",
})

var requestsDurationSecondsCount = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "weather_loader_http_request_duration_seconds_count",
	Help: "Count of  all requests",
})

var requestsDurationSecondsBucket = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "weather_loader_http_request_duration_seconds_bucket",
	Help: "group request by tiem repsonses tags",
},
	[]string{"le"})

var httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "weather_loader_http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func InitMetrics() {
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(requestsDurationSecondsSum)
	prometheus.MustRegister(requestsDurationSecondsCount)
	prometheus.MustRegister(requestsDurationSecondsBucket)
	prometheus.MustRegister(httpDuration)
	prometheus.MustRegister(weatherData)
}
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		c.Next()

		totalRequests.WithLabelValues(path, method).Inc()

		time := timer.ObserveDuration()
		incrementRequestOfDuration(time.Seconds())

	}
}
func incrementRequestOfDuration(d float64) {
	requestsDurationSecondsCount.Inc()
	requestsDurationSecondsSum.Add(d)
	if d <= 0.05 {
		requestsDurationSecondsBucket.WithLabelValues("0.05").Inc()
		return
	}
	if d <= 0.3 {
		requestsDurationSecondsBucket.WithLabelValues("0.3").Inc()
		return
	}
	if d <= 0.5 {
		requestsDurationSecondsBucket.WithLabelValues("0.5").Inc()
		return
	}
	if d <= 1 {
		requestsDurationSecondsBucket.WithLabelValues("1").Inc()
		return
	}
	if d <= 5 {
		requestsDurationSecondsBucket.WithLabelValues("5").Inc()
	}
}

var weatherData = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "weather_loader_weather_temperature",
		Help: "Temperature in Celsius",
	},
	[]string{"city", "timestamp"},
)

func SaveWeatherData(weather domain.Weather) {
	timestamp := weather.TimeStamp.Format(time.RFC3339)
	weatherData.WithLabelValues(weather.City, timestamp).Set(weather.Temperature)
}
