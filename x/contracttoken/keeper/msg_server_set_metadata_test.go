package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gluon-zone/gluon/x/contracttoken/keeper"
	"github.com/gluon-zone/gluon/x/contracttoken/types"
	"github.com/stretchr/testify/require"
)

func TestSetMetadata(t *testing.T) {
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
	err := fixture.keeper.ContractTokens.Set(fixture.ctx, contractAddr, token)
	require.NoError(t, err)

	// Get denom
	denom, err := types.GetTokenDenom(contractAddr.String())
	require.NoError(t, err)

	// Create metadata
	metadata := banktypes.Metadata{
		Description: "Test Token",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    denom,
				Exponent: 0,
			},
		},
		Base:    denom,
		Display: denom,
		Name:    "Test Token",
		Symbol:  "TEST",
	}

	tests := []struct {
		name          string
		msg           *types.MsgSetMetadata
		expectError   bool
		errorContains string
	}{
		{
			name: "success: set metadata",
			msg: &types.MsgSetMetadata{
				ContractAddress: contractAddr.String(),
				Metadata:        metadata,
			},
			expectError: false,
		},
		{
			name: "error: invalid contract address",
			msg: &types.MsgSetMetadata{
				ContractAddress: "invalid-address",
				Metadata:        metadata,
			},
			expectError:   true,
			errorContains: "invalid contract address",
		},
		{
			name: "error: token does not exist",
			msg: &types.MsgSetMetadata{
				ContractAddress: sdk.AccAddress("non-existent-addr").String(),
				Metadata:        metadata,
			},
			expectError:   true,
			errorContains: "token does not exist",
		},
		{
			name: "error: invalid denom",
			msg: &types.MsgSetMetadata{
				ContractAddress: contractAddr.String(),
				Metadata: banktypes.Metadata{
					Description: "Test Token",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    "invalid-denom",
							Exponent: 0,
						},
					},
					Base:    "invalid-denom",
					Display: "invalid-denom",
					Name:    "Test Token",
					Symbol:  "TEST",
				},
			},
			expectError:   true,
			errorContains: "base denom in the metadata does not match the denom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := msgServer.SetMetadata(fixture.ctx, tt.msg)
			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorContains)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)

			// Verify metadata was set
			fixture.mockBank.EXPECT().
				SetDenomMetaData(fixture.ctx, tt.msg.Metadata)
		})
	}
}
