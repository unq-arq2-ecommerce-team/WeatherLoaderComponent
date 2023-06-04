package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	app "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/application"
	infra "github.com/unq-arq2-ecommerce-team/WeatherLoaderComponent/internal/infrastructure"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	logger := log.Default()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoDbName, mongoConn := createMongoAtlasConnection()
	defer func() {
		log.Println("Closing mongo connection")
		if err = mongoConn.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	ticker := time.NewTicker(5 * time.Minute)

	saveCurrentWeatherUseCase := createSaveCurrentWeatherUseCase(mongoConn, mongoDbName, logger)

	// Start off by calling API immediately.
	err = saveCurrentWeatherUseCase.Do()
	if err != nil {
		logger.Println("Error saving current weather:", err)
	}
	go func() {
		for range ticker.C {
			err = saveCurrentWeatherUseCase.Do()
			if err != nil {
				logger.Println("Error saving current weather:", err)
			}
		}
	}()
	app := infra.NewGinApplication()
	logger.Fatal(app.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}

func createMongoAtlasConnection() (string, *mongo.Client) {
	mongoUri := os.Getenv("MONGO_URI")
	mongoDbName := os.Getenv("MONGO_DATABASE")
	if mongoUri == "" || mongoDbName == "" {
		panic("env vars MONGO_URI or MONGO_DATABASE not found")
	}
	mongoConn, err := createConnection(mongoUri)

	if err != nil {
		panic(err)
	}
	return mongoDbName, mongoConn
}

func createSaveCurrentWeatherUseCase(mongoConn *mongo.Client, mongoDbName string, logger *log.Logger) *app.SaveCurrentWeatherUseCase {
	apiUrl := os.Getenv("API_URL")
	lat := os.Getenv("LAT")
	long := os.Getenv("LONG")
	apiKey := os.Getenv("API_KEY")
	localRepository := infra.NewWeatherLocalRepository(mongoConn, mongoDbName, logger)
	remoteRepository := infra.NewWeatherRemoteRepository(apiKey, apiUrl, lat, long, logger)
	getCurrentWeatherQuery := app.NewGetCurrentWeatherUseCase(remoteRepository)
	saveWeatherQuery := app.NewSaveWeatherCommand(localRepository)
	saveCurrentWeatherUseCase := app.NewSaveCurrentWeatherUseCase(getCurrentWeatherQuery, saveWeatherQuery)
	return saveCurrentWeatherUseCase
}
func createConnection(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return client, nil
}
