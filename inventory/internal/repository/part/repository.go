package part

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	def "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(ctx context.Context, db *mongo.Database) *repository {
	collection := db.Collection("parts")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
	}

	indexCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(indexCtx, indexModels)
	if err != nil {
		panic(err)
	}
	r := &repository{
		collection: collection,
	}

	err = r.initParts(ctx, 100)
	if err != nil {
		log.Printf("failed to initialize parts: %v", err)
		return nil
	}
	return r
}
