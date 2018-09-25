package resourceadapter

import (
	"github.com/quantadex/stellar_go/services/horizon/internal/db2/core"
	. "github.com/quantadex/stellar_go/protocols/horizon"
)

func PopulateAccountFlags(dest *AccountFlags, row core.Account) {
	dest.AuthRequired = row.IsAuthRequired()
	dest.AuthRevocable = row.IsAuthRevocable()
}
