package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

var _ types.QueryServer = Keeper{}

// RegisteredZoneInfos provides information about registered zones
func (k Keeper) RegisteredZoneInfos(c context.Context, req *types.QueryRegisteredZonesInfoRequest) (*types.QueryRegisteredZonesInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var zones []types.RegisteredZone
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixZone)

	pageRes, err := query.Paginate(store, req.Pagination, func(_, value []byte) error {
		var zone types.RegisteredZone
		if err := k.cdc.Unmarshal(value, &zone); err != nil {
			return err
		}
		zones = append(zones, zone)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRegisteredZonesInfoResponse{
		Zones:      zones,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) DepositAccountFromAddress(c context.Context, req *types.QueryDepositAccountForChainRequest) (*types.QueryDepositAccountForChainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	zone, found := k.GetRegisteredZoneInfo(ctx, req.GetChainId())
	if !found {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("no zone found matching %s", req.GetChainId()))
	}

	return &types.QueryDepositAccountForChainResponse{
		DepositAccountAddress: zone.DepositAddress.Address,
	}, nil
}
