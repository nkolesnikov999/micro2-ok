package part

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/converter"
	repoModel "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/model"
)

func (r *repository) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	var repoPart repoModel.Part

	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&repoPart)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrPartNotFound
		}
		return model.Part{}, err
	}

	return repoConverter.ToModelPart(repoPart), nil
}
