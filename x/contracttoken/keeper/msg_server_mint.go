package keeper

import (
	"context"

	"github.com/gluon-zone/gluon/x/contracttoken/types"
)

func (m msgServer) Mint(ctx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	return &types.MsgMintResponse{}, nil
}
