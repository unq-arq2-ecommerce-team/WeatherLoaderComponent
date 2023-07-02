package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/application"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/dto"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/middleware"
	"net/http"
)

func CollectWeatherMetricsHandler(logger domain.Logger, findAllTemperature *application.GetFindAllWeathersQuery) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.WithRequestId(c)
		log.Debug("CollectWeatherMetrics init")

		allWeathers, err := findAllTemperature.Do(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorMessage("internal server error", err.Error()))
		}

		for _, weather := range *allWeathers {
			middleware.SaveWeatherData(weather)
		}

		c.JSON(http.StatusOK, "all weather loaded")
	}
}
