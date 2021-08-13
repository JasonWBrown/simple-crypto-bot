package main

import (
	"fmt"

	"github.com/preichenberger/go-coinbasepro/v2"
	"github.com/spf13/viper"
)

func main() {

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

	stats, err := client.GetStats("BTC-USD")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Printf("Stats: %+v\n", stats)

	ticker, err := client.GetTicker("BTC-USD")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	lastPrice := ticker.Price
	fmt.Printf("Last Price: %+v\n", lastPrice)
}
