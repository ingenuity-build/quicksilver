package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/ingenuity-build/quicksilver/x/interchainstaking/types"
)

const (
	FlagMultiSend = "multi-send"
)

// GetTxCmd returns a root CLI command handler for all x/bank transaction commands.
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Interchain staking transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetRegisterZoneTxCmd(),
		GetSignalIntentTxCmd(),
		GetRequestRedemptionTxCmd(),
	)

	return txCmd
}

// GetRegisterZoneTxCmd returns a CLI command handler for creating a MsgSend transaction.
func GetRegisterZoneTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [identifier] [connection_id] [chain_id] [local_denom] [remote_denom]",
		Short: `Register new zone with the chain.`,
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			identifier := args[0]
			connection_id := args[1]
			chain_id := args[2]
			local_denom := args[3]
			remote_denom := args[4]

			multi_send, _ := cmd.Flags().GetBool(FlagMultiSend)
			msg := types.NewMsgRegisterZone(identifier, connection_id, chain_id, local_denom, remote_denom, clientCtx.GetFromAddress(), multi_send)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().Bool(FlagMultiSend, false, "multi-send support")

	return cmd
}

// GetSignalIntentTxCmd returns a CLI command handler for signalling validator
// delegation intent.
func GetSignalIntentTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signal-intent [chain_id] [delegation_intent]",
		Short: `Signal validator delegation intent.`,
		Long: `signal validator delegation intent by providing a comma seperated string
containing a decimal weight and the bech32 validator address,
e.g. "0.3cosmosvaloper1xxxxxxxxx,0.3cosmosvaloper1yyyyyyyyy,0.4cosmosvaloper1zzzzzzzzz"`,
		Example: `signal-intent [chain_id] 0.3cosmosvaloper1xxxxxxxxx,0.3cosmosvaloper1yyyyyyyyy,0.4cosmosvaloper1zzzzzzzzz`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			chain_id := args[0]
			intents, err := types.IntentsFromString(args[1])
			if err != nil {
				return fmt.Errorf("%v, see example: %v", err, cmd.Example)
			}

			msg := types.NewMsgSignalIntent(chain_id, intents, clientCtx.GetFromAddress())

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetRegisterZoneTxCmd returns a CLI command handler for creating a MsgSend transaction.
func GetRequestRedemptionTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem [coins] [destination_address]",
		Short: `Redeem tokens.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)

			if err != nil {
				return err
			}
			coins := args[0]
			destination_address := args[1]

			msg := types.NewMsgRequestRedemption(coins, destination_address, clientCtx.GetFromAddress())

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
