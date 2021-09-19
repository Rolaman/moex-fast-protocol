package structs

import "github.com/shopspring/decimal"

type SecurityType string

const (
	Future   SecurityType = "Future"
	Option   SecurityType = "Option"
	Stock    SecurityType = "Stock"
	Currency SecurityType = "Currency"
)

type Book struct {
	Id   string
	Type SecurityType
	Ask  map[decimal.Decimal]*BookRow
	Bid  map[decimal.Decimal]*BookRow
}

type BookRow struct {
	Quantity decimal.Decimal
	Price    decimal.Decimal
}

type SecurityId struct {
	Symbol    string
	SessionId string
}

type SortedBook struct {
	Id   string
	Type SecurityType
	Ask  []*SortedBookRow
	Bid  []*SortedBookRow
}

type SortedBookRow struct {
	Quantity decimal.Decimal
	Price    decimal.Decimal
}
