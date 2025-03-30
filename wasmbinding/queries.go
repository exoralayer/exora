package wasmbinding

import (
	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkquery "github.com/cosmos/cosmos-sdk/types/query"

	"gluon/wasmbinding/bindings"

	contractmanagertypes "github.com/neutron-org/neutron/v5/x/contractmanager/types"
	icatypes "github.com/neutron-org/neutron/v5/x/interchaintxs/types"
)

func (qp *QueryPlugin) GetFailures(ctx sdk.Context, address string, pagination *sdkquery.PageRequest) (*bindings.FailuresResponse, error) {
	res, err := qp.contractmanagerQueryServer.AddressFailures(ctx, &contractmanagertypes.QueryFailuresRequest{
		Address:    address,
		Pagination: pagination,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get failures for address: %s", address)
	}

	return &bindings.FailuresResponse{Failures: res.Failures}, nil
}

func (qp *QueryPlugin) GetMinIbcFee(ctx sdk.Context, _ *bindings.QueryMinIbcFeeRequest) (*bindings.QueryMinIbcFeeResponse, error) {
	fee := qp.feeRefunderKeeper.GetMinFee(ctx)
	return &bindings.QueryMinIbcFeeResponse{MinFee: fee}, nil
}

func (qp *QueryPlugin) GetInterchainAccountAddress(ctx sdk.Context, req *bindings.QueryInterchainAccountAddressRequest) (*bindings.QueryInterchainAccountAddressResponse, error) {
	grpcReq := icatypes.QueryInterchainAccountAddressRequest{
		OwnerAddress:        req.OwnerAddress,
		InterchainAccountId: req.InterchainAccountID,
		ConnectionId:        req.ConnectionID,
	}

	grpcResp, err := qp.icaControllerKeeper.InterchainAccountAddress(ctx, &grpcReq)
	if err != nil {
		return nil, err
	}

	return &bindings.QueryInterchainAccountAddressResponse{InterchainAccountAddress: grpcResp.GetInterchainAccountAddress()}, nil
}

// GetDenomAdmin is a query to get denom admin.
func (qp QueryPlugin) GetDenomAdmin(ctx sdk.Context, denom string) (*bindings.DenomAdminResponse, error) {
	metadata, err := qp.tokenFactoryKeeper.GetAuthorityMetadata(ctx, denom)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get admin for denom: %s", denom)
	}

	return &bindings.DenomAdminResponse{Admin: metadata.Admin}, nil
}

// GetBeforeSendHook is a query to get denom before send hook.
func (qp QueryPlugin) GetBeforeSendHook(ctx sdk.Context, denom string) (*bindings.BeforeSendHookResponse, error) {
	contractAddr := qp.tokenFactoryKeeper.GetBeforeSendHook(ctx, denom)

	return &bindings.BeforeSendHookResponse{ContractAddr: contractAddr}, nil
}
