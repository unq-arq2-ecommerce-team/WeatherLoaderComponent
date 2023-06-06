package mongo

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const weatherCollection = "weathers"

type weatherLocalRepository struct {
	logger  domain.Logger
	db      *mongo.Database
	timeout time.Duration
}

func NewWeatherLocalRepository(client *mongo.Database, logger domain.Logger, timeout time.Duration) domain.WeatherLocalRepository {
	repo := &weatherLocalRepository{
		db:      client,
		logger:  logger.WithFields(domain.LoggerFields{"mongo.repository": "weatherLocalRepository", "collection": weatherCollection}),
		timeout: timeout,
	}
	repo.createIndexes()
	return repo
}

func (r *weatherLocalRepository) Save(ctx context.Context, weather *domain.Weather) error {
	ctxTimeout, cf := context.WithTimeout(ctx, r.timeout)
	defer cf()
	_, err := r.db.Collection(weatherCollection).InsertOne(ctxTimeout, weather)
	return err
}

func (r *weatherLocalRepository) FindCurrentByCity(ctx context.Context, city string) (*domain.Weather, error) {
	logger := r.logger.WithFields(domain.LoggerFields{"city": city})
	logger.Debugf("init find current by city...")
	ctxTimeout, cf := context.WithTimeout(ctx, r.timeout)
	defer cf()

	opts := options.Find().SetSort(bson.D{{"timestamp", -1}})
	cursor, err := r.db.Collection(weatherCollection).Find(ctxTimeout, createStringCaseInsensitiveFilter(city), opts)
	if err != nil {
		logger.WithFields(domain.LoggerFields{"error": err}).Errorf("error when execute find")
		return nil, err
	}

	var results []domain.Weather
	if err = cursor.All(ctxTimeout, &results); err != nil {
		logger.WithFields(domain.LoggerFields{"error": err}).Errorf("error when read cursor")
		return nil, err
	}

	logger.WithFields(domain.LoggerFields{"error": err}).Infof("successful find current weather by city")
	if len(results) == 0 {
		return nil, domain.WeatherNotFoundError{City: city}
	}
	return &results[0], nil
}

func (r *weatherLocalRepository) createIndexes() {
	ctxTimeout, cf := context.WithTimeout(context.Background(), r.timeout)
	defer cf()
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{"city", 1}},
		},
		{
			Keys: bson.D{{"timestamp", -1}},
		},
	}
	_, err := r.db.Collection(weatherCollection).Indexes().CreateMany(ctxTimeout, indexes)
	if err != nil {
		r.logger.WithFields(domain.LoggerFields{"error": err}).Fatalf("could not create mongo indexes")
	} else {
		r.logger.Infof("mongo indexes created")
	}
}
