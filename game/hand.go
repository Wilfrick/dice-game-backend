package game

import (
	"HigherLevelPerudoServer/messages"
	"math/rand"
)

type PlayerHand []int

func RandomPlayerHand(length int) PlayerHand {
	hand := make([]int, length)
	for i := range hand {
		hand[i] = rand.Intn(6) + 1
	}
	return hand
}

func (playerHand *PlayerHand) Randomise() {
	copy(*playerHand, RandomPlayerHand(len([]int(*playerHand))))
}

func (playerHand PlayerHand) AssembleHandMessage() []byte {
	encodedMessage := messages.PackMessage("PlayerHand", playerHand)
	return encodedMessage
}
