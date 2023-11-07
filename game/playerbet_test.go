package game

import "testing"

// The following tests how we update player turn
func Test_updatePlayerIndexRunsEmptyPlayerHands(t *testing.T) {
	var gamestate GameState
	var newmove PlayerMove

	err := gamestate.updatePlayerIndex(newmove) //Expecting success to be false
	if err == nil {
		t.Fail()
	}

}

func Test_updatePlayerIndexRunsNonEmptyPlayerHands(t *testing.T) {
	var gamestate GameState
	var newmove PlayerMove
	gamestate.PlayerHands = []PlayerHand{PlayerHand([]int{1, 3, 4, 5}), PlayerHand([]int{2, 4, 4}), PlayerHand([]int{4, 5, 4})}

	err := gamestate.updatePlayerIndex(newmove) //Expecting success to be false
	if !(err == nil) {
		t.Fail()
	}

}

func Test_checkPlayerIndexIncrementsClean(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{1, 3, 4, 5}), PlayerHand([]int{2, 4, 4}), PlayerHand([]int{4, 5, 4})}
	gameState.CurrentPlayerIndex = 0
	gameState.PrevMove = PlayerMove{MoveType: "Bet", Value: Bet{NumDice: 2, FaceVal: 2}}
	newBet := Bet{NumDice: 3, FaceVal: 2}

	gameState.updatePlayerIndex(PlayerMove{MoveType: "Bet", Value: newBet})

	expectedNewPlayerIndex := 1
	if !(gameState.CurrentPlayerIndex == expectedNewPlayerIndex) {
		t.Fail()
	}

}

func Test_checkPlayerIndexIncrementsWrapArround(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{1, 3, 4, 5}), PlayerHand([]int{2, 4, 4}), PlayerHand([]int{4, 5, 4})}
	gameState.CurrentPlayerIndex = 2
	gameState.PrevMove = PlayerMove{MoveType: "Bet", Value: Bet{NumDice: 2, FaceVal: 2}}
	newBet := Bet{NumDice: 3, FaceVal: 2}

	gameState.updatePlayerIndex(PlayerMove{MoveType: "Bet", Value: newBet})

	expectedNewPlayerIndex := 0
	if !(gameState.CurrentPlayerIndex == expectedNewPlayerIndex) {
		t.Fail()
	}

}

func Test_checkPlayerIndexIncrementsDeadPlayer(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{1, 3, 4, 5}), PlayerHand([]int{}), PlayerHand([]int{4, 5, 4})}
	gameState.CurrentPlayerIndex = 0
	gameState.PrevMove = PlayerMove{MoveType: "Bet", Value: Bet{NumDice: 2, FaceVal: 2}}
	newBet := Bet{NumDice: 3, FaceVal: 2}

	gameState.updatePlayerIndex(PlayerMove{MoveType: "Bet", Value: newBet})

	expectedNewPlayerIndex := 2
	if !(gameState.CurrentPlayerIndex == expectedNewPlayerIndex) {
		t.Fail()
	}

}

func Test_checkPlayerIndexAllPlayersDead(t *testing.T) {

	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{}), PlayerHand([]int{}), PlayerHand([]int{})}
	gameState.CurrentPlayerIndex = 0
	gameState.PrevMove = PlayerMove{MoveType: "Bet", Value: Bet{NumDice: 2, FaceVal: 2}}
	newBet := Bet{NumDice: 3, FaceVal: 2}

	err := gameState.updatePlayerIndex(PlayerMove{MoveType: "Bet", Value: newBet}) //We expect to fail
	if !(err.Error() == "all players are dead") {
		t.Fail()
	}
	if err == nil {
		t.Fail()
	}
}

func Test_checkPlayerIndexSinglePlayerAlive(t *testing.T) {

	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{5, 5}), PlayerHand([]int{}), PlayerHand([]int{})}
	gameState.CurrentPlayerIndex = 0
	gameState.PrevMove = PlayerMove{MoveType: "Bet", Value: Bet{NumDice: 2, FaceVal: 2}}
	newBet := Bet{NumDice: 3, FaceVal: 2}

	err := gameState.updatePlayerIndex(PlayerMove{MoveType: "Bet", Value: newBet}) //We expect to fail
	if !(err.Error() == "looped around to our initial player. all other players dead") {
		t.Fail()
	}
	if err == nil {
		t.Fail()
	}
}
