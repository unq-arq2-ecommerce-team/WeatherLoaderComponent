package infrastructure

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerDocs "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/docs"
	app "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/application"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/config"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/handlers"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"io"
	"net/http"
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
	findCityCurrentTemperatureQuery   *app.FindCityCurrentTemperatureQuery
	getCityDayTemperatureAverageQuery *app.GetCityTemperatureAverageQuery
}

func NewGinApplication(
	config config.Config,
	logger domain.Logger,
	findCityCurrentTemperatureQuery *app.FindCityCurrentTemperatureQuery,
	getCityDayTemperatureAverageQuery *app.GetCityTemperatureAverageQuery,
) Application {
	return &ginApplication{
		logger:                            logger,
		config:                            config,
		findCityCurrentTemperatureQuery:   findCityCurrentTemperatureQuery,
		getCityDayTemperatureAverageQuery: getCityDayTemperatureAverageQuery,
	}
}

func (app *ginApplication) Run() error {
	swaggerDocs.SwaggerInfo.Host = fmt.Sprintf("localhost:%v", app.config.Port)

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	router := gin.Default()
	router.Use(otelgin.Middleware("weather-loader"))

	router.GET("/", HealthCheck)

	routerApi := router.Group("/api")
	routerApi.Use(middleware.TracingRequestId())

	routerApi.GET("/weather/city/:city/temperature", handlers.FindCityCurrentTemperatureHandler(app.logger, app.findCityCurrentTemperatureQuery))
	routerApi.GET("/weather/city/:city/temperature/average", handlers.GetCityTemperatureAverageHandler(app.logger, app.getCityDayTemperatureAverageQuery))

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.logger.Infof("running http server on port %d", app.config.Port)
	return router.Run(fmt.Sprintf(":%v", app.config.Port))
}

// HealthCheck
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags Health check
// @Accept */*
// @Produce json
// @Success 200 {object} HealthCheckRes
// @Router / [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthCheckRes{Data: "Server is up and running"})
}

type HealthCheckRes struct {
	Data string `json:"data"`
}
