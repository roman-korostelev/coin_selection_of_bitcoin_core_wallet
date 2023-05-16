package selectionAlgorithms

import (
	"coursework/internal/selectionTypes"
	"math/rand"
)

func SelectCoinsSingleRandomDraw(params selectionTypes.CoinSelectionParams) *selectionTypes.CoinSelection {
	targetAfterFixedFee := params.TargetValue + params.GetFixedFee()
	utxoPool := params.UtxoPool

	rand.Shuffle(len(utxoPool), func(i, j int) {
		utxoPool[i], utxoPool[j] = utxoPool[j], utxoPool[i]
	})
	selectedOutputGroups := make([]selectionTypes.OutputGroup, 0)
	selectedValue := 0
	for _, outputGroup := range params.UtxoPool {
		selectedValue += outputGroup.EffectiveValue
		selectedOutputGroups = append(selectedOutputGroups, outputGroup)
		if selectedValue >= targetAfterFixedFee {
			return selectionTypes.NewCoinSelection(params, selectedOutputGroups, nil)
		}
	}
	ans := selectionTypes.AlgorithmFailure
	return selectionTypes.NewCoinSelection(params, nil, &ans)
}
