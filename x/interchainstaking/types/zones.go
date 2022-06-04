package types

import (
	"encoding/base64"
	fmt "fmt"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ingenuity-build/quicksilver/utils"
)

func (z RegisteredZone) GetDelegationAccountsByLowestBalance(qty uint64) []*ICAAccount {
	delegationAccounts := z.GetDelegationAccounts()
	sort.SliceStable(delegationAccounts, func(i, j int) bool {
		return delegationAccounts[i].DelegatedBalance.Amount.GT(delegationAccounts[j].DelegatedBalance.Amount)
	})
	if qty > 0 {
		return delegationAccounts[:int(utils.MinU64(append([]uint64{}, uint64(len(delegationAccounts)), qty)))]
	}
	return delegationAccounts
}

func (z RegisteredZone) SupportMultiSend() bool { return z.MultiSend }

func (z *RegisteredZone) GetValidatorByValoper(valoper string) (*Validator, error) {
	for _, v := range z.GetValidatorsSorted() {
		if v.ValoperAddress == valoper {
			return v, nil
		}
	}
	return nil, fmt.Errorf("invalid validator %s", valoper)
}

func (z *RegisteredZone) GetDelegationAccountByAddress(address string) (*ICAAccount, error) {
	if z.DelegationAddresses == nil {
		return nil, fmt.Errorf("no delegation accounts set: %v", z)
	}
	for _, account := range z.GetDelegationAccounts() {
		if account.GetAddress() == address {
			return account, nil
		}
	}
	return nil, fmt.Errorf("unable to find delegation account: %s", address)
}

func (z *RegisteredZone) ValidateCoinsForZone(ctx sdk.Context, coins sdk.Coins) error {

	zoneVals := z.GetValidatorsAddressesAsSlice()
COINS:
	for _, coin := range coins {
		if coin.Denom == z.BaseDenom {
			continue
		}

		for _, v := range zoneVals {
			if strings.HasPrefix(coin.Denom, v) {
				// continue 2 levels
				continue COINS
			}
		}
		return fmt.Errorf("invalid denom for zone: %s", coin.Denom)

	}
	return nil
}

func (z *RegisteredZone) ConvertCoinsToOrdinalIntents(coins sdk.Coins) ValidatorIntents {
	// should we be return DelegatorIntent here?
	out := make(ValidatorIntents)
	zoneVals := z.GetValidatorsAddressesAsSlice()
COINS:
	for _, coin := range coins {
		for _, v := range zoneVals {
			// if token share, add amount to
			if strings.HasPrefix(coin.Denom, v) {
				val, ok := out[v]
				if !ok {
					val = &ValidatorIntent{ValoperAddress: v, Weight: sdk.ZeroDec()}
				}
				val.Weight = val.Weight.Add(sdk.NewDecFromInt(coin.Amount))
				out[v] = val
				continue COINS
			}
		}
	}

	return out
}

func (z *RegisteredZone) ConvertMemoToOrdinalIntents(coins sdk.Coins, memo string) ValidatorIntents {
	// should we be return DelegatorIntent here?
	out := make(ValidatorIntents)

	if len(memo) == 0 {
		return out
	}

	memoBytes, err := base64.StdEncoding.DecodeString(memo)
	if err != nil {
		fmt.Println("Failed to decode base64 memo", err)
		return out
	}

	if len(memoBytes)%21 != 0 { // memo must be one byte (1-200) weight then 20 byte valoperAddress
		fmt.Println("Message was incorrect length", len(memoBytes))
		return out
	}

	for index := 0; index < len(memoBytes); {
		sdkWeight := sdk.NewDecFromInt(sdk.NewInt(int64(memoBytes[index]))).QuoInt(sdk.NewInt(200))
		coinWeight := sdkWeight.MulInt(coins.AmountOf(z.BaseDenom))
		index++
		address := memoBytes[index : index+20]
		index += 20
		valAddr, _ := bech32.ConvertAndEncode(z.AccountPrefix+"valoper", address)

		val, ok := out[valAddr]
		if !ok {
			val = &ValidatorIntent{ValoperAddress: valAddr, Weight: sdk.ZeroDec()}
		}
		val.Weight = val.Weight.Add(coinWeight)
		out[valAddr] = val
	}
	return out
}

func (z *RegisteredZone) GetValidatorsSorted() []*Validator {
	sort.Slice(z.Validators, func(i, j int) bool {
		return z.Validators[i].ValoperAddress < z.Validators[j].ValoperAddress
	})
	return z.Validators
}

func (z RegisteredZone) GetValidatorsAddressesAsSlice() []string {
	l := make([]string, 0)
	for _, v := range z.Validators {
		l = append(l, v.ValoperAddress)
	}

	sort.Strings(l)

	return l
}

func (z *RegisteredZone) GetDelegatedAmount() sdk.Coin {
	out := sdk.NewCoin(z.BaseDenom, sdk.ZeroInt())
	for _, da := range z.GetDelegationAccounts() {
		out = out.Add(da.DelegatedBalance)
	}
	return out
}

func (z *RegisteredZone) GetDelegationAccounts() []*ICAAccount {
	delegationAccounts := z.DelegationAddresses
	sort.Slice(delegationAccounts, func(i, j int) bool {
		return delegationAccounts[i].Address < delegationAccounts[j].Address
	})
	return delegationAccounts
}

func (z *RegisteredZone) GetAggregateIntentOrDefault() ValidatorIntents {
	if len(z.AggregateIntent) == 0 {
		return z.DefaultAggregateIntents()
	} else {
		return z.AggregateIntent
	}
}

// defaultAggregateIntents determines the default aggregate intent (for epoch 0)
func (z *RegisteredZone) DefaultAggregateIntents() ValidatorIntents {
	out := make(ValidatorIntents)
	for _, val := range z.GetValidatorsSorted() {
		if val.CommissionRate.LTE(sdk.NewDecWithPrec(5, 1)) { // 50%; make this a param.
			out[val.GetValoperAddress()] = &ValidatorIntent{ValoperAddress: val.GetValoperAddress(), Weight: sdk.OneDec()}
		}
	}

	valCount := sdk.NewInt(int64(len(out)))

	// normalise the array (divide everything by length of intent list)
	for _, key := range out.Keys() {
		if val, ok := out[key]; ok {
			val.Weight = val.Weight.Quo(sdk.NewDecFromInt(valCount))
			out[key] = val
		}
	}

	return out
}
