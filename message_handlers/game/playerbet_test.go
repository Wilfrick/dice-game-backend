package game

import (
	"HigherLevelPerudoServer/util"
	"testing"
)

// The following tests how we update player turn
func Test_updatePlayerIndexRunsEmptyPlayerHands(t *testing.T) {
	var gamestate GameState

	err := gamestate.updatePlayerIndex(BET) //Expecting success to be false
	if err == nil {
		t.Fail()
	}

}

func Test_updatePlayerIndexRunsNonEmptyPlayerHands(t *testing.T) {
	var gamestate GameState

	gamestate.PlayerHands = []PlayerHand{PlayerHand([]int{1, 3, 4, 5}), PlayerHand([]int{2, 4, 4}), PlayerHand([]int{4, 5, 4})}

	err := gamestate.updatePlayerIndex(BET) //Expecting success to be false
	if !(err == nil) {
		t.Fail()
	}

}

func Test_checkPlayerIndexIncrementsClean(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{1, 3, 4, 5}), PlayerHand([]int{2, 4, 4}), PlayerHand([]int{4, 5, 4})}
	gameState.CurrentPlayerIndex = 0
	gameState.PrevMove = PlayerMove{MoveType: "Bet", Value: Bet{NumDice: 2, FaceVal: 2}}

	gameState.updatePlayerIndex(BET)

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

	gameState.updatePlayerIndex(BET)

	expectedNewPlayerIndex := 0
	if !(gameState.CurrentPlayerIndex == expectedNewPlayerIndex) {
		t.Fail()
	}

}

func Test_checkPlayerIndexIncrementsDeadPlayer(t *testing.T) {
	var gameState GameState
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{1, 3, 4, 5}), PlayerHand([]int{}), PlayerHand([]int{4, 5, 4})}
	gameState.CurrentPlayerIndex = 0
	gameState.PrevMove = PlayerMove{MoveType: "Bet", Value: Bet{NumDice: 2, FaceVal: 2}}

	gameState.updatePlayerIndex(BET)

	expectedNewPlayerIndex := 2
	if !(gameState.CurrentPlayerIndex == expectedNewPlayerIndex) {
		t.Fail()
	}

}

func Test_checkPlayerIndexAllPlayersDead(t *testing.T) {

	var gameState GameState
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{}), PlayerHand([]int{}), PlayerHand([]int{})}
	gameState.CurrentPlayerIndex = 0
	gameState.PrevMove = PlayerMove{MoveType: "Bet", Value: Bet{NumDice: 2, FaceVal: 2}}

	err := gameState.updatePlayerIndex(BET) //We expect to fail
	if !(err.Error() == "all players are dead") {
		t.Fail()
	}
	if err == nil {
		t.Fail()
	}
}

func Test_checkPlayerIndexSinglePlayerAlive(t *testing.T) {

	var gameState GameState
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{5, 5}), PlayerHand([]int{}), PlayerHand([]int{})}
	gameState.CurrentPlayerIndex = 0
	gameState.PrevMove = PlayerMove{MoveType: "Bet", Value: Bet{NumDice: 2, FaceVal: 2}}

	err := gameState.updatePlayerIndex(BET) //We expect to fail
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
}

func Test_updatePlayerIndexDudoTrue(t *testing.T) {
	var gs GameState
	gs.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	gs.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 3}), PlayerHand([]int{1}), PlayerHand([]int{5})}
	gs.CurrentPlayerIndex = 1
	losing_player_index := 0
	err := gs.updatePlayerIndex(DUDO, losing_player_index)
	if err != nil {
		t.Fail()
	}
	t.Log(gs.CurrentPlayerIndex)
	util.Assert(t, gs.CurrentPlayerIndex == 0)
}

func Test_updatePlayerIndexDudoFalse(t *testing.T) {
	var gs GameState
	gs.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 3}), PlayerHand([]int{1, 2}), PlayerHand([]int{5})}
	gs.CurrentPlayerIndex = 1
	losing_player_index := 1
	err := gs.updatePlayerIndex(DUDO, losing_player_index)
	if err != nil {
		t.Fail()
	}
	t.Log(gs.CurrentPlayerIndex)
	util.Assert(t, gs.CurrentPlayerIndex == 1)
}

func Test_updatePlayerIndexCalza(t *testing.T) {
	var gs GameState
	gs.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	gs.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 3}), PlayerHand([]int{1, 2}), PlayerHand([]int{5})}
	gs.CurrentPlayerIndex = 1
	losing_player_index := 1
	err := gs.updatePlayerIndex(CALZA, losing_player_index)
	if err != nil {
		t.Fail()
	}
	t.Log(gs.CurrentPlayerIndex)
	util.Assert(t, gs.CurrentPlayerIndex == 1)
}
