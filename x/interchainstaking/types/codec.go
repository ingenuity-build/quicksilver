package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgRegisterZone{}, "cosmos-sdk/MsgRegisterZone", nil)
	cdc.RegisterConcrete(&MsgSignalIntent{}, "cosmos-sdk/MsgSignalIntent", nil)
	cdc.RegisterConcrete(&RegisterZoneProposal{}, "cosmos-sdk/RegisterZoneProposal", nil)
	cdc.RegisterConcrete(&UpdateZoneProposal{}, "cosmos-sdk/UpdateZoneProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	// cosmos.base.v1beta1.Msg
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgRegisterZone{},
		&MsgSignalIntent{},
	)

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&UpdateZoneProposal{},
		&RegisterZoneProposal{},
	)
	// registry.RegisterImplementations(
	// 	(*authz.Authorization)(nil),
	// 	&SendAuthorization{},
	// )

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
func init() {
	cryptocodec.RegisterCrypto(amino)
	govtypes.RegisterProposalType(ProposalTypeRegisterZone)
	govtypes.RegisterProposalTypeCodec(&RegisterZoneProposal{}, "cosmos-sdk/RegisterZoneProposal")

	govtypes.RegisterProposalType(ProposalTypeUpdateZone)
	govtypes.RegisterProposalTypeCodec(&UpdateZoneProposal{}, "cosmos-sdk/UpdateZoneProposal")
}
