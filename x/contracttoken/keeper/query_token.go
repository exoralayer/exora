package keeper

import (
	"context"

	"github.com/exoralayer/exora/x/contracttoken/types"
)

func (q queryServer) Token(ctx context.Context, req *types.QueryTokenRequest) (*types.QueryTokenResponse, error) {
	return &types.QueryTokenResponse{}, nil
}

func (q queryServer) Tokens(ctx context.Context, req *types.QueryTokensRequest) (*types.QueryTokensResponse, error) {
	return &types.QueryTokensResponse{}, nil
}
