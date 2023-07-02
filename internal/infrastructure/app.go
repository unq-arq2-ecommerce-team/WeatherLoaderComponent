package infrastructure

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerDocs "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/docs"
	app "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/application"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/config"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/handlers"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"io"
)

// Application
// @title Weather Loader Component API
// @version 1.0
// @description api for final tp arq2
// @contact.name API SUPPORT
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @BasePath /
// @query.collection.format multi
type Application interface {
	Run() error
}

type ginApplication struct {
	logger                            domain.Logger
	config                            config.Config
	mongoClient                       *mongo.Client
	findCityCurrentTemperatureQuery   *app.FindCityCurrentTemperatureQuery
	getCityDayTemperatureAverageQuery *app.GetCityTemperatureAverageQuery
	getFindAllWeathersQuery           *app.GetFindAllWeathersQuery
}

func NewGinApplication(
	config config.Config,
	logger domain.Logger,
	mongoClient *mongo.Client,
	findCityCurrentTemperatureQuery *app.FindCityCurrentTemperatureQuery,
	getCityDayTemperatureAverageQuery *app.GetCityTemperatureAverageQuery,
	getFindAllWeathersQuery *app.GetFindAllWeathersQuery,
) Application {
	return &ginApplication{
		logger:                            logger,
		config:                            config,
		mongoClient:                       mongoClient,
		findCityCurrentTemperatureQuery:   findCityCurrentTemperatureQuery,
		getCityDayTemperatureAverageQuery: getCityDayTemperatureAverageQuery,
		getFindAllWeathersQuery:           getFindAllWeathersQuery,
	}
}

func (app *ginApplication) Run() error {
	swaggerDocs.SwaggerInfo.Host = fmt.Sprintf("localhost:%v", app.config.Port)

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	router := gin.Default()

	router.GET("/", handlers.HealthCheck(app.mongoClient))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	routerApi := router.Group("/api")

	middleware.InitMetrics()
	routerApi.Use(middleware.TracingRequestId(), middleware.PrometheusMiddleware(), otelgin.Middleware(config.OtlServiceName))

	routerApi.GET("/weather/city/:city/temperature", handlers.FindCityCurrentTemperatureHandler(app.logger, app.findCityCurrentTemperatureQuery))
	routerApi.GET("/weather/city/:city/temperature/average", handlers.GetCityTemperatureAverageHandler(app.logger, app.getCityDayTemperatureAverageQuery))

	routerApi.GET("/weather/metrics/collector", handlers.CollectWeatherMetricsHandler(app.logger, app.getFindAllWeathersQuery))

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.logger.Infof("running http server on port %d, and env %s", app.config.Port, app.config.Environment)
	return router.Run(fmt.Sprintf(":%v", app.config.Port))
}
