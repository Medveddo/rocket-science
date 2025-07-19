package part

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Medveddo/rocket-science/inventory/internal/model"
	"github.com/Medveddo/rocket-science/inventory/internal/repository/converter"
	repoModel "github.com/Medveddo/rocket-science/inventory/internal/repository/model"
)

func (r *partsRepository) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	objectID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return model.Part{}, err
	}

	var part repoModel.Part
	filter := bson.M{"_id": objectID}
	err = r.collection.FindOne(ctx, filter).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrPartDoesNotExist
		}
		return model.Part{}, err
	}

	return converter.RepoPartToPart(part), nil
}
