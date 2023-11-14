package game

import (
	"HigherLevelPerudoServer/util"
	"testing"
)

// func (gameState *GameState) updatePlayerIndexFinalBet(player_dice_change_index, other_involved_player int, bet_true bool) {

func Test_updatePlayerIndexFinalBetDudoFalseNoDeath(t *testing.T) {
	gameState := GameState{PlayerHands: []PlayerHand{PlayerHand{1, 2}, PlayerHand{2, 3}, PlayerHand{3}},
		PrevMove: PlayerMove{MoveType: BET, Value: Bet{4, 2}}}
	// player_index 2 called Dudo
	// the bet was false
	// therefore player_1 loses a dice
	gameState.PlayerHands[1] = PlayerHand{2}

	gameState.updatePlayerIndexFinalBet(1, 2)

	// expect CPI = 1
	util.Assert(t, gameState.CurrentPlayerIndex == 1)
}
func Test_updatePlayerIndexFinalBetDudoFalseDeath(t *testing.T) {
	gameState := GameState{PlayerHands: []PlayerHand{PlayerHand{1, 2}, PlayerHand{2, 3}, PlayerHand{3}},
		PrevMove: PlayerMove{MoveType: BET, Value: Bet{4, 2}}}
	// player_index 0 called Dudo
	// the bet was false
	// therefore player_index 2 loses a dice
	gameState.PlayerHands[2] = PlayerHand{}
	// therefore player_index 2 dies

	gameState.updatePlayerIndexFinalBet(2, 0)

	// CPI = 0
	util.Assert(t, gameState.CurrentPlayerIndex == 0)
}

func Test_updatePlayerIndexFinalBetDudoTrueNoDeath(t *testing.T) {
	gameState := GameState{PlayerHands: []PlayerHand{PlayerHand{1, 2}, PlayerHand{2, 3}, PlayerHand{3}},
		PrevMove: PlayerMove{MoveType: BET, Value: Bet{3, 2}}}
	// player_index 2 called Dudo
	// the bet was true
	// therefore player_2 loses a dice
	gameState.PlayerHands[2] = PlayerHand{1}

	gameState.updatePlayerIndexFinalBet(2, 1)

	// expect CPI = 2
	util.Assert(t, gameState.CurrentPlayerIndex == 2)
}
func Test_updatePlayerIndexFinalBetDudoTrueDeath(t *testing.T) {
	gameState := GameState{PlayerHands: []PlayerHand{PlayerHand{3}, PlayerHand{2, 3}, PlayerHand{1, 2}},
		PrevMove: PlayerMove{MoveType: BET, Value: Bet{3, 2}}}
	// player_index 0 called Dudo
	// the bet was true
	// therefore player_index 0 loses a dice
	gameState.PlayerHands[0] = PlayerHand{}
	// therefore player_index 0 dies

	gameState.updatePlayerIndexFinalBet(0, 2)

	// CPI = 2
	util.Assert(t, gameState.CurrentPlayerIndex == 2)
}

func Test_updatePlayerIndexFinalBetCalzaTrueNoDeath(t *testing.T) { // should be explored for testing processPlayerCalza
	gameState := GameState{PlayerHands: []PlayerHand{PlayerHand{1, 2}, PlayerHand{2, 3}, PlayerHand{3, 4, 4, 4, 4}},
		PrevMove: PlayerMove{MoveType: BET, Value: Bet{3, 2}}}
	// player_index 2 called Calza
	// the bet was true
	// therefore player_2 tries to gain a dice, but can't
	gameState.PlayerHands[2] = PlayerHand{3, 4, 4, 4, 4}

	gameState.updatePlayerIndexFinalBet(2, 1)

	// expect CPI = 2
	util.Assert(t, gameState.CurrentPlayerIndex == 2)
}

func Test_processPlayerDudoFalse(t *testing.T) {
	var gs GameState
	gs.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 2}), PlayerHand([]int{1, 1}), PlayerHand([]int{4, 4})}
	gs.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	util.ChanSink(gs.PlayerChannels)
	gs.PrevMove = PlayerMove{MoveType: BET, Value: Bet{5, 2}}
	gs.CurrentPlayerIndex = 1
	// playerMove := PlayerMove{MoveType: DUDO}     // Dudo False, so P1 loses a dice
	validity := gs.processPlayerDudo() // not a big fan of 'validity', would rather 'move could be made' or similar
	if !validity {
		t.Fail()
	}
	// gs.updatePlayerIndexFinalBet(1, 0)
	util.Assert(t, gs.CurrentPlayerIndex == 1)
	t.Log(gs.CurrentPlayerIndex)
	t.Log(gs.PlayerHands)
	// util.Assert(t, gs.CurrentPlayerIndex == 1)  // not passing
	util.Assert(t, len(gs.PlayerHands[1]) == 1) // passes
}

func Test_DudoIdentifyLosersWinnersDudoFalse(t *testing.T) {
	var gs GameState
	gs.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 2}), PlayerHand([]int{1, 1}), PlayerHand([]int{4, 4})}
	gs.PrevMove = PlayerMove{MoveType: BET, Value: Bet{5, 2}}
	gs.CurrentPlayerIndex = 1
	// DUDO FALSE
	loser, winner, _ := gs.DudoIdentifyLosersWinners()
	util.Assert(t, loser == 1)
	util.Assert(t, winner == 0)
}

func Test_DudoIdentifyLosersWinnersDudoTrue(t *testing.T) {
	var gs GameState
	gs.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 2}), PlayerHand([]int{1, 1}), PlayerHand([]int{4, 4})}
	gs.PrevMove = PlayerMove{MoveType: BET, Value: Bet{6, 2}}
	gs.CurrentPlayerIndex = 1
	// DUDO TRUE
	loser, winner, _ := gs.DudoIdentifyLosersWinners()
	util.Assert(t, loser == 0)
	util.Assert(t, winner == 1)
}

func Test_DudoIdentifyLosersWinnersWrappers(t *testing.T) {
	var gs GameState
	gs.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 2}), PlayerHand([]int{1, 1}), PlayerHand([]int{4, 4})}
	gs.PrevMove = PlayerMove{MoveType: BET, Value: Bet{6, 2}}
	gs.CurrentPlayerIndex = 0
	// DUDO TRUE
	loser, winner, _ := gs.DudoIdentifyLosersWinners()
	util.Assert(t, loser == 2)
	util.Assert(t, winner == 0)
}
