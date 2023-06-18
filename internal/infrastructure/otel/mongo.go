package otel

import (
	"go.mongodb.org/mongo-driver/event"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

func GetMongoMonitor() *event.CommandMonitor {
	return otelmongo.NewMonitor()
}
