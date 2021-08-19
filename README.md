# Simple Crypto Bot
## Version 1 

Simple Program To Buy and Sell Crypto.
* Possible Buys are hard coded.
* Rules are hard coded. 

# Build

> ## Golang 1.16
> brew install go

# Run 
> go run main.go

# Algorthim
- [X] Authenticate
- [X] Get list of possible buys.  Hard Coded. Single Pair
- [X] Get price
- [X] Check price on schedule 
- [X] If price goes up 3% in a 2 hour period buy
- [X] Raise stop order every 1% gain after 3% gain
- [X] Sell at 8% after buy
- [X] Sell if price goes 10% below buy.  

- [ ] Figure out test strategy
- [ ] Error Handling

