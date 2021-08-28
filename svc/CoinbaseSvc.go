package svc

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/preichenberger/go-coinbasepro/v2"
)

type CoinbaseSvcInterface interface {
	Sell(product string, numberOwn, sellPrice float64) (float64, float64)
	Buy(product string, buyPrice, availablefunds float64) (float64, float64)
	GetLastPrice(product string) float64
	GetMarketConditions(product string, start, end time.Time) (float64, float64)
}

type CoinbaseSvc struct {
	Client  *coinbasepro.Client
	Timeout time.Duration
}

func NewCoinbaseSvc(client *coinbasepro.Client, d time.Duration) CoinbaseSvc {
	return CoinbaseSvc{
		Client:  client,
		Timeout: d,
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
	savedOrder, err := svc.Client.CreateOrder(&coinbasepro.Order{
		ProductID: product,
		Side:      "buy",
		Funds:     fmt.Sprintf("%f", availablefunds),
	})
	if err != nil {
		fmt.Printf("Failed to buy %s\n", err.Error())
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = time.Duration(svc.Timeout)
	totalPurchased := 0.0
	err = backoff.Retry(func() error {
		savedOrder, err = svc.Client.GetOrder(savedOrder.ID)
		if err != nil {
			fmt.Printf("Failed to get order %s\n", err.Error())
			return err
		}

		if savedOrder.Status != "done" && savedOrder.DoneReason != "filled" {
			errMessage := fmt.Sprintf("failed to get expected order Status got %s, want %s and DoneReason got %s, want %s\n", savedOrder.Status, "done", savedOrder.DoneReason, "filled")
			fmt.Println(errMessage)
			return fmt.Errorf(errMessage)
		}
		totalPurchased, err = strconv.ParseFloat(savedOrder.FilledSize, 64)
		if err != nil {
			fmt.Printf("Failed to parse float for filled order %s\n", err.Error())
			return err
		}
		return nil
	}, b)

	return totalPurchased, 0.0 //available funds may be pennies
}

func (svc CoinbaseSvc) GetLastPrice(product string) float64 {
	book, err := svc.Client.GetBook(product, 1)
	if err != nil {
		fmt.Println(err.Error())
	}

	lastPrice, err := strconv.ParseFloat(book.Bids[0].Price, 64)
	if err != nil {
		fmt.Println(err.Error())
	}
	return lastPrice
}

func (svc CoinbaseSvc) GetMarketConditions(product string, start, end time.Time) (float64, float64) {
	rates, err := svc.Client.GetHistoricRates(product, coinbasepro.GetHistoricRatesParams{
		Start:       start,
		End:         end,
		Granularity: 0,
	})
	if err != nil {
		fmt.Printf("failed to get historic rate %s\n", err.Error())
		panic(err)
	}

	lastPrice := svc.GetLastPrice(product)

	return rates[len(rates)-1].Open, lastPrice
}
