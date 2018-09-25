package resourceadapter

import (
	"context"

	"github.com/quantadex/stellar_go/services/horizon/internal/db2/history"
	. "github.com/quantadex/stellar_go/protocols/horizon"
)

func PopulateHistoryAccount(ctx context.Context, dest *HistoryAccount, row history.Account) {
	dest.ID = row.Address
	dest.AccountID = row.Address
}
