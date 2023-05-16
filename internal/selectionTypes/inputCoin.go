package selectionTypes

import "fmt"

type InputCoin struct {
	TxHash         string
	Vout           int
	Value          int
	InputBytes     int
	EffectiveValue int
	Fee            int
	LongTermFee    int
	Address        string
}

func (c *InputCoin) String() string {
	return fmt.Sprintf("txHash: %s\nvalue: %d", c.TxHash, c.Value)
}

func (c *InputCoin) SetFee(shortTermFeeRate int, longTermFeeRate int) {
	c.Fee = c.InputBytes * shortTermFeeRate
	c.LongTermFee = c.InputBytes * longTermFeeRate
	c.EffectiveValue = c.Value - c.Fee
}
