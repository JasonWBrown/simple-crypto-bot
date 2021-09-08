package svc

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_percentGrowth(t *testing.T) {
	type args struct {
		begin float64
		end   float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "valid 3%",
			args: args{
				begin: 100,
				end:   103,
			},
			want: 0.03,
		},
		{
			name: "divide by 0",
			args: args{
				begin: 0,
				end:   100,
			},
			want: 0,
		},
		{
			name: "negative growth",
			args: args{
				begin: 100,
				end:   97,
			},
			want: -0.03,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := percentGrowth(tt.args.begin, tt.args.end); got != tt.want {
				t.Errorf("percentGrowth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isGrowthGreater(t *testing.T) {
	type args struct {
		begin float64
		end   float64
		p     float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Greater",
			args: args{
				begin: 100,
				end:   103,
				p:     0.029,
			},
			want: true,
		},
		{
			name: "Less",
			args: args{
				begin: 100,
				end:   103,
				p:     0.031,
			},
			want: false,
		},
		{
			name: "Equal",
			args: args{
				begin: 100,
				end:   103,
				p:     0.03,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isGrowthGreater(tt.args.begin, tt.args.end, tt.args.p); got != tt.want {
				t.Errorf("isGrowthGreater() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLockPrice(t *testing.T) {
	type args struct {
		currentLockPrice float64
		currentClose     float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Current lock price higher, expect current lock price",
			args: args{
				currentLockPrice: 1.01,
				currentClose:     1.000001,
			},
			want: 1.01,
		},
		{
			name: "Current close higher, expect current close",
			args: args{
				currentLockPrice: 1.01,
				currentClose:     1.011,
			},
			want: 1.011,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLockPrice(tt.args.currentLockPrice, tt.args.currentClose); got != tt.want {
				t.Errorf("getLockPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestState_Lock(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		Product           string
		NumberOwn         float64
		BuyPrice          float64
		LockPrice         float64
		BottomPrice       float64
		LockPriceSet      bool
		AvailableUSDFunds float64
	}
	type wantFields struct {
		LockPrice    float64
		LockPriceSet bool
	}
	type args struct {
		close float64
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantFields wantFields
	}{
		{
			name: "When we haven't purchased anything. Don't change lockPrice or lockPriceSet flag",
			fields: fields{
				LockPriceSet:      false,
				AvailableUSDFunds: 0.01,
				LockPrice:         0.0,
			},
			args: args{
				close: 0.00,
			},
			wantFields: wantFields{
				LockPrice:    0.0,
				LockPriceSet: false,
			},
		},
		{
			name: "When we purchased something for the first time and growth is higher than 3% rate set the lockPrice flag and set lock price to close",
			fields: fields{
				LockPriceSet:      false,
				AvailableUSDFunds: 0.00,
				LockPrice:         0.0,
				BuyPrice:          1.0,
			},
			args: args{
				close: 1.031,
			},
			wantFields: wantFields{
				LockPrice:    1.031,
				LockPriceSet: true,
			},
		},
		{
			name: "When we purchased something for the first time and growth is lower than 3% rate, state is unchanged",
			fields: fields{
				LockPriceSet:      false,
				AvailableUSDFunds: 0.00,
				LockPrice:         0.0,
				BuyPrice:          1.0,
			},
			args: args{
				close: 1.029,
			},
			wantFields: wantFields{
				LockPrice:    0.0,
				LockPriceSet: false,
			},
		},
		{
			name: "When we have already lockedState and growth is an additional one percent, change lockPrice",
			fields: fields{
				LockPriceSet:      true,
				AvailableUSDFunds: 0.00,
				LockPrice:         1.0,
				BuyPrice:          1.0,
			},
			args: args{
				close: 1.011,
			},
			wantFields: wantFields{
				LockPrice:    1.011,
				LockPriceSet: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &State{
				Product:           tt.fields.Product,
				NumberOwn:         tt.fields.NumberOwn,
				BuyPrice:          tt.fields.BuyPrice,
				LockPrice:         tt.fields.LockPrice,
				BottomPrice:       tt.fields.BottomPrice,
				LockPriceSet:      tt.fields.LockPriceSet,
				AvailableUSDFunds: tt.fields.AvailableUSDFunds,
			}
			s.Lock(tt.args.close)
			assert.Equal(tt.wantFields.LockPrice, s.LockPrice, fmt.Sprintf("%s, Lock price is not equal", tt.name))
			assert.Equal(tt.wantFields.LockPriceSet, s.LockPriceSet, fmt.Sprintf("%s, Lock price flag is not equal", tt.name))
		})
	}
}

func TestState_Sell(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		Product           string
		NumberOwn         float64
		BuyPrice          float64
		LockPrice         float64
		BottomPrice       float64
		LockPriceSet      bool
		AvailableUSDFunds float64
	}
	type args struct {
		cbSvc CoinbaseSvcInterface
		close float64
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantFields fields
	}{
		{
			name: "8% sell will reset state",
			fields: fields{
				AvailableUSDFunds: 0.0,
				BuyPrice:          100.0,
				LockPrice:         1000.0,
				BottomPrice:       10.0,
				LockPriceSet:      true,
			},
			args: args{
				cbSvc: NewCoinbaseSvcMock(),
				close: 109.00,
			},
			wantFields: fields{
				BuyPrice:     0.0,
				LockPrice:    0.0,
				BottomPrice:  0.0,
				LockPriceSet: false,
			},
		},
		// {
		//This scenario should be handled by the api, if we close is less than bottom
		// additional revenue will be lost
		//name: "Close is less than bottom, will reset state",

		// },
		// {
		//This scenario should be handled by the api, if we close is less than lockPrice
		// additional revenue will be lost
		//name: "3% sell (close price less than lock price), will reset state",
		//
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &State{
				Product:           tt.fields.Product,
				NumberOwn:         tt.fields.NumberOwn,
				BuyPrice:          tt.fields.BuyPrice,
				LockPrice:         tt.fields.LockPrice,
				BottomPrice:       tt.fields.BottomPrice,
				LockPriceSet:      tt.fields.LockPriceSet,
				AvailableUSDFunds: tt.fields.AvailableUSDFunds,
			}
			s.Sell(tt.args.cbSvc, tt.args.close)
			assert.Equal(tt.wantFields.LockPrice, s.LockPrice, fmt.Sprintf("%s, Lock price is not equal", tt.name))
			assert.Equal(tt.wantFields.LockPriceSet, s.LockPriceSet, fmt.Sprintf("%s, Lock price flag is not equal", tt.name))
			assert.Equal(tt.wantFields.BottomPrice, s.LockPrice, fmt.Sprintf("%s, Bottom price is not equal", tt.name))
			assert.Equal(tt.wantFields.BuyPrice, s.BuyPrice, fmt.Sprintf("%s, Lock price flag is not equal", tt.name))
		})
	}
}

func TestState_Buy(t *testing.T) {
	invalidSaleTime := time.Now().Add(time.Minute * -119)
	assert := assert.New(t)
	type MockFields struct {
		Err            error
		TotalPurchased float64
		BuyPrice       float64
	}
	type fields struct {
		Product           string
		NumberOwn         float64
		BuyPrice          float64
		LockPrice         float64
		BottomPrice       float64
		LockPriceSet      bool
		AvailableUSDFunds float64
		LastSaleTime      time.Time
		MockFields        MockFields
		Executed          bool
	}
	type args struct {
		open  float64
		close float64
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantFields fields
	}{
		{
			name: "Happy Path. 3% growth, with available funds, and sold greater than 120 minutes will buy.",
			fields: fields{
				AvailableUSDFunds: 1.0, // not zero
				BuyPrice:          100.1,
				LockPrice:         1000.0,
				BottomPrice:       10.0,
				LockPriceSet:      true,
				LastSaleTime:      time.Now().Add(time.Minute * -121),
				MockFields: MockFields{
					Err:            nil,
					TotalPurchased: 9.9999,
					BuyPrice:       14.123,
				},
			},
			args: args{
				open:  100.00,
				close: 103.0001,
			},
			wantFields: fields{
				BuyPrice:     14.123,
				LockPrice:    0.0,
				BottomPrice:  12.7107,
				LockPriceSet: false,
				Executed:     true,
				NumberOwn:    9.9999,
				LastSaleTime: time.Time{},
			},
		},
		{
			name: "Sad Path. 3% growth, with available funds, and sold less than 120 minutes will not buy.",
			fields: fields{
				AvailableUSDFunds: 1.0, // not zero
				BuyPrice:          100.1,
				LockPrice:         1000.0,
				BottomPrice:       10.0,
				LockPriceSet:      true,
				LastSaleTime:      invalidSaleTime,
				MockFields: MockFields{
					Err:            nil,
					TotalPurchased: 9.9999,
					BuyPrice:       14.123,
				},
			},
			args: args{
				open:  100.00,
				close: 103.0001,
			},
			wantFields: fields{
				BuyPrice:     100.1,
				LockPrice:    1000.0,
				BottomPrice:  10.0,
				LockPriceSet: true,
				Executed:     false,
				NumberOwn:    0.0,
				LastSaleTime: invalidSaleTime,
			},
		},
		{
			name: "Sad Path. 3% growth, with available funds, and sold time is empty will not buy.",
			fields: fields{
				AvailableUSDFunds: 1.0, // not zero
				BuyPrice:          100.1,
				LockPrice:         1000.0,
				BottomPrice:       10.0,
				LockPriceSet:      true,
				LastSaleTime:      time.Time{},
				MockFields: MockFields{
					Err:            nil,
					TotalPurchased: 9.9999,
					BuyPrice:       14.123,
				},
			},
			args: args{
				open:  100.00,
				close: 103.0001,
			},
			wantFields: fields{
				BuyPrice:     100.1,
				LockPrice:    1000.0,
				BottomPrice:  10.0,
				LockPriceSet: true,
				Executed:     false,
				NumberOwn:    0.0,
				LastSaleTime: time.Time{},
			},
		},
	}
	for _, tt := range tests {
		cbSvcMock := NewCoinbaseSvcMock()
		cbSvcMock.BuyPrice = tt.fields.MockFields.BuyPrice
		cbSvcMock.Err = tt.fields.MockFields.Err
		cbSvcMock.TotalPurchased = tt.fields.MockFields.TotalPurchased

		t.Run(tt.name, func(t *testing.T) {

			s := &State{
				Product:           tt.fields.Product,
				NumberOwn:         tt.fields.NumberOwn,
				BuyPrice:          tt.fields.BuyPrice,
				LockPrice:         tt.fields.LockPrice,
				BottomPrice:       tt.fields.BottomPrice,
				LockPriceSet:      tt.fields.LockPriceSet,
				AvailableUSDFunds: tt.fields.AvailableUSDFunds,
				LastSaleTime:      tt.fields.LastSaleTime,
			}
			executed := s.Buy(cbSvcMock, tt.args.open, tt.args.close)
			assert.Equal(tt.wantFields.Executed, executed, fmt.Sprintf("%s, Lock price is not equal", tt.name))
			assert.Equal(tt.wantFields.LockPrice, s.LockPrice, fmt.Sprintf("%s, Lock price is not equal", tt.name))
			assert.Equal(tt.wantFields.LockPriceSet, s.LockPriceSet, fmt.Sprintf("%s, Lock price flag is not equal", tt.name))
			assert.Equal(tt.wantFields.BottomPrice, s.BottomPrice, fmt.Sprintf("%s, Bottom price is not equal", tt.name))
			assert.Equal(tt.wantFields.BuyPrice, s.BuyPrice, fmt.Sprintf("%s, Buy price is not equal", tt.name))
			assert.Equal(tt.wantFields.NumberOwn, s.NumberOwn, fmt.Sprintf("%s, NumberOwn is not equal", tt.name))
			assert.Equal(tt.wantFields.LastSaleTime, s.LastSaleTime, fmt.Sprintf("%s, LastSaleTime is not equal", tt.name))
		})
	}
}
