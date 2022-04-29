package keeper_test

import (
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	simapp "github.com/ingenuity-build/quicksilver/app"
	"github.com/ingenuity-build/quicksilver/x/mint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// returns context and an app with updated mint keeper
func createTestApp(isCheckTx bool) (*simapp.OsmosisApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.MintKeeper.SetParams(ctx, types.DefaultParams())
	app.MintKeeper.SetMinter(ctx, types.DefaultInitialMinter())

	return app, ctx
}
