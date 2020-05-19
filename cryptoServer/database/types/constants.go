package types

// Sell order type
const (
	SellOrder = iota + 1
	BuyOrder
)

// BTCXchangeRate  exchange rate 1 BTC 10$
const BTCXchangeRate float64 = 10

// EmptyString  string value to be used where no string value is needed
const EmptyString string = ""

// API operations
const (
	Register = iota + 1
	ListWalletBalances
	ListOrderBook
	PlaceOrder
	ListYourOrders
	CancelOrder
)

// OrderExecutedSuccessfully successful order execution message
const OrderExecutedSuccessfully string = "Your orer has been places sucessfully "

// OrderHold no matching offer for this order was found
const OrderHold string = "Your orer has been placed on hold see list of your orders"
