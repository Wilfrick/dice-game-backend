package game

import "testing"

func Test_handRandomise(t *testing.T) {
	playerHand := PlayerHand{0, 0, 0, 0}
	playerHand.Randomise()
	if playerHand[0] == 0 {
		t.Error("Failed to set a valid dice face value")
	}
}
