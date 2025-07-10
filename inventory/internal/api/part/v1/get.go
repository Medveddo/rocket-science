package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Medveddo/rocket-science/inventory/internal/converter"
	"github.com/Medveddo/rocket-science/inventory/internal/model"
	inventoryV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/inventory/v1"
)

func (a *partAPI) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.partService.GetPart(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrPartDoesNotExist) {
			return nil, status.Errorf(codes.NotFound, "part not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &inventoryV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
