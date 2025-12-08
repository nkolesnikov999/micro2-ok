package part

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	def "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/logger"
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
		logger.Error(ctx, "failed to initialize parts", zap.Error(err))
		return nil
	}
	return r
}
