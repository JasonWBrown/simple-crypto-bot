package svc

import (
	"fmt"
	"strconv"

	"github.com/preichenberger/go-coinbasepro/v2"
	"github.com/shopspring/decimal"
)

type CoinbaseSvcInterface interface {
	Sell(product string, numberOwn, sellPrice float64) (float64, float64)
	Buy(product string, buyPrice, availablefunds float64) (float64, float64)
	GetLastPrice(product string, ePrice float64) decimal.Decimal
}

type CoinbaseSvc struct {
	Client *coinbasepro.Client
}

func NewCoinbaseSvc(client *coinbasepro.Client) CoinbaseSvc {
	return CoinbaseSvc{
		Client: client,
	}
}

func (svc CoinbaseSvc) Sell(product string, numberOwn, sellPrice float64) (float64, float64) {
	savedOrder, err := svc.Client.CreateOrder(&coinbasepro.Order{
		ProductID: product,
		Side:      "sell",
		Size:      fmt.Sprintf("%f", numberOwn),
	})
	if err != nil {
		fmt.Printf("Failed to sell %s\n", err.Error())
	}

	//TODO back off here?
	savedOrder, err = svc.Client.GetOrder(savedOrder.ID)
	if err != nil {
		fmt.Printf("Failed to get order %s\n", err.Error())
	}

	//TODO how do I get my funds available and how do I know the order completed.
	accounts, err := svc.Client.GetAccounts()
	if err != nil {
		fmt.Printf("Failed to get accounts %s\n", err.Error())
	}
	var account coinbasepro.Account
	// these might be in order might not have to iterate every single time
	for _, a := range accounts {
		if a.Currency == "USD" {
			account = a
			break
		}
	}

	funds, err := strconv.ParseFloat(account.Balance, 64)
	if err != nil {
		fmt.Printf("Failed to parse float for account balance %s\n", err.Error())
	}

	return 0.0, funds
}

func (svc CoinbaseSvc) Buy(product string, buyPrice, availablefunds float64) (float64, float64) {
	orderAmount := availablefunds / buyPrice
	savedOrder, err := svc.Client.CreateOrder(&coinbasepro.Order{
		ProductID: product,
		Side:      "buy",
		Size:      fmt.Sprintf("%f", orderAmount),
	})
	if err != nil {
		fmt.Printf("Failed to buy %s\n", err.Error())
	}

	//TODO back off here?
	savedOrder, err = svc.Client.GetOrder(savedOrder.ID)
	if err != nil {
		fmt.Printf("Failed to get order %s\n", err.Error())
	}
	totalPurchased, err := strconv.ParseFloat(savedOrder.FilledSize, 64)
	if err != nil {
		fmt.Printf("Failed to parse float for filled order %s\n", err.Error())
	}
	return totalPurchased, 0.0 //available funds may be pennies
}

func (svc CoinbaseSvc) GetLastPrice(product string, ePrice float64) decimal.Decimal {
	book, err := svc.Client.GetBook(product, 1)
	if err != nil {
		fmt.Println(err.Error())
	}

	lastPrice, err := decimal.NewFromString(book.Bids[0].Price)
	if err != nil {
		fmt.Println(err.Error())
	}
	if !decimal.NewFromFloat(ePrice).Equal(lastPrice) {
		fmt.Printf("Last price not expected got %s, want %s\n", lastPrice.String(), decimal.NewFromFloat(ePrice))
	}
	return lastPrice
}
