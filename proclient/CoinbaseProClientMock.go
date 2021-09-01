package proclient

import (
	"github.com/preichenberger/go-coinbasepro/v2"
)

type MockClient struct {
	Err           error
	HistoricRates []coinbasepro.HistoricRate
	Book          coinbasepro.Book
	SavedOrder    coinbasepro.Order
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

// Product funcs
func (c *MockClient) GetBook(product string, level int) (coinbasepro.Book, error) {
	return c.Book, c.Err
}

func (c *MockClient) GetTicker(product string) (coinbasepro.Ticker, error) {
	var ticker coinbasepro.Ticker
	return ticker, nil
}

func (c *MockClient) ListTrades(product string,
	p ...coinbasepro.ListTradesParams) *coinbasepro.Cursor {
	return &coinbasepro.Cursor{}
}

func (c *MockClient) GetProducts() ([]coinbasepro.Product, error) {
	var products []coinbasepro.Product
	return products, nil
}

func (c *MockClient) GetHistoricRates(product string, p ...coinbasepro.GetHistoricRatesParams) ([]coinbasepro.HistoricRate, error) {
	if c.Err != nil {
		return c.HistoricRates, c.Err
	}

	return c.HistoricRates, nil
}

func (c *MockClient) GetStats(product string) (coinbasepro.Stats, error) {
	var stats coinbasepro.Stats
	return stats, nil
}

// Account Funcs
func (c *MockClient) GetAccounts() ([]coinbasepro.Account, error) {
	var accounts []coinbasepro.Account
	return accounts, nil
}

func (c *MockClient) GetAccount(id string) (coinbasepro.Account, error) {
	account := coinbasepro.Account{}
	return account, nil
}

func (c *MockClient) ListAccountLedger(id string,
	p ...coinbasepro.GetAccountLedgerParams) *coinbasepro.Cursor {
	return &coinbasepro.Cursor{}
}

func (c *MockClient) ListHolds(id string, p ...coinbasepro.ListHoldsParams) *coinbasepro.Cursor {
	return &coinbasepro.Cursor{}
}

//order funcs
func (c *MockClient) CreateOrder(newOrder *coinbasepro.Order) (coinbasepro.Order, error) {
	return c.SavedOrder, c.Err
}

func (c *MockClient) CancelOrder(id string) error {
	return nil
}

func (c *MockClient) CancelAllOrders(p ...coinbasepro.CancelAllOrdersParams) ([]string, error) {
	var orderIDs []string
	return orderIDs, nil
}

func (c *MockClient) GetOrder(id string) (coinbasepro.Order, error) {
	return c.SavedOrder, c.Err
}

func (c *MockClient) ListOrders(p ...coinbasepro.ListOrdersParams) *coinbasepro.Cursor {
	return &coinbasepro.Cursor{}
}
