package keeper

import (
	"context"

	"github.com/gluon-zone/gluon/x/contracttoken/types"
)

func (m msgServer) UpdateToken(ctx context.Context, msg *types.MsgUpdateToken) (*types.MsgUpdateTokenResponse, error) {
	return &types.MsgUpdateTokenResponse{}, nil
}
