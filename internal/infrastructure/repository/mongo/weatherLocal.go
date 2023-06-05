package mongo

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const collection = "weathers"

type weatherLocalRepository struct {
	logger  domain.Logger
	db      *mongo.Database
	timeout time.Duration
}

func NewWeatherLocalRepository(client *mongo.Database, logger domain.Logger, timeout time.Duration) domain.WeatherLocalRepository {
	repo := &weatherLocalRepository{
		db:      client,
		logger:  logger.WithFields(domain.LoggerFields{"mongo.repository": "weatherLocalRepository", "collection": collection}),
		timeout: timeout,
	}
	repo.createIndexes()
	return repo
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
	_, err := r.db.Collection(collection).Indexes().CreateMany(ctxTimeout, indexes)
	if err != nil {
		r.logger.WithFields(domain.LoggerFields{"error": err}).Fatalf("could not create mongo indexes")
	} else {
		r.logger.Infof("mongo indexes created")
	}
}

func (r *weatherLocalRepository) Save(ctx context.Context, weather *domain.Weather) error {
	ctxTimeout, cf := context.WithTimeout(ctx, r.timeout)
	defer cf()
	_, err := r.db.Collection(collection).InsertOne(ctxTimeout, weather)
	return err
}
