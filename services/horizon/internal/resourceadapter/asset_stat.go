package resourceadapter

import (
	"context"

	"github.com/quantadex/stellar_go/amount"
	. "github.com/quantadex/stellar_go/protocols/horizon"
	"github.com/quantadex/stellar_go/services/horizon/internal/db2/assets"
	"github.com/quantadex/stellar_go/support/errors"
	"github.com/quantadex/stellar_go/support/render/hal"
	"github.com/quantadex/stellar_go/xdr"
)

// PopulateAssetStat fills out the details
//func PopulateAssetStat(
func PopulateAssetStat(
	ctx context.Context,
	res *AssetStat,
	row assets.AssetStatsR,
) (err error) {
	res.Asset.Type = row.Type
	res.Asset.Code = row.Code
	res.Asset.Issuer = row.Issuer
	res.Amount, err = amount.IntStringToAmount(row.Amount)
	if err != nil {
		return errors.Wrap(err, "Invalid amount in PopulateAssetStat")
	}
	res.NumAccounts = row.NumAccounts
	res.Flags = AccountFlags{
		(row.Flags & int8(xdr.AccountFlagsAuthRequiredFlag)) != 0,
		(row.Flags & int8(xdr.AccountFlagsAuthRevocableFlag)) != 0,
	}
	res.PT = row.SortKey

	res.Links.Toml = hal.NewLink(row.Toml)
	return
}
