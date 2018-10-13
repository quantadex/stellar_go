package build

import (
	"github.com/stellar/go/xdr"
)

func Settlement(muts ...SettlementMutator) (*SettlementBuilder, error) {
	result := &SettlementBuilder{}
	err := result.Mutate(muts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type SettlementMutator interface {
	MutateSettlement(*xdr.SettlementOp) error
}

type SettlementBuilder struct {
	Settlement xdr.SettlementOp
	O          xdr.Operation
	Err        error
}

func (b *SettlementBuilder) Mutate(muts ...SettlementMutator) error {
	for _, m := range muts {
		err := m.MutateSettlement(&b.Settlement)
		if err != nil {
			return err
		}
	}

	return nil
}

func MatchedOrder(muts ...MatchedOrderMutator) (*MatchedOrderBuilder, error) {
	result := &MatchedOrderBuilder{}
	err := result.Mutate(muts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type MatchedOrderMutator interface {
	MutateMatchedOrder(*xdr.MatchedOrder) error
}

type MatchedOrderBuilder struct {
	Orders []xdr.MatchedOrder
}

func (b *MatchedOrderBuilder) Mutate(muts ...MatchedOrderMutator) error {
	for _, m := range muts {
		var o xdr.MatchedOrder
		err := m.MutateMatchedOrder(&o)
		if err != nil {
			return err
		}
		b.Orders = append(b.Orders, o)
	}

	return nil
}

// ------------------------------------------------------------
//
//   Mutator implementations
//
// ------------------------------------------------------------

func (b *MatchedOrderBuilder) MutateSettlement(op *xdr.SettlementOp) error {
	op.MatchedOrders = append(op.MatchedOrders, b.Orders...)
	return nil
}

type SettlementHash string

func (hash SettlementHash) MutateSettlement(op *xdr.SettlementOp) error {
	op.SettlementHash = string(hash)
	return nil
}

type ParentSettlementHash string

func (hash ParentSettlementHash) MutateSettlement(op *xdr.SettlementOp) error {
	op.ParentSettlementHash = string(hash)
	return nil
}

type OrderInfo struct {
	Code    string
	Issuer  string
	Amount  uint64
	Account string
}

type Order struct {
	Buyer  OrderInfo
	Seller OrderInfo
}

// MutateMatchedOrder for Asset sets the MatchedOrder's Asset field
func (m Order) MutateMatchedOrder(o *xdr.MatchedOrder) (err error) {
	o.Buyer = xdr.AccountId{}
	err = setAccountId(m.Buyer.Account, &o.Buyer)
	if err != nil {
		return err
	}
	o.Seller = xdr.AccountId{}
	err = setAccountId(m.Seller.Account, &o.Seller)
	if err != nil {
		return err
	}

	o.AmountBuy = xdr.Int64(m.Buyer.Amount)
	o.AmountSell = xdr.Int64(m.Seller.Amount)

	o.AssetBuy, err = createAlphaNumAsset(m.Buyer.Code, m.Buyer.Issuer)
	if err != nil {
		return
	}
	o.AssetSell, err = createAlphaNumAsset(m.Seller.Code, m.Seller.Issuer)
	return
}
