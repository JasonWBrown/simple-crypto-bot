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
			name: "Sed Path. Client does not return error but there are no bids will return error.",
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
