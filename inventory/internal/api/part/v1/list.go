package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Medveddo/rocket-science/inventory/internal/converter"
	inventoryV1 "github.com/Medveddo/rocket-science/shared/pkg/proto/inventory/v1"
)

func (a *partAPI) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter := converter.ListPartsFilterFromProto(req.Filter)
	parts, err := a.partService.ListParts(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	resp := &inventoryV1.ListPartsResponse{
		Parts: make([]*inventoryV1.Part, 0, len(parts)),
	}
	for _, part := range parts {
		resp.Parts = append(resp.Parts, converter.PartToProto(part))
	}
	return resp, nil
}
