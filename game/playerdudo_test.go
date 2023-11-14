package game

import (
	"HigherLevelPerudoServer/util"
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

func xTest_processPlayerMoveDudoFalse(t *testing.T) {
	var gs GameState
	gs.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 2}), PlayerHand([]int{1, 1}), PlayerHand([]int{4, 4})}
	gs.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	util.ChanSink(gs.PlayerChannels)
	gs.PrevMove = PlayerMove{MoveType: BET, Value: Bet{5, 2}}
	gs.CurrentPlayerIndex = 1
	playerMove := PlayerMove{MoveType: DUDO}     // Dudo False, so P1 loses a dice
	validity := gs.ProcessPlayerMove(playerMove) // not a big fan of 'validity', would rather 'move could be made' or similar
	if !validity {
		t.Fail()
	}
	t.Log(gs.CurrentPlayerIndex)
	t.Log(gs.PlayerHands)
	util.Assert(t, gs.CurrentPlayerIndex == 1) // ✓
	util.Assert(t, len(gs.PlayerHands[gs.CurrentPlayerIndex]) == 1)

}

func Test_processPlayerMoveDudoTrue(t *testing.T) {
	var gs GameState
	gs.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 3}), PlayerHand([]int{1, 1})}
	gs.PlayerChannels = util.InitialiseChans(make([]chan []byte, 2))
	go func() {
		for { // could use select to sink all messages without any possibility of blocking
			<-gs.PlayerChannels[0]
			<-gs.PlayerChannels[1]
		}
	}()
	gs.PrevMove = PlayerMove{MoveType: "Bet", Value: Bet{5, 2}}
	gs.CurrentPlayerIndex = 1
	playerMove := PlayerMove{MoveType: "Dudo"} // True only 4 2's
	validity := gs.ProcessPlayerMove(playerMove)
	if !validity {
		t.Error(validity)
	}
	t.Log(gs.CurrentPlayerIndex)
	util.Assert(t, gs.CurrentPlayerIndex == 0) // ✓

}

func Test_processPlayerMoveDudoFalseKilling(t *testing.T) {
	var gs GameState
	gs.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 3}), PlayerHand([]int{1}), PlayerHand([]int{5})}
	gs.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	// Have a waiter on each channel
	// Needed to deal with sends
	go func() {
		for {
			<-gs.PlayerChannels[0]
		}
	}()
	go func() {
		for {
			<-gs.PlayerChannels[1]
		}
	}()
	go func() {
		for {
			<-gs.PlayerChannels[2]
		}
	}()
	gs.PrevMove = PlayerMove{MoveType: "Bet", Value: Bet{2, 2}}
	gs.CurrentPlayerIndex = 1
	playerMove := PlayerMove{MoveType: "Dudo"} // False 3 2's, so P1 loses and dies
	validity := gs.ProcessPlayerMove(playerMove)
	if !validity {
		t.Error(validity)
	}
	t.Log(gs.CurrentPlayerIndex)
	t.Log(gs.PlayerHands)
	util.Assert(t, gs.CurrentPlayerIndex == 0) //If a player dies, the next player is the other player involved in the call
}
