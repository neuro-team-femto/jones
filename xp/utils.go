package xp

import (
	"math/rand"
	"time"
)

func truncatedInPlaceShuffle(input []string, max int) []string {
	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(input), func(i, j int) {
		input[i], input[j] = input[j], input[i]
	})
	return input[:max]
}
