package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ingenuity-build/quicksilver/x/participationrewards/types"
)

// HandleAddProtocolDataProposal is a handler for executing a passed add protocol data proposal
func HandleAddProtocolDataProposal(ctx sdk.Context, k Keeper, p *types.AddProtocolDataProposal) error {

	// handle me.

	ctx.EventManager().EmitEvents(sdk.Events{})

	return nil
}
