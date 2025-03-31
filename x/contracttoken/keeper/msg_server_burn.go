package keeper

import (
	"context"

	"github.com/gluon-zone/gluon/x/contracttoken/types"
)

func (m msgServer) Burn(ctx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	return &types.MsgBurnResponse{}, nil
}
