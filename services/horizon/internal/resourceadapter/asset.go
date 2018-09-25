package resourceadapter

import (
	"context"

	"github.com/quantadex/stellar_go/xdr"
	. "github.com/quantadex/stellar_go/protocols/horizon"

)

func PopulateAsset(ctx context.Context, dest *Asset, asset xdr.Asset) error {
	return asset.Extract(&dest.Type, &dest.Code, &dest.Issuer)
}
