package selectionAlgorithms

import (
	"coursework/internal/selectionTypes"
	"math/rand"
	"sort"
)

const defaultIterations = 1000

func ApproximateBestSubset(params selectionTypes.CoinSelectionParams, utxoPool []selectionTypes.OutputGroup,
	totalLower int, iterations int, adjustForMinChange bool) *selectionTypes.CoinSelection {
	targetAfterFixedFee := params.TargetValue + params.GetFixedFee()
	if adjustForMinChange {
		targetAfterFixedFee += int(selectionTypes.MinChange)
	}

	bestSelection := make([]bool, len(utxoPool))
	for i := 0; i < len(bestSelection); i++ {
		bestSelection[i] = true
	}
	bestValue := totalLower
	var reachedTarget bool
	for i := 0; i < iterations; i++ {
		if bestValue == targetAfterFixedFee {
			break
		}
		included := make([]bool, len(utxoPool))
		for j := 0; j < len(included); j++ {
			included[j] = false
		}
		totalValue := 0
		reachedTarget = false
		numberOfPasses := 2

		for passNumber := 0; passNumber < numberOfPasses && !reachedTarget; passNumber++ {
			for j := 0; j < len(utxoPool); j++ {
				/* The solver here uses a randomized algorithm,
				the randomness serves no real security purpose but is just
				needed to prevent degenerate behavior and it is important
				that the rng is fast. We do not use a constant random sequence,
				because there may be some privacy improvement by making
				the selection random. */
				var a bool
				if passNumber == 0 {
					a = rand.Float32() > 0.5
				} else {
					a = !included[j]
				}
				if a {
					totalValue += utxoPool[j].EffectiveValue
					included[j] = true
					if totalValue >= targetAfterFixedFee {
						reachedTarget = true
						if totalValue < bestValue {
							bestValue = totalValue
							bestSelection = included
						}
						totalValue -= utxoPool[j].EffectiveValue
						included[j] = false
					}
				}
			}
		}
	}
	if reachedTarget {
		return selectionTypes.NewCoinSelectionFromUtxoPool(params, utxoPool, bestSelection)
	} else {
		ans := selectionTypes.AlgorithmFailure
		return selectionTypes.NewCoinSelection(params, nil, &ans)
	}
}

func SelectCoinsKnapsackSolver(params selectionTypes.CoinSelectionParams, iterations *int) *selectionTypes.CoinSelection {
	if iterations == nil {
		temp := defaultIterations
		iterations = &temp
	}

	utxoPool := params.UtxoPool
	targetAfterFixedFee := params.TargetValue + params.GetFixedFee()

	// lowest output group larger than target_value
	var lowestLarger *selectionTypes.OutputGroup = nil
	applicableGroups := make([]selectionTypes.OutputGroup, 0)
	totalLower := 0

	rand.Shuffle(len(utxoPool), func(i, j int) {
		utxoPool[i], utxoPool[j] = utxoPool[j], utxoPool[i]
	})

	for _, outputGroup := range utxoPool {
		if outputGroup.EffectiveValue == targetAfterFixedFee {
			return selectionTypes.NewCoinSelection(params, []selectionTypes.OutputGroup{outputGroup}, nil)
		}

		if outputGroup.EffectiveValue < targetAfterFixedFee+int(selectionTypes.MinChange) {
			applicableGroups = append(applicableGroups, outputGroup)
			totalLower += outputGroup.EffectiveValue
		} else if lowestLarger == nil || outputGroup.EffectiveValue < lowestLarger.EffectiveValue {
			lowestLarger = &outputGroup
		}
	}
	if totalLower == targetAfterFixedFee {
		return selectionTypes.NewCoinSelection(params, applicableGroups, nil)
	}

	if totalLower < targetAfterFixedFee && lowestLarger != nil {
		return selectionTypes.NewCoinSelection(params, []selectionTypes.OutputGroup{*lowestLarger}, nil)
	}

	sort.Slice(utxoPool, func(i, j int) bool {
		return utxoPool[i].EffectiveValue < utxoPool[j].EffectiveValue
	})
	bestSelection := ApproximateBestSubset(params, applicableGroups, totalLower, *iterations, false)
	if bestSelection.EffectiveValue != targetAfterFixedFee &&
		totalLower >= targetAfterFixedFee+int(selectionTypes.MinChange) {
		bestSelection = ApproximateBestSubset(params, applicableGroups, totalLower, *iterations, true)
	}
	/* If we have a bigger coin and (either the stochastic approximation didn't find
	a good solution, or the next bigger coin is closer), return the bigger coin*/
	if lowestLarger != nil {
		a := bestSelection.EffectiveValue != targetAfterFixedFee &&
			bestSelection.EffectiveValue < targetAfterFixedFee+int(selectionTypes.MinChange)
		b := lowestLarger.EffectiveValue <= bestSelection.EffectiveValue
		if a || b {
			bestSelection = selectionTypes.NewCoinSelection(params, []selectionTypes.OutputGroup{*lowestLarger}, nil)
		}
	}
	return bestSelection
}
