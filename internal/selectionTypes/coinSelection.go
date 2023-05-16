package selectionTypes

type CoinSelection struct {
	Outcome        Outcome
	Outputs        []*InputCoin
	TargetValue    int
	EffectiveValue int
	Value          int
	Fee            int
	ChangeValue    int
}

func NewCoinSelection(params CoinSelectionParams, selectedOutputGroups []OutputGroup, outcome *Outcome) *CoinSelection {
	if outcome == nil {
		tempOutput := Success
		outcome = &tempOutput
	}
	cs := &CoinSelection{
		TargetValue: params.TargetValue,
		Outputs:     make([]*InputCoin, 0),
		Outcome:     *outcome,
	}

	for _, outputGroup := range selectedOutputGroups {
		for _, output := range outputGroup.Outputs {
			cs.Insert(output)
		}
	}

	cs.Fee = cs.CalculateFee(params.GetFixedFee())
	cs.ChangeValue = cs.CalculateChangeValue(params.GetCostOfChange())

	return cs
}

func NewCoinSelectionFromUtxoPool(params CoinSelectionParams, utxoPool []OutputGroup, bestSelection []bool) *CoinSelection {
	selectedGroups := make([]OutputGroup, 0)
	for i, output := range utxoPool {
		if bestSelection[i] {
			selectedGroups = append(selectedGroups, output)
		}
	}
	return NewCoinSelection(params, selectedGroups, nil)
}

func (cs *CoinSelection) Insert(coin *InputCoin) {
	cs.Outputs = append(cs.Outputs, coin)
	cs.EffectiveValue += coin.EffectiveValue
	cs.Value += coin.EffectiveValue
}

func (cs *CoinSelection) CalculateChangeValue(costOfChange int) int {
	if cs.Outcome != Success {
		return 0
	}
	var changeValue int = 0
	if cs.EffectiveValue-cs.TargetValue > 0 && changeValue > costOfChange {
		changeValue = cs.EffectiveValue - cs.TargetValue
	}
	return changeValue
}

func (cs *CoinSelection) CalculateFee(fixedFee int) int {
	if cs.Outcome != Success {
		return 0
	}
	return fixedFee + cs.Value - cs.EffectiveValue
}
