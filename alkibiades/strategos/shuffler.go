package strategos

// ShuffleAndDistributeSteps takes a slice of RandomSteps and returns a new slice with steps shuffled
// and similar types distributed to avoid having the same type in consecutive positions
func (j *JourneyServiceImpl) shuffleAndDistributeSteps(steps []RandomSteps) []RandomSteps {
	// If we have 0 or 1 steps, no shuffling needed
	if len(steps) <= 1 {
		return steps
	}

	// Create a copy of the steps to avoid modifying the original
	shuffled := make([]RandomSteps, len(steps))
	copy(shuffled, steps)

	// Step 1: Shuffle the steps initially
	for i := range shuffled {
		k := j.Randomizer.RandomNumberBaseZero(i + 1)
		shuffled[i], shuffled[k] = shuffled[k], shuffled[i]
	}

	// Step 2: Try to distribute similar types
	result := make([]RandomSteps, 0, len(steps))
	typeUsed := make(map[string]bool)

	// Add the first step to the result
	result = append(result, shuffled[0])
	typeUsed[shuffled[0].Type] = true

	// Remove the first step from available steps
	available := make([]RandomSteps, 0, len(steps)-1)
	for i := 1; i < len(shuffled); i++ {
		available = append(available, shuffled[i])
	}

	// While we have steps to place
	for len(available) > 0 {
		lastType := result[len(result)-1].Type

		// Find a step with a different type than the last one added
		differentTypeIndex := -1
		for i, step := range available {
			if step.Type != lastType {
				differentTypeIndex = i
				break
			}
		}

		// If we found a different type, use it
		if differentTypeIndex != -1 {
			result = append(result, available[differentTypeIndex])
			// Remove the used step
			available = append(available[:differentTypeIndex], available[differentTypeIndex+1:]...)
		} else {
			// If we couldn't find a different type, just use the first available
			result = append(result, available[0])
			available = available[1:]
		}
	}

	// Additional optimization - try to separate identical types further
	optimized := j.improveTypeDistribution(result)

	return optimized
}

// improveTypeDistribution tries to improve the distribution of types by looking for
// clusters of the same type and breaking them up if possible
func (j *JourneyServiceImpl) improveTypeDistribution(steps []RandomSteps) []RandomSteps {
	if len(steps) <= 2 {
		return steps
	}

	// Create a copy of the input
	result := make([]RandomSteps, len(steps))
	copy(result, steps)

	// Try a few iterations of improvement
	for attempt := 0; attempt < 3; attempt++ {
		improved := false

		// Look through all steps (except last one)
		for i := 0; i < len(result)-1; i++ {
			// If current type matches next type
			if result[i].Type == result[i+1].Type {
				// Try to find a different type further in the sequence
				swapIndex := -1
				for k := i + 2; k < len(result); k++ {
					// Make sure swapping wouldn't create another cluster
					if result[k].Type != result[i].Type &&
						(k == len(result)-1 || result[k].Type != result[k+1].Type) &&
						(i == 0 || result[k].Type != result[i-1].Type) {
						swapIndex = k
						break
					}
				}

				// If we found a suitable swap, do it
				if swapIndex != -1 {
					result[i+1], result[swapIndex] = result[swapIndex], result[i+1]
					improved = true
				}
			}
		}

		// If we couldn't improve further, stop trying
		if !improved {
			break
		}
	}

	// In some rare cases with many of the same type, we might still have clusters
	// Let's do a final check and apply a random swap if needed
	for i := 0; i < len(result)-1; i++ {
		if result[i].Type == result[i+1].Type {
			// Just try a random swap to break this up
			k := j.Randomizer.RandomNumberBaseZero(len(result))
			if k != i && k != i+1 {
				result[i+1], result[k] = result[k], result[i+1]
			}
		}
	}

	return result
}

func (j *JourneyServiceImpl) randomizeOptions(options []string, answer string) []string {
	// Make a copy to avoid modifying the original
	optionsCopy := make([]string, len(options))
	copy(optionsCopy, options)

	// Shuffle the options
	j.ShuffleStrings(optionsCopy)

	// Ensure the answer is in the options
	answerIncluded := false
	for _, opt := range optionsCopy {
		if opt == answer {
			answerIncluded = true
			break
		}
	}

	// If answer is not included (which shouldn't happen but just in case),
	// replace the last option with the answer
	if !answerIncluded {
		optionsCopy[len(optionsCopy)-1] = answer
	}

	return optionsCopy
}

func (j *JourneyServiceImpl) limitAndRandomizeOptions(options []string, answer string, maxOptions int) []string {
	// First check if we need to do anything
	if len(options) <= maxOptions {
		return j.randomizeOptions(options, answer)
	}

	// Create a map to track if answer is included
	includedOptions := make(map[string]bool)
	includedOptions[answer] = true

	// Initialize result with the answer
	result := []string{answer}

	// Create a copy of options without the answer
	var optionsWithoutAnswer []string
	for _, opt := range options {
		if opt != answer {
			optionsWithoutAnswer = append(optionsWithoutAnswer, opt)
		}
	}

	// Shuffle the remaining options
	j.ShuffleStrings(optionsWithoutAnswer)

	// Add options until we reach maxOptions
	for _, opt := range optionsWithoutAnswer {
		if len(result) >= maxOptions {
			break
		}
		result = append(result, opt)
	}

	// Shuffle the final result
	j.ShuffleStrings(result)

	return result
}

func (j *JourneyServiceImpl) ShuffleStrings(slice []string) {
	for i := range slice {
		j := j.Randomizer.RandomNumberBaseZero(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func (j *JourneyServiceImpl) ShuffleRandomSteps(steps []RandomSteps) {
	for i := range steps {
		j := j.Randomizer.RandomNumberBaseZero(i + 1)
		steps[i], steps[j] = steps[j], steps[i]
	}
}
