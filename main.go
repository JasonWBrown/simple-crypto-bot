package main

import (
	"fmt"
	"math"
	"time"

	"github.com/JasonWBrown/svc"
	"github.com/preichenberger/go-coinbasepro/v2"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

//TODO I don't like the global but, this is supposed to be quick and simple
//I can go back to this and make it better on version 2
var isTest bool

func main() {
	availableFunds := 1000.00
	numberOwn := 0.0
	buyPrice := 0.00
	lockPrice, bottomPrice := resetLockAndBottomPrice()
	lockPriceSet := false
	//Read in Configuration
	viper.AddConfigPath(".conf")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("failed to read config", err)
		panic(err) // this is a simple tool, this is fine
	}

	//set config parameters
	key := viper.GetString("api_key")
	passphrase := viper.GetString("api_passphrase")
	secret := viper.GetString("api_secret")
	isTest = viper.GetBool("is_test")
	product := viper.GetString("product")

	//create coinbase pro client
	client := coinbasepro.NewClient()
	client.UpdateConfig(&coinbasepro.ClientConfig{
		BaseURL: "https://api-public.sandbox.pro.coinbase.com",
		// BaseURL:    "https://api.pro.coinbase.com",
		Key:        key,
		Passphrase: passphrase,
		Secret:     secret,
	})

	var tSvc svc.TimeSvcInterface
	var cbSvc svc.CoinbaseSvcInterface
	if isTest {
		tSvc = svc.NewTimeSvcMock()
		cbSvc = svc.NewCoinbaseSvcMock()
	} else {
		tSvc = svc.NewTimeSvc()
		cbSvc = svc.NewCoinbaseSvc(client)
	}

	t := tSvc.SetInitialTime()
	rates := []coinbasepro.HistoricRate{}
	printTime(t, lockPrice, rates, decimal.Zero)
	for {
		t, start, end := tSvc.GetStartAndEnd(t)
		rates, err = client.GetHistoricRates(product, coinbasepro.GetHistoricRatesParams{
			Start:       start,
			End:         end,
			Granularity: 0,
		})
		if err != nil {
			fmt.Printf("failed to get historic rate %s\n", err.Error())
			panic(err)
		}

		lastPrice := cbSvc.GetLastPrice(product, rates[0].Close)
		printTime(t, lockPrice, rates, lastPrice)

		// buy Conditions
		if isGreaterThanPercentGrowth(rates, 0.03) && availableFunds != 0 {
			fmt.Printf("buy time %s\n", t.String())
			buyPrice = rates[0].Close
			numberOwn, availableFunds = cbSvc.Buy(product, buyPrice, availableFunds)
			bottomPrice = buyPrice - (buyPrice * .10)
			lockPriceSet = false
			continue // jump out sell and buy should not happen in the same loop
		}

		if lockPriceSet && isGrowthGreater(lockPrice, rates[0].Close, 0.01) {
			lockPrice = getLockPrice(lockPrice, rates[0].Close)
		}

		if !lockPriceSet && availableFunds == 0.0 && isGrowthGreater(buyPrice, rates[0].Close, 0.03) {
			lockPrice = rates[0].Close
			lockPriceSet = true
		}

		// sell Conditions
		if availableFunds == 0.0 && isGrowthGreater(buyPrice, rates[0].Close, 0.08) {
			fmt.Printf(".08 percent time gain %s\n", t.String())
			numberOwn, availableFunds = cbSvc.Sell(product, numberOwn, rates[0].Close)
			lockPrice, bottomPrice = resetLockAndBottomPrice()
		} else if availableFunds == 0.0 && lockPrice != 0.0 && rates[0].Close < lockPrice { //This could be set by the coinbase API
			fmt.Printf(".03 percent or greater gain time gain %s\n", t.String())
			numberOwn, availableFunds = cbSvc.Sell(product, numberOwn, lockPrice)
			lockPrice, bottomPrice = resetLockAndBottomPrice()
		} else if availableFunds == 0.0 && rates[0].Close < bottomPrice { //This could be set by the coinbase API.
			fmt.Printf("**big loss** sell time %s\n", t.String())
			numberOwn, availableFunds = cbSvc.Sell(product, numberOwn, bottomPrice)
			lockPrice, bottomPrice = resetLockAndBottomPrice()
		}
	}
}

func resetLockAndBottomPrice() (float64, float64) {
	return 0.0, 0.0
}

func getLockPrice(currentLockPrice, currentClose float64) float64 {
	return math.Max(currentLockPrice, currentClose)
}

func printTime(t time.Time, lockPrice float64, rates []coinbasepro.HistoricRate, lastPrice decimal.Decimal) {
	if isTest && t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 && t.Nanosecond() == 0 {
		fmt.Printf("t: %s\tlockPrice %f\n", t.String(), lockPrice)
	} else if len(rates) == 0 {
		fmt.Printf("t: %s\tlockPrice %f\n", t.String(), lockPrice)
	} else {
		fmt.Printf("t: %s\tlockPrice %f open %f close %f ticker.Price %s\n", t.String(), lockPrice, rates[len(rates)-1].Open, rates[0].Close, lastPrice.String())
	}
}

func isGrowthGreater(begin, end, p float64) bool {
	return percentGrowth(begin, end) > p
}

func percentGrowth(begin, end float64) float64 {
	return (end - begin) / begin
}

func isGreaterThanPercentGrowth(hr []coinbasepro.HistoricRate, p float64) bool {
	return isGrowthGreater(hr[len(hr)-1].Open, hr[0].Close, p)
}
