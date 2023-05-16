package selectionTypes

type Money int

var (
	Coin      Money = 1000000
	Cent            = Coin / 100
	MinChange       = Cent
	MaxMoney        = 21000000 * Coin
)

type Outcome int

const (
	Success Outcome = iota
	InsufficientFunds
	InsufficientFundsAfterFees
	AlgorithmFailure
	InvalidSpend
)
