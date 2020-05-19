package dbResponses

import "cryptoServer/database/types"

// DBOrder contains a struct of type Order and the Struct's ID
// the ID is used as a key in a in-memory map
type DBOrder struct {
	ID    string      `json:"ID"`
	Order types.Order `json:"order"`
}

// NewDBOrder returns a pointer to DBOrder
func NewDBOrder(ID string, order types.Order) *DBOrder {

	return &DBOrder{
		ID:    ID,
		Order: order,
	}
}
