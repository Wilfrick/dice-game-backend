package game

import (
	"testing"
)

func Test_generateNewHands(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{7, 7, 7, 7}), PlayerHand([]int{7, 7, 7}), PlayerHand([]int{7, 7, 7})}
	gameState.randomiseCurrentHands()

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

func Test_isBetTrueSimpleCase(t *testing.T) {
	var gameState GameState
	gameState.PrevMove = PlayerMove{MoveType: "Bet", Value: Bet{5, 2}}
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 2}), PlayerHand([]int{2, 2})} //Exactly true
	if !gameState.isBetTrue() {
		t.Fail()
	}
	if !gameState.isBetExactlyTrue() {
		t.Fail()
	}
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 2}), PlayerHand([]int{3, 3})} //Not true
	if gameState.isBetTrue() {
		t.Fail()
	}
	if gameState.isBetExactlyTrue() {
		t.Fail()
	}
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 2}), PlayerHand([]int{2, 2, 2})} //True but not exact
	if !gameState.isBetTrue() {
		t.Fail()
	}
	if gameState.isBetExactlyTrue() {
		t.Fail()
	}
}

func Test_isBetTrueOnesCase(t *testing.T) {
	var gameState GameState
	gameState.PrevMove = PlayerMove{MoveType: "Bet", Value: Bet{5, 2}}
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 2}), PlayerHand([]int{1, 1})} //Exactly true
	if !gameState.isBetTrue() {
		t.Fail()
	}
	if !gameState.isBetExactlyTrue() {
		t.Fail()
	}
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 2}), PlayerHand([]int{1, 1}), PlayerHand([]int{1, 3})} //True not exactly
	if !gameState.isBetTrue() {
		t.Fail()
	}
	if gameState.isBetExactlyTrue() {
		t.Fail()
	}
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{2, 4, 4}), PlayerHand([]int{1, 1}), PlayerHand([]int{1, 3})} //Not true
	if gameState.isBetTrue() {
		t.Fail()
	}
	if gameState.isBetExactlyTrue() {
		t.Fail()
	}
}
