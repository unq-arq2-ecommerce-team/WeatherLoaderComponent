package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/repository/mongo"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

const serverOkMsg = "Server is up and running"

// HealthCheck
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags Health check
// @Accept */*
// @Produce json
// @Success 200 {object} HealthCheckRes
// @Failure 500 {object} HealthCheckRes
// @Router / [get]
func HealthCheck(mongoClient *mongoDriver.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		mongoMsg, err := mongo.HealthCheck(c.Request.Context(), mongoClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, NewHealthCheckRes(serverOkMsg, mongoMsg))
			return
		}
		c.JSON(http.StatusOK, NewHealthCheckRes(serverOkMsg, mongoMsg))
	}
}

type (
	HealthCheckRes struct {
		Data *HealthData `json:"data"`
	}
	HealthData struct {
		Status     string            `json:"status"`
		Components *HealthComponents `json:"components"`
	}
	HealthComponents struct {
		MongoDB string `json:"mongoDB"`
	}
)

func NewHealthCheckRes(serverMsg, mongoMsg string) *HealthCheckRes {
	return &HealthCheckRes{
		Data: &HealthData{
			Status: serverMsg,
			Components: &HealthComponents{
				MongoDB: mongoMsg,
			},
		},
	}
}
