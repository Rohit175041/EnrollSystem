package storage

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
)

const dbName = "studentsdb"
const studentCollection = "students"

func GetCollection() *mongo.Collection {
	return Client.Database(dbName).Collection(studentCollection)
}

func Init(mongoURI string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}

	// Ping to verify connection
	if err = client.Ping(ctx, nil); err != nil {
		return err
	}

	Client = client
	slog.Info("✅ Connected to MongoDB")
	return nil
}

// Disconnect closes the mongo client connection gracefully
func Disconnect() {
	if Client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(ctx); err != nil {
		slog.Error("failed to disconnect MongoDB client", slog.String("error", err.Error()))
	} else {
		slog.Info("MongoDB client disconnected")
	}
}
