package selectionAlgorithms

import (
	"coursework/internal/selectionTypes"
	"sort"
)

const totalTries = 10

func SelectCoinsBranchAndBounds(params selectionTypes.CoinSelectionParams) *selectionTypes.CoinSelection {
	utxoPool := params.UtxoPool
	targetAfterFixedFees := params.TargetValue + params.GetFixedFee()
	curValue := 0
	curSelection := make([]bool, 0)
	curAvailableValue := params.GetTotalEffectiveValue()
	curWaste := 0
	bestWaste := selectionTypes.MaxMoney
	bestSelection := make([]bool, 0)
	sort.Slice(utxoPool, func(i, j int) bool {
		return utxoPool[i].EffectiveValue < utxoPool[j].EffectiveValue
	})

	for i := 0; i < totalTries; i++ {
		// Cannot possibly reach target with the amount remaining in the curr_available_value
		a := curValue+curAvailableValue < targetAfterFixedFees
		// Selected value is out of range, go back and try other branch
		b := curValue > targetAfterFixedFees+params.GetCostOfChange()
		// Don't select things which we know will be more wasteful if the waste is increasing
		c := curWaste > int(bestWaste) && utxoPool[0].Fee-utxoPool[0].LongTermFee > 0

		shouldBacktrack := a || b || c
		if !shouldBacktrack && curValue >= targetAfterFixedFees {
			/*This is the excess value which is added to the waste for the below comparison
			Adding another UTXO after this check could bring the waste down if the long term fee is higher than the current fee.
			However we are not going to explore that because this optimization for the waste is only done when we have hit our target
			value. Adding any more UTXOs will be just burning the UTXO; it will go entirely to fees. Thus we aren't going to
			explore any more UTXOs to avoid burning money like that */

			curWaste += curValue - targetAfterFixedFees
			if curWaste <= int(bestWaste) {
				bestSelection = curSelection
				bestWaste = selectionTypes.Money(curWaste)
				if bestWaste == 0 {
					return selectionTypes.NewCoinSelectionFromUtxoPool(params, params.UtxoPool, bestSelection)
				}
			}
			// Remove the excess value as we will be selecting different coins now
			curWaste -= curValue - targetAfterFixedFees
			shouldBacktrack = true
		}
		if shouldBacktrack {
			/* Walk backwards to find the last included UTXO
			that still needs to have its omission branch traversed*/
			for len(curSelection) > 0 && !curSelection[len(curSelection)-1] {
				curSelection = curSelection[:len(curSelection)-1]
				curAvailableValue += utxoPool[len(curSelection)].EffectiveValue
			}
			if len(curSelection) == 0 {
				break
			}
		} else {
			utxo := utxoPool[len(curSelection)]
			curAvailableValue -= utxo.EffectiveValue
			/* Avoid searching a branch if the previous UTXO has the same value and
			same waste and was excluded. Since the ratio of fee to long term fee
			is the same, we only need to check if one of those values match in
			order to know that the waste is the same */
			if len(curSelection) > 0 && !curSelection[len(curSelection)-1] &&
				utxo.EffectiveValue == utxoPool[len(curSelection)-1].EffectiveValue &&
				utxo.Fee == utxoPool[len(curSelection)-1].Fee {
				curSelection = append(curSelection, false)
			} else {
				// Inclusion branch first (Largest First Exploration)
				curSelection = append(curSelection, true)
				curValue += utxo.EffectiveValue
				curWaste += utxo.Fee - utxo.LongTermFee
			}
		}
	}

	if len(bestSelection) == 0 {
		ans := selectionTypes.AlgorithmFailure
		return selectionTypes.NewCoinSelection(params, nil, &ans)
	}
	return selectionTypes.NewCoinSelectionFromUtxoPool(params, params.UtxoPool, bestSelection)
}
