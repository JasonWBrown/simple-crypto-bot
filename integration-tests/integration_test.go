package integrationtests_test

import (
	"fmt"
	"testing"

	"github.com/walkerus/go-wiremock"
)

func TestBuyAt03GrowthRateAndSellAtLossOf10Percent(t *testing.T) {
	//given there exists apis in wiremock
	wiremockClient := wiremock.NewClient("http://0.0.0.0:8080")

	//when the 2 hour historical data shows a .03% gain
	wiremockClient.StubFor(wiremock.Get(wiremock.URLPathMatching("/products/BTC-USD/candles")).
		WithQueryParam("end", wiremock.Contains("-")).
		WithQueryParam("start", wiremock.Contains("-")).
		WillReturn(
			`[`+
				// [ time, low, high, open, close, volume ],
				`[ 1415398769, 0.32, 4.2, 0.35, 103, 12.3 ],`+
				`[ 1415398768, 0.32, 4.2, 0.35, 4.2, 12.3 ],`+
				`[ 1415398767, 0.32, 4.2, 0.35, 100, 12.3 ]`+
				`]`,
			map[string]string{"Content-Type": "application/json"},
			200,
		).
		AtPriority(1))

	fmt.Println("running...")

	//and there is a valid ticker response

	//then there will be a buy request in wiremock

	// and there is a loss of 10 percent

	// there will be a sell order for the 10 percent loss
}

func VerifyBuyAt03GrowthRateAndSellAfterLockingIn3PercentGrowth(t *testing.T) {
	//given there exists apis in wiremock

	//when the 2 hour historical data shows a .03% gain
	//and there is a valid ticker response

	//then there will be a buy request in wiremock

	//and there is a gain of 3 percent

	//when there is a loss of 10 percent

	//there will be a sell order for the 3 percent gain
}

func VerifyBuyAt03GrowthRateAndSellAfterLockingIn3PercentGrowthPlusAlmost1percent(t *testing.T) {
	//given there exists apis in wiremock

	//when the 2 hour historical data shows a .03% gain
	//and there is a valid ticker response

	//then there will be a buy request in wiremock

	//and there is a gain of 3 percent

	//and there is a gain of an additional .9 percent

	//when there is a loss of 10 percent

	//there will be a sell order for the 4 percent gain
}

func VerifyBuyAt03GrowthRateAndSellAfterLockingIn3PercentGrowthPlus1percent(t *testing.T) {
	//given there exists apis in wiremock

	//when the 2 hour historical data shows a .03% gain
	//and there is a valid ticker response

	//then there will be a buy request in wiremock

	//and there is a gain of 3 percent

	//and there is a gain of an additional 1 percent

	//when there is a loss of 10 percent

	//there will be a sell order for the 4 percent gain
}

func VerifyBuyAt03GrowthRateAndSellAt8Percent(t *testing.T) {
	//given there exists apis in wiremock

	//when the 2 hour historical data shows a .03% gain
	//and there is a valid ticker response

	//then there will be a buy request in wiremock

	//and there is a gain of 8 percent

	//there will be a sell order for the 8 percent gain
}
