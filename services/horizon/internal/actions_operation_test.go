package horizon

import (
	"encoding/json"
	"fmt"
	"github.com/stellar/go.bak/_vendor-20181016153202/github.com/go-errors/errors"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/services/horizon/internal/resource"
	"github.com/stellar/go/support/log"
	"net/http"
	//"net/url"
	"strconv"
	"testing"
	"time"

	//"time"

	//"github.com/stellar/go/services/horizon/internal/db2/history"
	"github.com/stellar/go/services/horizon/internal/resource/operations"
	"github.com/stellar/go/services/horizon/internal/test"
)

func TestOperationActions_Index(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	// no filter
	w := ht.Get("/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(4, w.Body)
	}

	// filtered by ledger sequence
	w = ht.Get("/ledgers/1/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(0, w.Body)
	}

	w = ht.Get("/ledgers/2/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(3, w.Body)
	}

	w = ht.Get("/ledgers/3/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
	}

	// filtered by account
	w = ht.Get("/accounts/GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(3, w.Body)
	}

	w = ht.Get("/accounts/GA5WBPYA5Y4WAEHXWR2UKO2UO4BUGHUQ74EUPKON2QHV4WRHOIRNKKH2/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
	}

	w = ht.Get("/accounts/GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(2, w.Body)
	}

	// filtered by transaction
	w = ht.Get("/transactions/2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
	}

	w = ht.Get("/transactions/164a5064eba64f2cdbadb856bf3448485fc626247ada3ed39cddf0f6902133b6/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
	}

	// filtered by ledger
	w = ht.Get("/ledgers/3/operations")
	if ht.Assert.Equal(200, w.Code) {
		ht.Assert.PageOf(1, w.Body)
	}

	// missing ledger
	w = ht.Get("/ledgers/100/operations")
	ht.Assert.Equal(404, w.Code)
}

func TestOperationActions_Show(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	// exists
	w := ht.Get("/operations/8589938689")
	if ht.Assert.Equal(200, w.Code) {
		var result operations.Base
		err := json.Unmarshal(w.Body.Bytes(), &result)
		ht.Require.NoError(err, "failed to parse body")
		ht.Assert.Equal("8589938689", result.PT)
		ht.Assert.Equal("2374e99349b9ef7dba9a5db3339b78fda8f34777b1af33ba468ad5c0df946d4d", result.TransactionHash)
	}

	// doesn't exist
	w = ht.Get("/operations/9589938689")
	ht.Assert.Equal(404, w.Code)

	// before history
	ht.ReapHistory(1)
	w = ht.Get("/operations/8589938689")
	ht.Assert.Equal(410, w.Code)
}

func TestOperationActions_Regressions(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	// ensure that trying to stream ops from an account that doesn't exist
	// fails before streaming the hello message.  Regression test for #285
	w := ht.Get("/accounts/foo/operations?limit=1", test.RequestHelperStreaming)
	if ht.Assert.Equal(404, w.Code) {
		ht.Assert.ProblemType(w.Body, "not_found")
	}

	// #202 - price is not shown on manage_offer operations
	test.LoadScenario("trades")
	w = ht.Get("/operations/21474840577")
	if ht.Assert.Equal(200, w.Code) {
		var result operations.ManageOffer
		err := json.Unmarshal(w.Body.Bytes(), &result)
		ht.Require.NoError(err, "failed to parse body")
		ht.Assert.Equal("1.0000000", result.Price)
		ht.Assert.Equal(int32(1), result.PriceR.N)
		ht.Assert.Equal(int32(1), result.PriceR.D)
	}
}

func TestOperation_CreatedAt(t *testing.T) {

	Horizon := &horizon.Client{
		URL:  "http://localhost:8000",
		HTTP: http.DefaultClient,
	}

	order := &build.Order{
		Buyer: build.OrderInfo{
			Code:    "ETH",
			Issuer:  "QAIVSB7SVO76Z3SHYCU4NCYJ44WLA2Z6DVJAX3ICEM4RYSPWS3KJ6XXP",
			Amount:  10,
			Account: "QDNDQ74KK5XIT7RPOB7BB7STTZHCXU6GEXLPXIMW4ICQHGV3BTLGMOV3",
		},
		Seller: build.OrderInfo{
			Code:    "BTC",
			Issuer:  "QAIVSB7SVO76Z3SHYCU4NCYJ44WLA2Z6DVJAX3ICEM4RYSPWS3KJ6XXP",
			Amount:  10,
			Account: "QCVABKNIUKXJSJMXGHPQC2ZRK6J26PKIM6DYFB4WTJ5YU5W3BKZ3BXZS",
		},
	}

	orderBuilder, err := build.MatchedOrder(order)
	if err != nil {
		log.Errorf("submit err %s", err.Error())
		return
	}
	settlementBuilder, err := build.Settlement(orderBuilder)
	if err != nil {
		log.Errorf("submit err %s", err.Error())
		return
	}

	w, err := Horizon.HTTP.Get("http://localhost:8000/accounts/QBJUMEXBXBQN6NGHVGE4U3IWLYKF6RZ2JXY67CBSLKPXNZPLJ2KZ5DWH")
	var result resource.Account
	if err == nil {
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			errors.New("failed to decode operations: " + err.Error())
		}
	}
	seq, err := strconv.ParseUint(result.Sequence, 10, 64)

	fmt.Println("Seq = ", seq)
	tx, err := build.Transaction(
		build.Network{"QUANTA Test Network ; September 2018"},
		build.SourceAccount{"QBJUMEXBXBQN6NGHVGE4U3IWLYKF6RZ2JXY67CBSLKPXNZPLJ2KZ5DWH"},
		build.Sequence{seq + 1},
		settlementBuilder,
	)
	if err != nil {
		log.Errorf("submit err %s", err.Error())
		return
	}

	sign, err := tx.Sign("ZCQ4DDYJMXPLAN5ZUAQ5UXH325BRLOPI2326XUYU3LDVQZEIBIJJZCGA")
	if err != nil {
		log.Errorf("submit err %s", err.Error())
		return
	}

	signed, err := sign.Base64()
	if err != nil {
		log.Errorf("submit err %s", err.Error())
		return
	}

	fmt.Println(signed)
	Horizon.SubmitTransaction(signed)

	time.Sleep(time.Second * 4)

	//w = ht.Post("/transactions", form)
	//w := ht.Get("/ledgers/3/operations")
	//records := []operations.Base{}
	//ht.UnmarshalPage(w.Body, &records)

	//l := history.Ledger{}
	//hq := history.Q{Session: ht.HorizonSession()}
	//ht.Require.NoError(hq.LedgerBySequence(&l, 3))

	//ht.Assert.WithinDuration(l.ClosedAt, records[0].LedgerCloseTime, 1*time.Second)
}
