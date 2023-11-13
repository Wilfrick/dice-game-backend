package game

import (
	"HigherLevelPerudoServer/util"
	"testing"
)

func Test_RemoveDice(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{1, 2, 3}), PlayerHand([]int{5, 4, 1}), PlayerHand([]int{4, 5, 6})}

	PLAYER_INDEX := 1
	ORIGINAL_LENGTH := len(gameState.PlayerHands[PLAYER_INDEX])
	death, err := gameState.removeDice(PLAYER_INDEX)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
	if death {
		t.Fail()
	}
	util.Assert(t, len(gameState.PlayerHands[PLAYER_INDEX]) != ORIGINAL_LENGTH)
}

func Test_RemoveDiceKilling(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{1, 2, 3}), PlayerHand([]int{5}), PlayerHand([]int{4, 5, 6})}
	PLAYER_INDEX := 1
	death, err := gameState.removeDice(PLAYER_INDEX)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
	if !death {
		t.Fail()
	}
}
