package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

type LiquidTokensModule struct {
}

var _ Submodule = &LiquidTokensModule{}

func (m *LiquidTokensModule) Hooks(ctx sdk.Context, k Keeper) {

}

func (m *LiquidTokensModule) IsActive() bool {
	return true
}

func (m *LiquidTokensModule) IsReady() bool {
	return true
}

func (m *LiquidTokensModule) VerifyClaim(ctx sdk.Context, k Keeper, msg types.MsgSubmitClaim) bool {
	// message
	// check denom is valid vs allowed

	epochInfo := k.epochsKeeper.GetEpochInfo(ctx, "epoch")
	epoch_id := epochInfo.CurrentEpoch - 1

	_, found := k.GetProtocolData(ctx, fmt.Sprintf("liquid/%s/%s", msg.Zone, msg.AssetType))
	if !found {
		k.Logger(ctx).Error("unable to query " + fmt.Sprintf("liquid/%s/%s", msg.Zone, msg.AssetType) + " in LiquidTokensModule hook")
		return false
	}

	// check height is valid for epoch
	_, found = k.GetProtocolData(ctx, fmt.Sprintf("epoch/%d/%s", epoch_id, msg.Zone))
	if !found {
		k.Logger(ctx).Error("unable to query " + fmt.Sprintf("epoch/%d/%s", epoch_id, msg.Zone) + " in LiquidTokensModule hook")
		return false
	}

	// check proof is valid

	return true
}
