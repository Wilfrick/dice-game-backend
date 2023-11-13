package game

import (
	"HigherLevelPerudoServer/util"
	"testing"
)

func Test_alivePlayerIndices(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{1, 2, 3}), PlayerHand([]int{5, 4, 1}), PlayerHand([]int{4, 5, 6})}
	alivePlayerIndices := gameState.alivePlayerIndices()
	t.Log(alivePlayerIndices)
	util.Assert(t, len(alivePlayerIndices) == 3)
	truth := []int{0, 1, 2}
	for i, player := range alivePlayerIndices {
		util.Assert(t, player == truth[i])
	}

}

func Test_previousAlivePlayer(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{1, 2, 3}), PlayerHand([]int{5, 4, 1}), PlayerHand([]int{4, 5, 6})}
	gameState.CurrentPlayerIndex = 1
	previousAlivePlayer, err := gameState.PreviousAlivePlayer()
	if !(err == nil) {
		t.Fail()
	}
	t.Log(previousAlivePlayer)
	util.Assert(t, previousAlivePlayer == 0)

}

func Test_previousAlivePlayerCurrentPlayer0(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{1, 2, 3}), PlayerHand([]int{5, 4, 1}), PlayerHand([]int{4, 5, 6})}
	gameState.CurrentPlayerIndex = 0
	previousAlivePlayer, err := gameState.PreviousAlivePlayer()
	if !(err == nil) {
		t.Fail()
	}
	t.Log(previousAlivePlayer)
	util.Assert(t, previousAlivePlayer == 2)

}

func Test_previousAlivePlayerOneDead(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{1, 2, 3}), PlayerHand([]int{}), PlayerHand([]int{4, 5, 6})}
	gameState.CurrentPlayerIndex = 0
	previousAlivePlayer, err := gameState.PreviousAlivePlayer()
	if !(err == nil) {
		t.Fail()
	}
	util.Assert(t, previousAlivePlayer == 2)

}

func Test_previousAlivePlayerAllDead(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{}), PlayerHand([]int{}), PlayerHand([]int{})}
	gameState.CurrentPlayerIndex = 0
	_, err := gameState.PreviousAlivePlayer()
	if err == nil {
		t.FailNow()
	}
	util.Assert(t, err.Error() == "not enough alive players")

}

func Test_previousAlivePlayerOnePlayerAlive(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{5, 6, 2}), PlayerHand([]int{}), PlayerHand([]int{})}
	gameState.CurrentPlayerIndex = 0
	_, err := gameState.PreviousAlivePlayer()
	if err == nil {
		t.FailNow()
	}
	util.Assert(t, err.Error() == "not enough alive players")

}
