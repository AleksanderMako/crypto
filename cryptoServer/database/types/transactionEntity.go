package types

type TransactionEntity struct {
	Wallet   Wallet
	Order    Order
	OrderID  string
	NewOrder bool
}

func NewTransactionEntity(wallet Wallet, order Order, OrderID string) *TransactionEntity {

	return &TransactionEntity{
		Order:   order,
		Wallet:  wallet,
		OrderID: OrderID,
	}
}
