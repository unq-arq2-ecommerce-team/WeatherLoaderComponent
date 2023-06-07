package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/application"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/dto"
	"net/http"
)

// GetCityTemperatureAverageHandler
// @Summary      Endpoint get city average temperature between dates
// @Description  get city current temperature in range of dates
// @Param city path string true "City" example("Quilmes")
// @Param        dateFrom   	query	string	true	"range dateFrom in UTC to get average temp"		example("2023-06-06T00:00:00.000Z")
// @Param        dateTo		   	query	string	true	"range dateTo in UTC to get average temp"		example("2023-06-07T00:00:00.000Z")
// @Tags         Weather
// @Produce json
// @Success 200 {object} domain.AverageTemperature
// @Success 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router       /api/weather/city/{city}/temperature/average [get]
func GetCityTemperatureAverageHandler(logger domain.Logger, getCityDayTemperatureAverageQuery *application.GetCityTemperatureAverageQuery) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.WithRequestId(c)
		log.Debug("findCityCurrentTemperatureHandler init")
		cityParam, _ := c.Params.Get("city")
		if cityParam == "" {
			log.WithFields(domain.LoggerFields{"cityParam": cityParam}).Errorf("city param is empty")
			c.JSON(http.StatusBadRequest, dto.NewErrorMessage("bad request", "empty city path param"))
			return
		}
		var queryParams dto.TemperatureAverageQuery
		if err := c.ShouldBindQuery(&queryParams); err != nil || queryParams.InvalidDates() {
			log.WithFields(domain.LoggerFields{"cityParam": cityParam, "error": err}).Errorf("queryParams cannot be bound")
			c.JSON(http.StatusBadRequest, dto.NewErrorMessage("invalid query params or date range", ""))
			return
		}

		avgTemp, err := getCityDayTemperatureAverageQuery.Do(c.Request.Context(), cityParam, queryParams.GetDateFrom(), queryParams.GetDateTo())
		if err != nil {
			switch err.(type) {
			case domain.AverageTemperatureNotFoundErr:
				c.JSON(http.StatusNotFound, dto.NewErrorMessage("average temperature of city not found", err.Error()))
			default:
				c.JSON(http.StatusInternalServerError, dto.NewErrorMessage("internal server error", err.Error()))
			}
			return
		}
		c.JSON(http.StatusOK, avgTemp)

	}
}
