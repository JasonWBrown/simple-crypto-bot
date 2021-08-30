package svc

import (
	"fmt"
	"math"
	"time"
)

type StateSvc struct {
	//postgres client
}

func NewStateSvc() *StateSvc {
	return &StateSvc{}
}

type State struct {
	Product           string
	NumberOwn         float64
	BuyPrice          float64
	LockPrice         float64
	BottomPrice       float64
	LockPriceSet      bool
	AvailableUSDFunds float64
}

func (s *State) ResetState() {
	s.BuyPrice = 0.0
	s.LockPrice = 0.0
	s.BottomPrice = 0.0
	s.LockPriceSet = false
}

func (svc StateSvc) NewState(product string, funds float64) *State {
	return &State{
		Product:           product,
		NumberOwn:         0.0,
		BuyPrice:          0.0,
		LockPrice:         0.0,
		BottomPrice:       0.0,
		LockPriceSet:      false,
		AvailableUSDFunds: funds,
	}
}

func (s *State) PrintStateChange(trigger string) {
	fmt.Printf("%s %s state change, %+v", time.Now().Format(time.RFC822), trigger, s)
}

func (s *State) Buy(cbSvc CoinbaseSvcInterface, open, close float64) bool {
	if isGrowthGreater(open, close, .03) && s.AvailableUSDFunds != 0 {
		nOwn, funds, err := cbSvc.Buy(s.Product, close, s.AvailableUSDFunds)
		if err != nil {
			return false
		}
		s.NumberOwn = nOwn
		s.AvailableUSDFunds = funds
		s.BuyPrice = close
		s.BottomPrice = s.BuyPrice - (s.BuyPrice * .10)
		s.LockPriceSet = false
		s.PrintStateChange("buy")
		return true
	}
	return false
}

func (s *State) Lock(close float64) {
	if s.LockPriceSet && isGrowthGreater(s.LockPrice, close, 0.01) {
		s.LockPrice = getLockPrice(s.LockPrice, close)
		s.PrintStateChange("Lock growth of 1%")
	}

	if !s.LockPriceSet && s.AvailableUSDFunds == 0.0 && isGrowthGreater(s.BuyPrice, close, 0.03) {
		s.LockPrice = close
		s.LockPriceSet = true
		s.PrintStateChange("Lock growth of 3%")
	}
}

func (s *State) Sell(cbSvc CoinbaseSvcInterface, close float64) {
	if s.AvailableUSDFunds == 0.0 && isGrowthGreater(s.BuyPrice, close, 0.08) {
		s.NumberOwn, s.AvailableUSDFunds = cbSvc.Sell(s.Product, s.NumberOwn, close)
		s.ResetState()
		s.PrintStateChange("Sell 8% gain")
	} else if s.AvailableUSDFunds == 0.0 && s.LockPrice != 0.0 && close < s.LockPrice { //This could be set by the coinbase API
		fmt.Printf(".03 percent or greater gain time gain %s\n", time.Now())
		s.NumberOwn, s.AvailableUSDFunds = cbSvc.Sell(s.Product, s.NumberOwn, s.LockPrice)
		s.ResetState()
		s.PrintStateChange("Sell 3% gain")
	} else if s.AvailableUSDFunds == 0.0 && close < s.BottomPrice { //This could be set by the coinbase API.
		s.NumberOwn, s.AvailableUSDFunds = cbSvc.Sell(s.Product, s.NumberOwn, s.BottomPrice)
		s.ResetState()
		s.PrintStateChange("*Sell 10% loss*")
	}
}

func getLockPrice(currentLockPrice, currentClose float64) float64 {
	return math.Max(currentLockPrice, currentClose)
}

func isGrowthGreater(begin, end, p float64) bool {
	return percentGrowth(begin, end) > p
}

func percentGrowth(begin, end float64) float64 {
	if begin == 0 {
		return 0.0
	}
	return (end - begin) / begin
}
