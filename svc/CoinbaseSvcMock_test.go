package svc

import (
	"fmt"
	"time"
)

type CoinbaseSvcMock struct {
	Err               error
	TotalPurchased    float64
	AvailableUSDFunds float64
	BuyPrice          float64
}

func NewCoinbaseSvcMock() CoinbaseSvcMock {
	return CoinbaseSvcMock{}
}

func (svc CoinbaseSvcMock) Sell(product string, numberOwn, sellPrice float64) (float64, float64, error) {
	funds := numberOwn * sellPrice
	fmt.Printf("sold %f at price %f, funds available %f\n", numberOwn, sellPrice, funds)
	return 0.0, funds, nil
}

func (svc CoinbaseSvcMock) Buy(product string, buyPrice, availablefunds float64) (float64, float64, error) {
	return svc.TotalPurchased, svc.BuyPrice, svc.Err
}

func (svc CoinbaseSvcMock) GetLastPrice(product string) (float64, error) {
	return 1.01, nil
}

func (svc CoinbaseSvcMock) GetMarketConditions(product string, start, end time.Time) (float64, float64, error) {
	return 2.02, 4.04, nil
}
