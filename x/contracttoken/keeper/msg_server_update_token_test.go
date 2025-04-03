package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/exoralayer/exora/x/contracttoken/keeper"
	"github.com/exoralayer/exora/x/contracttoken/types"
	"github.com/stretchr/testify/require"
)

func TestUpdateToken(t *testing.T) {
	// Setup test keeper
	fixture := initFixture(t)
	msgServer := keeper.NewMsgServerImpl(fixture.keeper)

	// Generate test addresses
	contractAddr := sdk.AccAddress("test-contract-addr")

	// Register token
	token := types.Token{
		ContractAddress:       contractAddr.String(),
		BeforeSendHookEnabled: false,
	}
	err := fixture.keeper.Tokens.Set(fixture.ctx, contractAddr, token)
	require.NoError(t, err)

	tests := []struct {
		name          string
		msg           *types.MsgUpdateToken
		expectError   bool
		errorContains string
	}{
		{
			name: "success: update token",
			msg: &types.MsgUpdateToken{
				ContractAddress:       contractAddr.String(),
				BeforeSendHookEnabled: true,
			},
			expectError: false,
		},
		{
			name: "error: invalid contract address",
			msg: &types.MsgUpdateToken{
				ContractAddress:       "invalid-address",
				BeforeSendHookEnabled: true,
			},
			expectError:   true,
			errorContains: "invalid contract address",
		},
		{
			name: "error: token does not exist",
			msg: &types.MsgUpdateToken{
				ContractAddress:       sdk.AccAddress("non-existent-addr").String(),
				BeforeSendHookEnabled: true,
			},
			expectError:   true,
			errorContains: "token does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := msgServer.UpdateToken(fixture.ctx, tt.msg)
			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorContains)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)

			// Verify token was updated
			token, err := fixture.keeper.Tokens.Get(fixture.ctx, contractAddr)
			require.NoError(t, err)
			require.Equal(t, tt.msg.ContractAddress, token.ContractAddress)
			require.Equal(t, tt.msg.BeforeSendHookEnabled, token.BeforeSendHookEnabled)
		})
	}
}
