package selectionAlgorithms

import "coursework/internal/selectionTypes"

func SelectCoins(params selectionTypes.CoinSelectionParams) *selectionTypes.CoinSelection {
	//Validate that target value is correct

	if params.TargetValue == 0 || params.TargetValue > int(selectionTypes.MaxMoney) {
		ans := selectionTypes.InvalidSpend
		return selectionTypes.NewCoinSelection(params, nil, &ans)
	}

	// check for insufficient funds
	if params.GetTotalValue() < params.TargetValue {
		ans := selectionTypes.InsufficientFunds
		return selectionTypes.NewCoinSelection(params, nil, &ans)
	}

	if params.GetTotalValue() < params.TargetValue+params.GetFixedFee() {
		ans := selectionTypes.InsufficientFundsAfterFees
		return selectionTypes.NewCoinSelection(params, nil, &ans)
	}

	// Return branch and bound selection (more optimized) if possible
	bnbSelection := SelectCoinsBranchAndBounds(params)
	if bnbSelection.Outcome == selectionTypes.Success {
		return bnbSelection
	}

	// Otherwise return knapsack_selection (less optimized) if possible
	knapsackSelection := SelectCoinsKnapsackSolver(params, nil)
	if knapsackSelection.Outcome == selectionTypes.Success {
		return knapsackSelection
	}

	// If all else fails, return single random draw selection (not optomized) as a fallback
	return SelectCoinsSingleRandomDraw(params)
}
