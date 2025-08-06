//go:build integration
package integration

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	objectID := primitive.NewObjectID()
	now := time.Now()

	partDoc := bson.M{
		"_id":            objectID,
		"name":           "Test Engine Part",
		"description":    "A test engine part for e2e testing",
		"price":          150000.00,
		"stock_quantity": 10,
		"category":       1, // CATEGORY_ENGINE
		"dimensions": bson.M{
			"length": 100.0,
			"width":  50.0,
			"height": 75.0,
			"weight": 200.0,
		},
		"manufacturer": bson.M{
			"name":    "Test Manufacturer",
			"country": "USA",
			"website": "https://test.example.com",
		},
		"tags": []string{"engine", "test", "e2e"},
		"metadata": bson.M{
			"power_output":    bson.M{"double_value": 8.5},
			"is_experimental": bson.M{"bool_value": false},
		},
		"created_at": primitive.NewDateTimeFromTime(now),
		"updated_at": primitive.NewDateTimeFromTime(now),
	}

	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return objectID.Hex(), nil
}

func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
