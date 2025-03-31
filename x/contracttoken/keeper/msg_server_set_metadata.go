package keeper

import (
	"context"

	"github.com/gluon-zone/gluon/x/contracttoken/types"
)

func (m msgServer) SetMetadata(ctx context.Context, msg *types.MsgSetMetadata) (*types.MsgSetMetadataResponse, error) {
	return &types.MsgSetMetadataResponse{}, nil
}
