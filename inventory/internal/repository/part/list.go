package part

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
	repoConverter "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/converter"
	repoModel "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/model"
)

func (r *repository) ListParts(ctx context.Context) ([]model.Part, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Printf("Ошибка закрытия курсора: %v\n", cerr)
		}
	}()

	var repoParts []repoModel.Part
	if err = cursor.All(ctx, &repoParts); err != nil {
		return nil, err
	}

	// Если нет документов, возвращаем пустой слайс
	if len(repoParts) == 0 {
		return []model.Part{}, nil
	}

	parts := make([]model.Part, 0, len(repoParts))
	for _, repoPart := range repoParts {
		parts = append(parts, repoConverter.ToModelPart(repoPart))
	}

	return parts, nil
}
