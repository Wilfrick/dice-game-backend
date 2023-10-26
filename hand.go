package main

import "math/rand"

type PlayerHand []int

func randomPlayerHand(length int) PlayerHand {
	hand := make([]int, length)
	for i := range hand {
		hand[i] = rand.Intn(6) + 1
	}
	return hand
}
