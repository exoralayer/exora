package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gluon-zone/gluon/x/contracttoken/keeper"
	"github.com/gluon-zone/gluon/x/contracttoken/types"
	"github.com/stretchr/testify/require"
)

func TestMint(t *testing.T) {
	// Setup test keeper
	fixture := initFixture(t)
	msgServer := keeper.NewMsgServerImpl(fixture.keeper)

	// Generate test addresses
	contractAddr := sdk.AccAddress("test-contract-addr")
	recipientAddr := sdk.AccAddress("test-recipient-addr")

	// Register token
	token := types.Token{
		ContractAddress:       contractAddr.String(),
		BeforeSendHookEnabled: false,
	}
	err := fixture.keeper.Tokens.Set(fixture.ctx, contractAddr, token)
	require.NoError(t, err)

	// Set up mock expectations
	denom, err := types.GetTokenDenom(contractAddr.String())
	require.NoError(t, err)
	coin := sdk.NewCoin(denom, math.NewInt(100))
	coins := sdk.NewCoins(coin)

	fixture.mockBank.EXPECT().
		MintCoins(fixture.ctx, types.ModuleName, coins).
		Return(nil)

	fixture.mockBank.EXPECT().
		SendCoinsFromModuleToAccount(fixture.ctx, types.ModuleName, recipientAddr, coins).
		Return(nil)

	tests := []struct {
		name          string
		msg           *types.MsgMint
		expectError   bool
		errorContains string
	}{
		{
			name: "success: mint token",
			msg: &types.MsgMint{
				ContractAddress: contractAddr.String(),
				Amount:          math.NewInt(100),
				Recipient:       recipientAddr.String(),
			},
			expectError: false,
		},
		{
			name: "error: non-existent contract address",
			msg: &types.MsgMint{
				ContractAddress: sdk.AccAddress("non-existent-addr").String(),
				Amount:          math.NewInt(100),
				Recipient:       recipientAddr.String(),
			},
			expectError:   true,
			errorContains: "token does not exist",
		},
		{
			name: "error: invalid recipient address",
			msg: &types.MsgMint{
				ContractAddress: contractAddr.String(),
				Amount:          math.NewInt(100),
				Recipient:       "invalid-address",
			},
			expectError:   true,
			errorContains: "invalid recipient address",
		},
		{
			name: "error: invalid contract address",
			msg: &types.MsgMint{
				ContractAddress: "invalid-address",
				Amount:          math.NewInt(100),
				Recipient:       recipientAddr.String(),
			},
			expectError:   true,
			errorContains: "invalid contract address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := msgServer.Mint(fixture.ctx, tt.msg)
			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorContains)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
		})
	}
}
