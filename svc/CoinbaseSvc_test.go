package svc

import (
	"fmt"
	"testing"
	"time"

	"github.com/JasonWBrown/proclient"
	"github.com/preichenberger/go-coinbasepro/v2"
)

func TestCoinbaseSvc_GetMarketConditions(t *testing.T) {
	type fields struct {
		WantErr error
		rates   []coinbasepro.HistoricRate
		book    coinbasepro.Book
	}
	type args struct {
		product string
		start   time.Time
		end     time.Time
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantStart float64
		wantEnd   float64
		wantErr   error
	}{
		{
			name: "Happy Path. No Error from client.",
			fields: fields{
				WantErr: nil,
				rates: []coinbasepro.HistoricRate{
					{
						Time:   time.Now(),
						Low:    0.1,
						High:   0.2,
						Open:   0.3,
						Close:  0.4,
						Volume: 1000.1,
					},
					{
						Time:   time.Now(),
						Low:    1.1,
						High:   1.2,
						Open:   1.3,
						Close:  1.4,
						Volume: 2000.1,
					},
				},
				book: coinbasepro.Book{
					Bids: []coinbasepro.BookEntry{
						{
							Price:          "4.0",
							Size:           "1.0",
							NumberOfOrders: 400,
							OrderID:        "GUID-1",
						},
					},
				},
			},
			args: args{
				product: "BTC-USD",
				start:   time.Now().Add(time.Hour * -1),
				end:     time.Now(),
			},
			wantStart: 1.3, // len(historicrates)-1.open
			wantEnd:   4.0,
			wantErr:   nil,
		},
		{
			name: "Sand Path. Error From Client.",
			fields: fields{
				WantErr: fmt.Errorf("Its broke"),
				rates: []coinbasepro.HistoricRate{
					{
						Time:   time.Now(),
						Low:    0.1,
						High:   0.2,
						Open:   0.3,
						Close:  0.4,
						Volume: 1000.1,
					},
					{
						Time:   time.Now(),
						Low:    1.1,
						High:   1.2,
						Open:   1.3,
						Close:  1.4,
						Volume: 2000.1,
					},
				},
				book: coinbasepro.Book{
					Bids: []coinbasepro.BookEntry{
						{
							Price:          "4.0",
							Size:           "1.0",
							NumberOfOrders: 400,
							OrderID:        "GUID-1",
						},
					},
				},
			},
			args: args{
				product: "BTC-USD",
				start:   time.Now().Add(time.Hour * -1),
				end:     time.Now(),
			},
			wantStart: 0.0,
			wantEnd:   0.0,
			wantErr:   fmt.Errorf("Its broke"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := proclient.NewMockClient()
			c.Err = tt.fields.WantErr
			c.HistoricRates = tt.fields.rates
			c.Book = tt.fields.book

			svc := CoinbaseSvc{
				Client:  c,
				Timeout: time.Duration(time.Millisecond), // for all tests, no backoff necessary
			}
			start, end, err := svc.GetMarketConditions(tt.args.product, tt.args.start, tt.args.end)
			if start != tt.wantStart {
				t.Errorf("CoinbaseSvc.GetMarketConditions() got = %v, want %v", start, tt.wantStart)
			}
			if end != tt.wantEnd {
				t.Errorf("CoinbaseSvc.GetMarketConditions() got1 = %v, want %v", end, tt.wantEnd)
			}

			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CoinbaseSvc.GetMarketConditions() err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestCoinbaseSvc_GetLastPrice(t *testing.T) {
	type fields struct {
		WantErr error
		book    coinbasepro.Book
	}
	type args struct {
		product string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr error
	}{
		{
			name: "Happy Path. Last Price Returns expected values.",
			fields: fields{
				WantErr: nil,
				book: coinbasepro.Book{
					Bids: []coinbasepro.BookEntry{
						{
							Price:          "44.0",
							Size:           "1.0",
							NumberOfOrders: 400,
							OrderID:        "GUID-1",
						},
						{
							Price:          "9.0",
							Size:           "2.0",
							NumberOfOrders: 800,
							OrderID:        "GUID-2",
						},
					},
				},
			},
			args: args{
				"BTC-USD",
			},
			want:    44.0,
			wantErr: nil,
		},
		{
			name: "Sad Path. Client does not return error but there are no bids will return error.",
			fields: fields{
				WantErr: nil,
				book: coinbasepro.Book{
					Bids: []coinbasepro.BookEntry{},
				},
			},
			args: args{
				"BTC-USD",
			},
			want:    -100.0,
			wantErr: fmt.Errorf("failed to get books expecting array to be populated"), //tight coupling between error messages is by design
		},
		{
			name: "Sad Path. Client return error return error.",
			fields: fields{
				WantErr: fmt.Errorf("its broke"),
			},
			args: args{
				"BTC-USD",
			},
			want:    -100.0,
			wantErr: fmt.Errorf("its broke"), //tight coupling between error messages is by design
		},
		{
			name: "Sad Path. Client does not return error but bids does not have price.",
			fields: fields{
				WantErr: nil,
				book: coinbasepro.Book{
					Bids: []coinbasepro.BookEntry{
						{
							Size: "100",
						},
						{
							Price:          "9.0",
							Size:           "2.0",
							NumberOfOrders: 800,
							OrderID:        "GUID-2",
						},
					},
				},
			},
			args: args{
				"BTC-USD",
			},
			want:    -100.0,
			wantErr: fmt.Errorf(`strconv.ParseFloat: parsing "": invalid syntax`), //tight coupling between error messages is by design
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := proclient.NewMockClient()
			c.Err = tt.fields.WantErr
			c.Book = tt.fields.book
			svc := CoinbaseSvc{
				Client:  c,
				Timeout: time.Duration(time.Millisecond), // for all tests, no backoff necessary
			}
			got, err := svc.GetLastPrice(tt.args.product)
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CoinbaseSvc.GetLastPrice() err = %v, want %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("CoinbaseSvc.GetLastPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoinbaseSvc_Buy(t *testing.T) {
	type fields struct {
		wantErr error
		order   coinbasepro.Order
	}
	type args struct {
		product        string
		buyPrice       float64
		availablefunds float64
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantTotalPurchased float64
		wantBuyPrice       float64
		wantErr            error
	}{
		{
			name: "Happy Path. No error from client. Order is fullfilled and done",
			fields: fields{
				wantErr: nil,
				order: coinbasepro.Order{
					ID:            "GUID-99",
					FilledSize:    "2.0",
					Status:        "done",
					DoneReason:    "filled",
					ExecutedValue: "1000.00",
				},
			},
			args: args{
				product:        "SOME-PRODUCT",
				buyPrice:       99.998,
				availablefunds: 99.997,
			},
			wantTotalPurchased: 2,
			wantBuyPrice:       500.00,
			wantErr:            nil,
		},
		{
			name: "Sad Path.  Error from client.",
			fields: fields{
				wantErr: fmt.Errorf("this is broke"),
				order: coinbasepro.Order{
					ID:            "GUID-99",
					FilledSize:    "99.9999",
					Status:        "done",
					DoneReason:    "filled",
					ExecutedValue: "1000.00",
				},
			},
			args: args{
				product:        "SOME-PRODUCT",
				buyPrice:       99.996,
				availablefunds: 99.997,
			},
			wantTotalPurchased: 0.0,
			wantBuyPrice:       0.0,
			wantErr:            fmt.Errorf("this is broke"),
		},
		{
			name: "Sad Path.  No error from client. Status not done",
			fields: fields{
				wantErr: nil,
				order: coinbasepro.Order{
					ID:            "GUID-99",
					FilledSize:    "99.9999",
					Status:        "not_done",
					DoneReason:    "filled",
					ExecutedValue: "1000.00",
				},
			},
			args: args{
				product:        "SOME-PRODUCT",
				buyPrice:       99.996,
				availablefunds: 99.99788,
			},
			wantTotalPurchased: 0.0,
			wantBuyPrice:       0.0,
			wantErr:            fmt.Errorf("failed to get expected order Status got not_done, want done and DoneReason got filled, want filled"),
		},
		{
			name: "Sad Path.  No error from client. Done reason not filled",
			fields: fields{
				wantErr: nil,
				order: coinbasepro.Order{
					ID:            "GUID-99",
					FilledSize:    "99.9999",
					Status:        "done",
					DoneReason:    "not_filled",
					ExecutedValue: "1000.00",
				},
			},
			args: args{
				product:        "SOME-PRODUCT",
				buyPrice:       99.996,
				availablefunds: 99.99788,
			},
			wantTotalPurchased: 0.0,
			wantBuyPrice:       0.0,
			wantErr:            fmt.Errorf("failed to get expected order Status got done, want done and DoneReason got not_filled, want filled"),
		},
		{
			name: "Sad Path.  No fill size returns error",
			fields: fields{
				wantErr: nil, //no mocked error still returns error
				order: coinbasepro.Order{
					ID:            "GUID-99",
					Status:        "done",
					DoneReason:    "filled",
					ExecutedValue: "1000.00",
				},
			},
			args: args{
				product:        "SOME-PRODUCT",
				buyPrice:       99.996,
				availablefunds: 99.997,
			},
			wantTotalPurchased: 0.0,
			wantBuyPrice:       0.0,
			wantErr:            fmt.Errorf(`strconv.ParseFloat: parsing "": invalid syntax`),
		},
		{
			name: "Sad Path.  No ExecutedValue returns error",
			fields: fields{
				wantErr: nil, //no mocked error still returns error
				order: coinbasepro.Order{
					ID:         "GUID-99",
					Status:     "done",
					DoneReason: "filled",
					FilledSize: "1.0",
				},
			},
			args: args{
				product:        "SOME-PRODUCT",
				buyPrice:       99.996,
				availablefunds: 99.997,
			},
			wantTotalPurchased: 0.0,
			wantBuyPrice:       0.0,
			wantErr:            fmt.Errorf(`strconv.ParseFloat: parsing "": invalid syntax`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := proclient.NewMockClient()
			c.Err = tt.fields.wantErr
			c.SavedOrder = tt.fields.order
			svc := CoinbaseSvc{
				Client:  c,
				Timeout: time.Duration(time.Millisecond), // for all tests, no backoff necessary
			}
			totalPurchased, buyPrice, err := svc.Buy(tt.args.product, tt.args.buyPrice, tt.args.buyPrice)
			if totalPurchased != tt.wantTotalPurchased {
				t.Errorf("CoinbaseSvc.Buy() totalPurchased = %v, want %v", totalPurchased, tt.wantTotalPurchased)
			}
			if buyPrice != tt.wantBuyPrice {
				t.Errorf("CoinbaseSvc.Buy() buyPrice = %v, want %v", buyPrice, tt.wantBuyPrice)
			}
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CoinbaseSvc.Buy() err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestCoinbaseSvc_Sell(t *testing.T) {
	type fields struct {
		wantErr  error
		order    coinbasepro.Order
		accounts []coinbasepro.Account
	}
	type args struct {
		product   string
		numberOwn float64
		sellPrice float64
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantNumberOwn      float64
		wantAvailableFunds float64
		wantErr            error
	}{
		{
			name: "Happy Path. No error from client.",
			fields: fields{
				wantErr: nil,
				order: coinbasepro.Order{
					ID:         "GUID-99",
					FilledSize: "99.9999",
					Status:     "done",
					DoneReason: "filled",
				},
				accounts: []coinbasepro.Account{
					{
						Currency: "USD",
						Balance:  "1000.0551",
					},
				},
			},
			wantNumberOwn:      0.0,
			wantAvailableFunds: 1000.05, //truncated
			wantErr:            nil,
		},
		{
			name: "Sad Path. Error from client returns error.",
			fields: fields{
				wantErr: fmt.Errorf("this is so broke"),
				order: coinbasepro.Order{
					ID:         "GUID-99",
					FilledSize: "99.9999",
					Status:     "done",
					DoneReason: "filled",
				},
				accounts: []coinbasepro.Account{
					{
						Currency: "USD",
						Balance:  "1000.0001",
					},
				},
			},
			wantNumberOwn:      0.0,
			wantAvailableFunds: 0.0,
			wantErr:            fmt.Errorf("this is so broke"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := proclient.NewMockClient()
			c.Err = tt.fields.wantErr
			c.SavedOrder = tt.fields.order
			c.Accounts = tt.fields.accounts
			svc := CoinbaseSvc{
				Client:  c,
				Timeout: time.Duration(time.Millisecond), // for all tests, no backoff necessary
			}
			numberOwn, availablefunds, err := svc.Sell(tt.args.product, tt.args.numberOwn, tt.args.sellPrice)
			if numberOwn != tt.wantNumberOwn {
				t.Errorf("CoinbaseSvc.Sell() numberOwn = %v, want %v", numberOwn, tt.wantNumberOwn)
			}
			if availablefunds != tt.wantAvailableFunds {
				t.Errorf("CoinbaseSvc.Sell() availablefunds = %v, want %v", availablefunds, tt.wantAvailableFunds)
			}

			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CoinbaseSvc.Sell() err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
