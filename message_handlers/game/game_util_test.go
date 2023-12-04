package game

import (
	"HigherLevelPerudoServer/util"
	"slices"
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
		t.FailNow() // 🤔
	}
	if death {
		t.Fail()
	}
	util.Assert(t, len(gameState.PlayerHands[PLAYER_INDEX]) == ORIGINAL_LENGTH-1) // ✓
	// could check that other player hands are still intact
	util.Assert(t, len(gameState.PlayerHands[0]) == 3 && len(gameState.PlayerHands[2]) == 3)
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

	util.Assert(t, len(gameState.PlayerHands[1]) == 0)
	// could also assert the player hands after this.
	util.Assert(t, len(gameState.PlayerHands[0]) == 3 && len(gameState.PlayerHands[2]) == 3)
}

func Test_nextPlayerAliveAlivePlayers(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{1, 2, 3}), PlayerHand([]int{5}), PlayerHand([]int{4, 5, 6})}
	gameState.CurrentPlayerIndex = 0
	err := gameState.findNextAlivePlayerInclusive()
	if err != nil {
		t.Fail()
	}
	t.Log(gameState.CurrentPlayerIndex)
	util.Assert(t, gameState.CurrentPlayerIndex == 0)
}
func Test_nextPlayerAliveDeadPlayer(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{}), PlayerHand([]int{5}), PlayerHand([]int{4, 5, 6})}
	gameState.CurrentPlayerIndex = 0
	err := gameState.findNextAlivePlayerInclusive()
	if err != nil {
		t.Fail()
	}
	t.Log(gameState.CurrentPlayerIndex)
	util.Assert(t, gameState.CurrentPlayerIndex == 1)
}

func Test_removePlayer(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{2}), PlayerHand([]int{5}), PlayerHand([]int{4, 5, 6})}
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	gameState.CurrentPlayerIndex = 0
	err := gameState.removePlayer(1)
	if err != nil {
		t.Error("Couldn't remove the player")
	}
	// util.Assert(t, slices.Equal(gameState.PlayerHands,[]PlayerHand{PlayerHand([]int{2}),  PlayerHand([]int{4, 5, 6})}))
	// util.Assert(t, len(gameState.PlayerChannels) == 2)
	util.Assert(t, len(gameState.PlayerHands) == 2)
	util.Assert(t, slices.Equal(gameState.PlayerHands[0], PlayerHand([]int{2})) && slices.Equal(gameState.PlayerHands[1], PlayerHand([]int{4, 5, 6})))

}

func Test_removePlayerInvalid(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{2}), PlayerHand([]int{5}), PlayerHand([]int{4, 5, 6})}
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	gameState.CurrentPlayerIndex = 0
	err := gameState.removePlayer(4)
	if err == nil {
		t.Error("Allowed the removal of a nonexistent player")
	}
	util.Assert(t, len(gameState.PlayerHands) == len(gameState.PlayerChannels) && len(gameState.PlayerHands) == 3)
}

func Test_removePlayerCurrentTurn(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{2}), PlayerHand([]int{5}), PlayerHand([]int{4, 5, 6})}
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 4))
	gameState.CurrentPlayerIndex = 0
	err := gameState.removePlayer(0)
	if err != nil {
		t.Fail()
	}
	util.Assert(t, gameState.CurrentPlayerIndex == 1)
	// util.Assert(t, len(gameState.PlayerChannels) == 2)
	util.Assert(t, len(gameState.PlayerHands) == 2)
	util.Assert(t, slices.Equal(gameState.PlayerHands[0], PlayerHand([]int{5})) && slices.Equal(gameState.PlayerHands[1], PlayerHand([]int{4, 5, 6})))
}
