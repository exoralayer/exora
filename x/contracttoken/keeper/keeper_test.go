package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/core/address"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/golang/mock/gomock"

	"github.com/exoralayer/exora/x/contracttoken/keeper"
	module "github.com/exoralayer/exora/x/contracttoken/module"
	contracttokentestutil "github.com/exoralayer/exora/x/contracttoken/testutil"
	"github.com/exoralayer/exora/x/contracttoken/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
	ctrl         *gomock.Controller
	mockAuth     *contracttokentestutil.MockAuthKeeper
	mockBank     *contracttokentestutil.MockBankKeeper
	mockWasm     *contracttokentestutil.MockWasmKeeper
}

//go:generate mockgen -destination ../testutil/expected_keepers.go -package testutil github.com/exoralayer/exora/x/contracttoken/types AuthKeeper,BankKeeper,WasmKeeper

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	authority := authtypes.NewModuleAddress(types.GovModuleName)

	ctrl := gomock.NewController(t)
	mockAuth := contracttokentestutil.NewMockAuthKeeper(ctrl)
	mockBank := contracttokentestutil.NewMockBankKeeper(ctrl)
	mockWasm := contracttokentestutil.NewMockWasmKeeper(ctrl)

	k := keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		log.NewNopLogger(),
		authority.String(),
		mockAuth,
		mockBank,
		mockWasm,
	)

	// Initialize params
	if err := k.Params.Set(ctx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return &fixture{
		ctx:          ctx,
		keeper:       k,
		addressCodec: addressCodec,
		ctrl:         ctrl,
		mockAuth:     mockAuth,
		mockBank:     mockBank,
		mockWasm:     mockWasm,
	}
}
