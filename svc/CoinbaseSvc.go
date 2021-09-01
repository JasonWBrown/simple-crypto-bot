package svc

import (
	"fmt"
	"strconv"
	"time"

	"github.com/JasonWBrown/proclient"
	"github.com/cenkalti/backoff/v4"
	"github.com/preichenberger/go-coinbasepro/v2"
)

type CoinbaseSvcInterface interface {
	Sell(product string, numberOwn, sellPrice float64) (float64, float64, error)
	Buy(product string, buyPrice, availablefunds float64) (float64, float64, error)
	GetLastPrice(product string) (float64, error)
	GetMarketConditions(product string, start, end time.Time) (float64, float64, error)
}

type CoinbaseSvc struct {
	Client  proclient.ProClientInterface
	Timeout time.Duration
}

func NewCoinbaseSvc(client proclient.ProClientInterface, d time.Duration) CoinbaseSvc {
	return CoinbaseSvc{
		Client:  client,
		Timeout: d,
	}
}

//Sell
//NumberOwn, AvailableUSDFunds, error := Sell()
func (svc CoinbaseSvc) Sell(product string, numberOwn, sellPrice float64) (float64, float64, error) {
	savedOrder, err := svc.Client.CreateOrder(&coinbasepro.Order{
		ProductID: product,
		Side:      "sell",
		Size:      fmt.Sprintf("%f", numberOwn),
	})
	if err != nil {
		fmt.Printf("Failed to sell %s\n", err.Error())
		return numberOwn, 0.0, err
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = time.Duration(svc.Timeout)
	funds := 0.0
	err = backoff.Retry(func() error {
		fmt.Printf("Entering backoff.")
		savedOrder, err = svc.Client.GetOrder(savedOrder.ID)
		if err != nil {
			fmt.Printf("Failed to get order %s\n", err.Error())
			return err
		}

		//TODO how do I get my funds available and how do I know the order completed.
		accounts, err := svc.Client.GetAccounts()
		if err != nil {
			fmt.Printf("Failed to get accounts %s\n", err.Error())
			return err
		}
		var account coinbasepro.Account
		// these might be in order might not have to iterate every single time
		for _, a := range accounts {
			if a.Currency == "USD" {
				account = a
				break
			}
		}

		funds, err = strconv.ParseFloat(account.Balance, 64)
		if err != nil {
			fmt.Printf("Failed to parse float for account balance %s\n", err.Error())
			return err
		}

		return nil
	}, b)

	if err != nil {
		fmt.Printf("Failed to sell %s\n", err.Error())
		return numberOwn, 0.0, err
	}

	return 0.0, funds, nil
}

func (svc CoinbaseSvc) Buy(product string, buyPrice, availablefunds float64) (float64, float64, error) {
	savedOrder, err := svc.Client.CreateOrder(&coinbasepro.Order{
		ProductID: product,
		Side:      "buy",
		Funds:     fmt.Sprintf("%f", availablefunds),
	})
	if err != nil {
		fmt.Printf("Failed to buy %s\n", err.Error())
		return 0.0, availablefunds, err
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = time.Duration(svc.Timeout)
	totalPurchased := 0.0
	err = backoff.Retry(func() error {
		fmt.Printf("Entering backoff.")
		so, err := svc.Client.GetOrder(savedOrder.ID)
		if err != nil {
			fmt.Printf("Failed to get order %s\n", err.Error())
			return err
		}
		fmt.Printf("Saved order %+v", so)
		if so.Status != "done" || so.DoneReason != "filled" {
			errMessage := fmt.Sprintf("failed to get expected order Status got %s, want %s and DoneReason got %s, want %s", so.Status, "done", so.DoneReason, "filled")
			fmt.Println(errMessage)
			return fmt.Errorf(errMessage)
		}
		totalPurchased, err = strconv.ParseFloat(so.FilledSize, 64)
		if err != nil {
			fmt.Printf("Failed to parse float for filled order %s\n", err.Error())
			return err
		}
		return nil
	}, b)
	if err != nil {
		return 0.0, availablefunds, err
	}
	return totalPurchased, 0.0, nil //available funds may be pennies
}

func (svc CoinbaseSvc) GetLastPrice(product string) (float64, error) {
	book, err := svc.Client.GetBook(product, 1)
	if err != nil {
		fmt.Println(err.Error())
		return -100.0, err
	}

	if len(book.Bids) == 0 {
		return -100.0, fmt.Errorf("failed to get books expecting array to be populated")
	}

	lastPrice, err := strconv.ParseFloat(book.Bids[0].Price, 64)
	if err != nil {
		fmt.Println(err.Error())
		return -100.0, err
	}
	return lastPrice, err
}

func (svc CoinbaseSvc) GetMarketConditions(product string, start, end time.Time) (float64, float64, error) {
	rates, err := svc.Client.GetHistoricRates(product, coinbasepro.GetHistoricRatesParams{
		Start:       start,
		End:         end,
		Granularity: 0,
	})
	if err != nil {
		fmt.Printf("failed to get historic rate %s\n", err.Error())
		return 0.0, 0.0, err
	}

	lastPrice, err := svc.GetLastPrice(product)
	if err != nil {
		fmt.Printf("failed to get last price %s\n", err.Error())
		return 0.0, 0.0, err
	}

	return rates[len(rates)-1].Open, lastPrice, nil
}
