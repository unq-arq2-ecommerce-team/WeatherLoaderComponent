package mongo

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func Connect(ctx context.Context, baseLogger domain.Logger, uri, database string) *mongo.Database {
	log := baseLogger.WithFields(domain.LoggerFields{"logger": "mongo", "database": database})
	ctx, cf := context.WithTimeout(ctx, 10*time.Second)
	defer cf()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Fatalf("an error has occurred while trying to connect to mongo cluster")
	}

	// check connection
	if err = client.Ping(ctx, nil); err != nil {
		log.WithFields(domain.LoggerFields{"error": err}).Fatalf("could not connect to mongo cluster")
		panic(err)
	}

	log.Info("successfully connected to mongo cluster")
	return client.Database(database)
}
