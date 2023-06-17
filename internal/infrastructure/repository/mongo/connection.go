package mongo

import (
	"context"
	"github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"time"
)

func Connect(ctx context.Context, baseLogger domain.Logger, uri, database string) *mongo.Database {
	log := baseLogger.WithFields(domain.LoggerFields{"logger": "mongo", "database": database})
	ctx, cf := context.WithTimeout(ctx, 10*time.Second)
	defer cf()

	mongoOptions := options.Client()
	mongoOptions.Monitor = otelmongo.NewMonitor()
	mongoOptions.ApplyURI(uri)

	client, err := mongo.Connect(ctx, mongoOptions)

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
