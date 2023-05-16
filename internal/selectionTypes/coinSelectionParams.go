package selectionTypes

type CoinSelectionParams struct {
	UtxoPool                []OutputGroup
	TargetValue             int
	ShortTermFeeRate        int
	LongTermFeeRate         int
	ChangeOutputSizeInBytes int
	ChangeSpendSizeInBytes  int
	NotInputSizeInBytes     int
}

func (p *CoinSelectionParams) Init() {
	for _, output := range p.UtxoPool {
		output.SetFee(p.ShortTermFeeRate, p.LongTermFeeRate)
	}
}

func (p *CoinSelectionParams) GetTotalValue() (ans int) {
	for _, output := range p.UtxoPool {
		ans += output.Value
	}
	return ans
}

func (p *CoinSelectionParams) GetTotalEffectiveValue() (ans int) {
	for _, output := range p.UtxoPool {
		ans += output.EffectiveValue
	}
	return ans
}

func (p *CoinSelectionParams) GetFixedFee() int {
	return p.ShortTermFeeRate * p.NotInputSizeInBytes
}

func (p *CoinSelectionParams) GetCostOfCreatingChange() int {
	return p.ShortTermFeeRate * p.ChangeOutputSizeInBytes
}

func (p *CoinSelectionParams) GetCostOfSpendingChange() int {
	return p.LongTermFeeRate * p.ChangeSpendSizeInBytes
}

func (p *CoinSelectionParams) GetCostOfChange() int {
	return p.GetCostOfSpendingChange() + p.GetCostOfCreatingChange()
}
