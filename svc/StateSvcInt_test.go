package svc

import (
	"fmt"
	"testing"
	"time"

	"github.com/preichenberger/go-coinbasepro/v2"
	"github.com/stretchr/testify/assert"
	"github.com/walkerus/go-wiremock"
)

func TestBuy(t *testing.T) {
	// Set up the cbSvc this will be an integration test so we will create the actual client
	// the client will call an imposter through wiremock
	client := coinbasepro.NewClient()
	client.UpdateConfig(&coinbasepro.ClientConfig{BaseURL: "http://0.0.0.0:8080"})
	cbSvc := NewCoinbaseSvc(client, time.Duration(time.Millisecond*1))
	assert := assert.New(t)

	//set up wire mock
	//given there exists apis in wiremock
	wiremockClient := wiremock.NewClient("http://0.0.0.0:8080")
	wiremockClient.Clear()

	type fields struct {
		OrderID           string
		Product           string
		NumberOwn         float64
		BuyPrice          float64
		LockPrice         float64
		BottomPrice       float64
		LockPriceSet      bool
		AvailableUSDFunds float64
		Executed          bool
	}
	type args struct {
		svc   CoinbaseSvcInterface
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
			name: "Execute buy when growth is greater than expected and have USD",
			fields: fields{
				OrderID:           "d0c5340b-6d6c-49d9-b567-48c4bfca13d2",
				Product:           "BTC-USD",
				NumberOwn:         0,
				BuyPrice:          0,
				LockPrice:         0,
				BottomPrice:       0,
				LockPriceSet:      false,
				AvailableUSDFunds: 1000.0,
			},
			args: args{
				svc:   cbSvc,
				open:  100.0,
				close: 103.1,
			},
			wantFields: fields{
				Product:           "BTC-USD",
				NumberOwn:         0,
				BuyPrice:          103.1,
				LockPrice:         0,
				BottomPrice:       0,
				LockPriceSet:      false,
				AvailableUSDFunds: 0.0,
				Executed:          true,
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

			// Set wire mock to return an expected order id
			wiremockClient.StubFor(wiremock.Post(wiremock.URLPathMatching("/orders")).
				WillReturn(
					`{`+
						`"id": "`+tt.fields.OrderID+`"`+
						`}`,
					map[string]string{"Content-Type": "application/json"},
					200,
				).
				AtPriority(1))

			//Set wire mock to return a fullfilled status at some point
			wiremockClient.StubFor(wiremock.Get(wiremock.URLPathMatching(fmt.Sprintf("/orders/%s", tt.fields.OrderID))).
				WillReturn(
					`{`+
						`"id": "`+tt.fields.OrderID+`"`+
						`"product_id": "BTC-USD",`+
						`"type": "market",`+
						`"post_only": false,`+
						`"created_at": "2016-12-08T20:09:05.508883Z",`+
						`"done_at": "2016-12-08T20:09:05.527Z",`+
						`"done_reason": "filled",`+
						`"fill_fees": "0.0249376391550000",`+
						`"filled_size": "0.01291771",`+
						`"executed_value": "9.9750556620000000",`+
						`"status": "done",`+
						`"settled": true`+
						`}`,
					map[string]string{"Content-Type": "application/json"},
					200,
				).
				AtPriority(2))

			executed := s.Buy(tt.args.svc, tt.args.open, tt.args.close)
			assert.Equal(tt.wantFields.AvailableUSDFunds, s.AvailableUSDFunds, fmt.Sprintf("%s, AvailableUSDFunds is not equal", tt.name))
			assert.Equal(tt.wantFields.NumberOwn, s.NumberOwn, fmt.Sprintf("%s, NumberOwn is not equal", tt.name))
			assert.Equal(tt.wantFields.BuyPrice, s.BuyPrice, fmt.Sprintf("%s, BuyPrice is not equal", tt.name))
			assert.Equal(tt.wantFields.LockPrice, s.LockPrice, fmt.Sprintf("%s, LockPrice is not equal", tt.name))
			assert.Equal(tt.wantFields.LockPriceSet, s.LockPriceSet, fmt.Sprintf("%s, LockPriceSet flag is not equal", tt.name))
			assert.Equal(tt.wantFields.Executed, executed)
		})
	}
}
