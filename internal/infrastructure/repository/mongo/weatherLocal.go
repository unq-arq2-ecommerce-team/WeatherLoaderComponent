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
	log := r.logger.WithRequestId(ctx).WithFields(domain.LoggerFields{"city": city})
	log.Debugf("init find current by city...")
	ctxTimeout, cf := context.WithTimeout(ctx, r.timeout)
	defer cf()

	opts := options.Find().SetSort(bson.D{{"timestamp", -1}}).SetLimit(1)
	cursor, err := r.db.Collection(weatherCollection).Find(ctxTimeout, bson.M{"city": createStringCaseInsensitiveFilter(city)}, opts)
	if err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("error when execute find")
		return nil, err
	}

	var results []domain.Weather
	if err = cursor.All(ctxTimeout, &results); err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("error when read cursor")
		return nil, err
	}

	log.WithFields(domain.LoggerFields{"error": err}).Infof("successful find current weather by city")
	if len(results) == 0 {
		return nil, domain.WeatherNotFoundError{City: city}
	}
	return &results[0], nil
}

func (r *weatherLocalRepository) GetAverageTemperatureByCityAndDateRange(ctx context.Context, city string, dateFrom, dateTo time.Time) (*domain.AverageTemperature, error) {
	log := r.logger.WithRequestId(ctx).WithFields(domain.LoggerFields{"city": city, "dateFrom": dateFrom, "dateTo": dateTo})
	log.Debugf("init GetAverageTemperatureByCityAndDateRange...")
	ctxTimeout, cancelFn := context.WithTimeout(ctx, r.timeout)
	defer cancelFn()

	pipeline := getAverageTemperatureByCityAndDatesPipeline(city, dateFrom, dateTo)

	log = log.WithFields(domain.LoggerFields{"pipelineParam": pipeline})
	cursor, err := r.db.Collection(weatherCollection).Aggregate(ctxTimeout, pipeline)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.NewAverageTemperatureNotFoundError(city, dateFrom, dateTo)
		}
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("error when execute aggregate")
		return nil, err
	}

	res := make([]domain.AverageTemperature, 0)
	for cursor.Next(ctxTimeout) {
		var elem domain.AverageTemperature
		if err := cursor.Decode(&elem); err != nil {
			log.WithFields(domain.LoggerFields{"error": err}).Errorf("error when decode to AverageTemperature")
			return nil, err
		}
		res = append(res, elem)
	}

	if len(res) == 0 {
		log.Errorf("error when get average temperature")
		return nil, domain.NewAverageTemperatureNotFoundError(city, dateFrom, dateTo)
	}

	avgTemp := &res[0]
	avgTemp.Set(dateFrom, dateTo)

	log.Infof("successful get average temperature city %s", avgTemp)
	return avgTemp, nil
}

func (r *weatherLocalRepository) FindAllWeathers(ctx context.Context) (*[]domain.Weather, error) {
	log := r.logger.WithRequestId(ctx)
	log.Debugf("init find all weathers...")
	ctxTimeout, cf := context.WithTimeout(ctx, r.timeout)
	defer cf()

	opts := options.Find().SetSort(bson.D{{"timestamp", -1}}).SetLimit(1)
	cursor, err := r.db.Collection(weatherCollection).Find(ctxTimeout, bson.M{}, opts)
	if err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("error when execute find")
		return nil, err
	}

	var results []domain.Weather
	if err = cursor.All(ctxTimeout, &results); err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Errorf("error when read cursor")
		return nil, err
	}

	log.Infof("successful find all weathers")

	return &results, nil
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
