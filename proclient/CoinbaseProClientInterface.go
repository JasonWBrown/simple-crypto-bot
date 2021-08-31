package proclient

import "github.com/preichenberger/go-coinbasepro/v2"

type ProClientInterface interface {
	//Product
	GetBook(product string, level int) (coinbasepro.Book, error)
	GetTicker(product string) (coinbasepro.Ticker, error)
	ListTrades(product string, p ...coinbasepro.ListTradesParams) *coinbasepro.Cursor
	GetProducts() ([]coinbasepro.Product, error)
	GetHistoricRates(product string, p ...coinbasepro.GetHistoricRatesParams) ([]coinbasepro.HistoricRate, error)
	GetStats(product string) (coinbasepro.Stats, error)

	//Account
	GetAccounts() ([]coinbasepro.Account, error)
	GetAccount(id string) (coinbasepro.Account, error)
	ListAccountLedger(id string, p ...coinbasepro.GetAccountLedgerParams) *coinbasepro.Cursor
	ListHolds(id string, p ...coinbasepro.ListHoldsParams) *coinbasepro.Cursor

	//Order
	CreateOrder(newOrder *coinbasepro.Order) (coinbasepro.Order, error)
	CancelOrder(id string) error
	CancelAllOrders(p ...coinbasepro.CancelAllOrdersParams) ([]string, error)
	GetOrder(id string) (coinbasepro.Order, error)
	ListOrders(p ...coinbasepro.ListOrdersParams) *coinbasepro.Cursor
}
