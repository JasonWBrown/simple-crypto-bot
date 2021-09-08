# Simple Momentum Crypto Bot
## In Beta Testing
Use at your own risk! 
Simple Program To Buy and Sell Crypto.
* Possible Buys are hard coded.
* Rules are hard coded. 

# Build

> ## Golang 1.16
> brew install go

# Test
Always run the unit test coverage
> make cover

## HTML Coverage Heat Map
> make html

# Run 
> make run

# Algorithm
- [X] Authenticate
- [X] Get list of possible buys.  Hard Coded. Single Pair
- [X] Get price
- [X] Check price on schedule 
- [X] If price goes up 3% in a 2 hour period buy
- [X] Raise stop order every 1% gain after 3% gain
- [X] Sell at 8% after buy
- [X] Sell if price goes 10% below buy.  

# Test Strategy
- Unit Testing 70% requirement
- Use Mock/Imposter Interfaces where available to test packages in issolation.

# Error Handling
- Errors will percolate to the top level processor.
- Errors at the top will result in no processing.


# Logging pattern
- [ ] What are the logging patterns here? 
- [ ] Daily rotations? 