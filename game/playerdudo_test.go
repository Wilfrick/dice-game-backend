package game

import "testing"

func Test_generateNewHands(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{7, 7, 7, 7}), PlayerHand([]int{7, 7, 7}), PlayerHand([]int{7, 7, 7})}
	gameState.generateNewHands()

	lengths := []int{4, 3, 3}
	for i, hand := range gameState.PlayerHands {
		hand_length := len(hand)
		if hand_length != lengths[i] {
			t.Error("Generated hand of inappropriate length")
		}
		for _, val := range hand {
			if (val < 1) || (val > 6) {
				t.Error("Generated dice with non dice value")
			}
		}
	}
}
