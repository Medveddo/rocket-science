package part

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Medveddo/rocket-science/inventory/internal/model"
	"github.com/Medveddo/rocket-science/inventory/internal/repository/converter"
	repoModel "github.com/Medveddo/rocket-science/inventory/internal/repository/model"
)

func (r *partsRepository) ListParts(ctx context.Context, filter *model.PartsFilter) ([]model.Part, error) {
	mongoFilter := bson.M{}
	if filter != nil {
		if len(filter.UUIDs) > 0 {
			objIDs := make([]primitive.ObjectID, 0, len(filter.UUIDs))
			for _, id := range filter.UUIDs {
				objID, err := primitive.ObjectIDFromHex(id)
				if err == nil {
					objIDs = append(objIDs, objID)
				}
			}
			if len(objIDs) > 0 {
				mongoFilter["_id"] = bson.M{"$in": objIDs}
			}
		}
		if len(filter.Names) > 0 {
			mongoFilter["name"] = bson.M{"$in": filter.Names}
		}
		if len(filter.Categories) > 0 {
			mongoFilter["category"] = bson.M{"$in": filter.Categories}
		}
		if len(filter.ManufacturerCountries) > 0 {
			mongoFilter["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
		}
		if len(filter.Tags) > 0 {
			mongoFilter["tags"] = bson.M{"$all": filter.Tags}
		}
	}

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("error while closing list parts curstor: %v\n", err)
		}
	}()

	var results []model.Part
	for cursor.Next(ctx) {
		var part repoModel.Part
		if err := cursor.Decode(&part); err != nil {
			return nil, err
		}
		results = append(results, converter.RepoPartToPart(part))
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
