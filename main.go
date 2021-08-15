package main

import (
	"fmt"
	"math"
	"time"

	"github.com/preichenberger/go-coinbasepro/v2"
	"github.com/spf13/viper"
)

func main() {

	availableFunds := 1000.00
	numberOwn := 0.0
	buyPrice := 0.00
	lockPrice := 0.00
	//Read in Configuration
	viper.AddConfigPath(".conf")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("failed to read config", err)
		panic(err) // this is a simple tool, this is fine
	}

	key := viper.GetString("api_key")
	passphrase := viper.GetString("api_passphrase")
	secret := viper.GetString("api_secret")

	client := coinbasepro.NewClient()

	// optional, configuration can be updated with ClientConfig
	client.UpdateConfig(&coinbasepro.ClientConfig{
		BaseURL:    "https://api.pro.coinbase.com",
		Key:        key,
		Passphrase: passphrase,
		Secret:     secret,
	})

	t := time.Date(2020, time.December, 17, 0, 0, 0, 0, time.UTC)
	for {
		// ticker, err := client.GetTicker("BTC-USD")
		// if err != nil {
		// 	fmt.Println(err)
		// 	panic(err)
		// }
		printTime(t, lockPrice)
		t = t.Add(time.Minute * 20)
		if t.Year() == 2021 && t.Month() == time.August {
			break
		}
		rates, err := client.GetHistoricRates("BTC-USD", coinbasepro.GetHistoricRatesParams{
			Start:       t,
			End:         t.Add(time.Hour * 2),
			Granularity: 0,
		})

		if err != nil {
			fmt.Printf("failed to get historic rate %s\n", err.Error())
			panic(err)
		}

		// buy Conditions
		if isGreaterThanPercentGrowth(rates, 0.03) && availableFunds != 0 {
			fmt.Printf("buy time %s\n", t.String())
			buyPrice = rates[0].Close
			numberOwn, availableFunds = buy(buyPrice, availableFunds)
			lockPrice = buyPrice - (buyPrice * .05)
			continue // jump out sell and buy should not happen in the same loop
		}

		lockPrice = getLockPrice(lockPrice, rates[0].Close)

		// sell Conditions
		if availableFunds == 0.0 && isGrowthGreater(buyPrice, rates[0].Close, 0.08) {
			fmt.Printf(".08 percent gain boom sell time gain %s\n", t.String())
			numberOwn, availableFunds = sell(numberOwn, rates[0].Close)
		} else if availableFunds == 0.0 && rates[0].Close < lockPrice {
			fmt.Printf("sell time gain %s\n", t.String())
			numberOwn, availableFunds = sell(numberOwn, lockPrice)
		}
		// } else if availableFunds == 0.0 && isGrowthLess(buyPrice, rates[0].Close, -0.02) {
		// 	fmt.Printf("sell time loss %s\n", t.String())
		// 	numberOwn, availableFunds = sell(numberOwn, rates[0].Close)
		// }
	}
}

func getLockPrice(currentLockPrice, currentClose float64) float64 {
	possibleLockPrice := currentClose - (currentClose * .05)
	return math.Max(currentLockPrice, possibleLockPrice)
}

func printTime(t time.Time, f float64) {
	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 && t.Nanosecond() == 0 {
		fmt.Printf("t: %s\tlockPrice %f\n", t.String(), f)
	}
}

func sell(numberOwn, sellPrice float64) (float64, float64) {
	funds := numberOwn * sellPrice
	fmt.Printf("sold %f at price %f, funds available %f\n", numberOwn, sellPrice, funds)
	return 0.0, funds
}

func buy(buyPrice, availablefunds float64) (float64, float64) {
	totalPurchased := availablefunds / buyPrice
	fmt.Printf("purchased %f with funds %f at price %f\n", totalPurchased, availablefunds, buyPrice)
	return totalPurchased, 0.0
}

func isGrowthLess(begin, end, p float64) bool {
	return percentGrowth(begin, end) < p
}

func isGrowthGreater(begin, end, p float64) bool {
	return percentGrowth(begin, end) > p
}

func percentGrowth(begin, end float64) float64 {
	return (end - begin) / begin
}

func isGreaterThanPercentGrowth(hr []coinbasepro.HistoricRate, p float64) bool {
	// fmt.Printf("Historic Rates %+v\n", hr)
	// fmt.Printf("open %s %.2f\n", hr[len(hr)-1].Time.String(), hr[len(hr)-1].Open)
	// fmt.Printf("close %s %.2f\n", hr[0].Time.String(), hr[0].Close)
	growth := (hr[0].Close - hr[len(hr)-1].Open) / hr[0].Close
	// fmt.Printf("Growth is %.3f\t open %f close %f\n", growth, hr[len(hr)-1].Open, hr[0].Close)
	return growth > p
}
