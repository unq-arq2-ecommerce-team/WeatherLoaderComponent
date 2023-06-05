package mongo

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	loggerPkg "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

const collection = "weathers"

type weatherLocalRepository struct {
	db         *mongo.Database
	baseLogger domain.Logger
}

func NewWeatherLocalRepository(client *mongo.Database, baseLogger domain.Logger) domain.WeatherLocalRepository {
	return &weatherLocalRepository{
		db:         client,
		baseLogger: baseLogger.WithFields(loggerPkg.Fields{"mongo.repository": "weatherLocalRepository"}),
	}
}

func (r *weatherLocalRepository) Save(weather *domain.Weather) error {
	_, err := r.db.Collection(collection).InsertOne(context.Background(), weather)
	return err
}
