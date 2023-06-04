package infrastructure

import (
	"context"
	"log"

	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

const collection = "weathers"

type weatherLocalRepository struct {
	client *mongo.Client
	dbName string
	logger *log.Logger
}

func NewWeatherLocalRepository(client *mongo.Client, dbName string, logger *log.Logger) domain.WeatherLocalRepository {
	return &weatherLocalRepository{
		client: client,
		dbName: dbName,
	}
}

func (r *weatherLocalRepository) Save(weather *domain.Weather) error {
	_, err := r.client.Database(r.dbName).Collection(collection).InsertOne(context.Background(), weather)
	return err
}
