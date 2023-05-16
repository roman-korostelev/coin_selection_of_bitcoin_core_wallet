package selectionTypes

type OutputGroup struct {
	Outputs        []*InputCoin
	Value          int
	EffectiveValue int
	Fee            int
	LongTermFee    int
	Address        string
}

func (og *OutputGroup) SetFee(shortTermFeeRate int, longTermFeeRate int) {
	og.EffectiveValue = 0
	og.Fee = 0
	og.LongTermFee = 0

	nonNegativeOutputs := make([]*InputCoin, 0)

	for _, coin := range og.Outputs {
		coin.SetFee(shortTermFeeRate, longTermFeeRate)
		if coin.EffectiveValue > 0 {
			og.Fee += coin.Fee
			og.LongTermFee += coin.LongTermFee
			og.EffectiveValue += coin.EffectiveValue
			nonNegativeOutputs = append(nonNegativeOutputs, coin)
		}
	}

	og.Outputs = nonNegativeOutputs
}

func (og *OutputGroup) InsertCoin(coin *InputCoin) {
	coin.Address = og.Address
	og.Outputs = append(og.Outputs, coin)
	og.Value += coin.Value
}
