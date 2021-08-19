package svc

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type CoinbaseSvcMock struct{}

func NewCoinbaseSvcMock() CoinbaseSvcMock {
	return CoinbaseSvcMock{}
}

func (svc CoinbaseSvcMock) Sell(product string, numberOwn, sellPrice float64) (float64, float64) {
	funds := numberOwn * sellPrice
	fmt.Printf("sold %f at price %f, funds available %f\n", numberOwn, sellPrice, funds)
	return 0.0, funds
}

func (svc CoinbaseSvcMock) Buy(product string, buyPrice, availablefunds float64) (float64, float64) {
	totalPurchased := availablefunds / buyPrice
	fmt.Printf("purchased %f with funds %f at price %f\n", totalPurchased, availablefunds, buyPrice)
	return totalPurchased, 0.0
}

func (svc CoinbaseSvcMock) GetLastPrice(product string, ePrice float64) decimal.Decimal {
	return decimal.NewFromFloat(ePrice)
}
