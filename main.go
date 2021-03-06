package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JasonWBrown/svc"
	"github.com/motemen/go-loghttp"
	_ "github.com/motemen/go-loghttp/global"
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

	//set config parameters
	key := viper.GetString("api_key")
	passphrase := viper.GetString("api_passphrase")
	secret := viper.GetString("api_secret")
	product := viper.GetString("product")
	funds := viper.GetFloat64("seed")

	//create coinbase pro client
	client := coinbasepro.NewClient()
	client.UpdateConfig(&coinbasepro.ClientConfig{
		// BaseURL: "http://0.0.0.0:8080",
		BaseURL:    "https://api.pro.coinbase.com",
		Key:        key,
		Passphrase: passphrase,
		Secret:     secret,
	})

	client.HTTPClient.Transport = &loghttp.Transport{
		LogRequest: func(req *http.Request) {
			log.Printf("[%p] %s %s", req, req.Method, req.URL)
		},
		LogResponse: func(resp *http.Response) {
			log.Printf("[%p] %d %s", resp.Request, resp.StatusCode, resp.Request.URL)
		},
	}

	tSvc := svc.NewTimeSvc()
	cbSvc := svc.NewCoinbaseSvc(client, time.Duration(time.Minute*5))
	stSvc := svc.NewStateSvc()

	//in memory state tracker
	state := stSvc.NewState(product, funds)

	t := tSvc.SetInitialTime()
	for {
		state.PrintStateChange("loop begin")
		_, start, end := tSvc.GetStartAndEnd(t)
		open, close, err := cbSvc.GetMarketConditions(product, start, end)
		if err != nil {
			continue
		}

		if state.Buy(cbSvc, open, close) {
			continue
		}

		state.Lock(close)

		state.Sell(cbSvc, close)
	}
}
