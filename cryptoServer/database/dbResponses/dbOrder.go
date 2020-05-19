package dbResponses

import "cryptoServer/database/types"

type DBOrder struct {
	ID    string      `json:"ID"`
	Order types.Order `json:"order"`
}

func NewDBOrder(ID string, order types.Order) *DBOrder {

	return &DBOrder{
		ID:    ID,
		Order: order,
	}
}
