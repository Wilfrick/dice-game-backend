package game

import (
	"HigherLevelPerudoServer/util"
	"slices"
	"testing"
)

// func (gameState *GameState) updatePlayerIndexFinalBet(player_dice_change_index, other_involved_player int, bet_true bool) {

func Test_updatePlayerIndexFinalBetDudoFalseNoDeath(t *testing.T) {
	gameState := GameState{PlayerHands: []PlayerHand{PlayerHand{1, 2}, PlayerHand{2, 3}, PlayerHand{3}},
		RoundMoveHistory: []PlayerMove{{MoveType: BET, Value: Bet{4, 2}}}}
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
		RoundMoveHistory: []PlayerMove{{MoveType: BET, Value: Bet{4, 2}}}}
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
		RoundMoveHistory: []PlayerMove{{MoveType: BET, Value: Bet{3, 2}}}}
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
		RoundMoveHistory: []PlayerMove{{MoveType: BET, Value: Bet{3, 2}}}}
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
		RoundMoveHistory: []PlayerMove{{MoveType: BET, Value: Bet{3, 2}}}}
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
	gs.RoundMoveHistory = []PlayerMove{{MoveType: BET, Value: Bet{5, 2}}}
	gs.CurrentPlayerIndex = 1
	gs.PalacifoablePlayers = []bool{true, true, true}
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
	gs.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	gs.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 2}), PlayerHand([]int{1, 1}), PlayerHand([]int{4, 4})}
	gs.RoundMoveHistory = []PlayerMove{{MoveType: BET, Value: Bet{5, 2}}}
	gs.CurrentPlayerIndex = 1
	// DUDO FALSE
	loser, winner, _ := gs.DudoIdentifyLosersWinners()
	util.Assert(t, loser == 1)
	util.Assert(t, winner == 0)
}

func Test_DudoIdentifyLosersWinnersDudoTrue(t *testing.T) {
	var gs GameState
	gs.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	gs.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 2}), PlayerHand([]int{1, 1}), PlayerHand([]int{4, 4})}
	gs.RoundMoveHistory = []PlayerMove{{MoveType: BET, Value: Bet{6, 2}}}
	gs.CurrentPlayerIndex = 1
	// DUDO TRUE
	loser, winner, _ := gs.DudoIdentifyLosersWinners()
	util.Assert(t, loser == 0)
	util.Assert(t, winner == 1)
}

func Test_DudoIdentifyLosersWinnersWrappers(t *testing.T) {
	var gs GameState
	gs.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	gs.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 2}), PlayerHand([]int{1, 1}), PlayerHand([]int{4, 4})}
	gs.RoundMoveHistory = []PlayerMove{{MoveType: BET, Value: Bet{6, 2}}}
	gs.CurrentPlayerIndex = 0
	// DUDO TRUE
	loser, winner, _ := gs.DudoIdentifyLosersWinners()
	util.Assert(t, loser == 2)
	util.Assert(t, winner == 0)
}

func Test_CalzaShouldntBePossibleOnFirstMoveOfRoundNoPrevMove(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1, 1, 1}, {3, 3, 3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3))}
	util.ChanSink(gs.PlayerChannels)

	res := gs.processPlayerCalza()
	util.Assert(t, res == false)
}

func Test_CalzaShouldntBePossibleOnFirstMoveOfRoundDudoPrevMove(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1, 1, 1}, {3, 3, 3}, {4, 4, 4}}, RoundMoveHistory: []PlayerMove{{MoveType: DUDO}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3))}
	util.ChanSink(gs.PlayerChannels)

	res := gs.processPlayerCalza()
	util.Assert(t, res == false)
}

func Test_CalzaShouldntBePossibleOnFirstMoveOfRoundCalzaPrevMove(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1, 1, 1}, {3, 3, 3}, {4, 4, 4}}, RoundMoveHistory: []PlayerMove{{MoveType: CALZA}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3))}
	util.ChanSink(gs.PlayerChannels)

	res := gs.processPlayerCalza()
	util.Assert(t, res == false)
}

func Test_DudoShouldntBePossibleOnFirstMoveOfRoundNoPrevMove(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1, 1, 1}, {3, 3, 3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3))}
	util.ChanSink(gs.PlayerChannels)

	res := gs.processPlayerDudo()
	util.Assert(t, res == false)
}

func Test_DudoShouldntBePossibleOnFirstMoveOfRoundDudoPrevMove(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1, 1, 1}, {3, 3, 3}, {4, 4, 4}}, RoundMoveHistory: []PlayerMove{{MoveType: DUDO}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3))}
	util.ChanSink(gs.PlayerChannels)

	res := gs.processPlayerDudo()
	util.Assert(t, res == false)
}

func Test_DudoShouldntBePossibleOnFirstMoveOfRoundCalzaPrevMove(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1, 1, 1}, {3, 3, 3}, {4, 4, 4}}, RoundMoveHistory: []PlayerMove{{MoveType: CALZA}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3))}
	util.ChanSink(gs.PlayerChannels)

	res := gs.processPlayerDudo()
	util.Assert(t, res == false)
}

func Test_rejectOnesOnTheFirstMove(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1, 1, 1}, {3, 3, 3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3))}
	util.ChanSink(gs.PlayerChannels)
	bet := Bet{FaceVal: 1, NumDice: 1}
	res := gs.processPlayerBet(PlayerMove{MoveType: BET, Value: bet})
	util.Assert(t, res == false)
}
func Test_rejectOnesOnTheFirstMoveAfterADudo(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1, 1, 1}, {3, 3, 3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3))}
	gs.RoundMoveHistory = []PlayerMove{{MoveType: DUDO}}
	util.ChanSink(gs.PlayerChannels)
	bet := Bet{FaceVal: 1, NumDice: 1}
	res := gs.processPlayerBet(PlayerMove{MoveType: BET, Value: bet})
	util.Assert(t, res == false)
}
func Test_rejectOnesOnTheFirstMoveAfterACalza(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1, 1, 1}, {3, 3, 3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3))}
	gs.RoundMoveHistory = []PlayerMove{{MoveType: CALZA}}
	util.ChanSink(gs.PlayerChannels)
	bet := Bet{FaceVal: 1, NumDice: 1}
	res := gs.processPlayerBet(PlayerMove{MoveType: BET, Value: bet})
	util.Assert(t, res == false)
}

func Test_checkPalacifoBettingRuleOnePlayerFirstTime(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1}, {3, 3, 3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3)), PalacifoablePlayers: []bool{false, true, true}}
	gs.RoundMoveHistory = []PlayerMove{{MoveType: DUDO}}
	util.ChanSink(gs.PlayerChannels)
	gs.IsPalacifoRound = true
	bet := Bet{1, 1}
	res := gs.processPlayerBet(PlayerMove{MoveType: BET, Value: bet})
	util.Assert(t, res == true) // because this is a palacifo round
	// util.Assert(t, gs.PalacifoablePlayers[0] == false) // this player can't trigger another palacifo round
}

func Test_checkPalacifoBettingRulePalacifoFollowingMove(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1}, {3, 3, 3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3)), PalacifoablePlayers: []bool{false, true, true}}
	gs.IsPalacifoRound = true
	gs.RoundMoveHistory = []PlayerMove{{MoveType: DUDO}}
	gs.CurrentPlayerIndex = 0
	util.ChanSink(gs.PlayerChannels)
	bet := Bet{1, 1}
	res := gs.processPlayerBet(PlayerMove{MoveType: BET, Value: bet})
	// util.Assert(t, gs.IsPalacifoRound == true)
	util.Assert(t, res == true) // because this is a palacifo round
	// util.Assert(t, gs.PalacifoablePlayers[0] == false) // this player can't trigger another palacifo round
	bet = Bet{1, 2}
	res = gs.processPlayerBet(PlayerMove{MoveType: BET, Value: bet})
	util.Assert(t, res == false) // because this is a palacifo round and you must follow the rank
	// util.Assert(t, gs.PalacifoablePlayers[0] == false)
}

func Test_checkPalacifoBettingRulePalacifoFollowingMovePlayerWithOneDice(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1}, {3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3)), PalacifoablePlayers: []bool{false, false, true}}
	gs.IsPalacifoRound = true
	bet := Bet{1, 1}
	gs.RoundMoveHistory = []PlayerMove{{MoveType: BET, Value: bet}}
	gs.CurrentPlayerIndex = 1
	util.ChanSink(gs.PlayerChannels)

	bet = Bet{1, 6}

	res := gs.processPlayerBet(PlayerMove{MoveType: BET, Value: bet})
	util.Assert(t, res == true) // because this is a palacifo round but a player on one dice can change the suit
	// util.Assert(t, gs.PalacifoablePlayers[0] == false) // this player can't trigger another palacifo round
	bet = Bet{2, 2}
	res = gs.processPlayerBet(PlayerMove{MoveType: BET, Value: bet})
	util.Assert(t, res == false) // same as above
}
func Test_checkPalacifoOnesBiddableFirstTurn(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1}, {3, 3, 3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3)), PalacifoablePlayers: []bool{false, true, true}}
	gs.RoundMoveHistory = []PlayerMove{{MoveType: DUDO}}
	// t.Log(gs.PrevMove.Value.FaceVal)
	gs.IsPalacifoRound = true
	gs.CurrentPlayerIndex = 0
	util.ChanSink(gs.PlayerChannels)
	bet := Bet{1, 1}
	res := gs.processPlayerBet(PlayerMove{MoveType: BET, Value: bet})
	util.Assert(t, res == true) // because this is a palacifo round
	// util.Assert(t, gs.PalacifoablePlayers[0] == false) // this player can't trigger another palacifo round
}

func Test_onPalacifoOnesNotWildcard(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1}, {3, 3, 3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3)), PalacifoablePlayers: []bool{false, true, true}}
	gs.IsPalacifoRound = true
	bet_dependant_on_one := Bet{FaceVal: 3, NumDice: 4}
	gs.RoundMoveHistory = []PlayerMove{{MoveType: BET, Value: bet_dependant_on_one}}
	gs.CurrentPlayerIndex = 1
	util.ChanSink(gs.PlayerChannels)
	res := gs.ProcessPlayerMove(PlayerMove{MoveType: DUDO})
	util.Assert(t, res == true)
	util.Assert(t, len(gs.PlayerHands[0]) == 0)
}

func Test_onPalacifoCalzaNotAllowed(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1}, {3, 3, 3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3)), PalacifoablePlayers: []bool{false, true, true}}
	util.ChanSink(gs.PlayerChannels)
	gs.IsPalacifoRound = true
	bet_dependant_on_one := Bet{FaceVal: 3, NumDice: 3}
	gs.RoundMoveHistory = []PlayerMove{{MoveType: BET, Value: bet_dependant_on_one}}
	gs.CurrentPlayerIndex = 1
	res := gs.ProcessPlayerMove(PlayerMove{MoveType: CALZA})
	util.Assert(t, res == false)
	// util.Assert(t, len(gs.PlayerHands[1]) == 2)
}

func Test_updatePlayerIndexFinalBetDudoTrueNoDeathNotPalacifo(t *testing.T) { // not sure what this tests
	gameState := GameState{PlayerHands: []PlayerHand{PlayerHand{1, 2, 4}, PlayerHand{2, 3, 5}, PlayerHand{3, 5, 6}},
		RoundMoveHistory: []PlayerMove{{MoveType: BET, Value: Bet{3, 2}}}}
	gameState.processPlayerBet(PlayerMove{MoveType: BET, Value: Bet{2, 2}})

	util.Assert(t, gameState.IsPalacifoRound == false)
}

func Test_gameRoundStartSetsIsPalacifoRound(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1}, {3, 3, 3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3)), PalacifoablePlayers: []bool{true, true, true}}
	util.ChanSink(gs.PlayerChannels)
	gs.IsPalacifoRound = false
	gs.startNewRound()
	util.Assert(t, gs.IsPalacifoRound == true)
}

func Test_gameRoundStartUpdatesPalacifoablePlayers(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1}, {3, 3, 3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3)), PalacifoablePlayers: []bool{true, true, true}}
	util.ChanSink(gs.PlayerChannels)
	gs.IsPalacifoRound = false
	gs.startNewRound()
	util.Assert(t, gs.PalacifoablePlayers[0] == false)
}

func Test_PalacifoBettingOrder(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1}, {2}, {6}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3)), PalacifoablePlayers: []bool{false, false, false}, RoundMoveHistory: []PlayerMove{{MoveType: DUDO}}}
	util.ChanSink(gs.PlayerChannels)
	gs.IsPalacifoRound = true
	// On a Palaifo round, if a player has 1 dice, they can change faceval following
	// 1 < ..< 6
	// annoying cases below
	bet := Bet{1, 1}
	res := gs.processPlayerBet(PlayerMove{MoveType: BET, Value: bet})
	util.Assert(t, res == true)

	bet = Bet{1, 2}
	res = gs.processPlayerBet(PlayerMove{MoveType: BET, Value: bet})
	util.Assert(t, res == true)

	bet = Bet{1, 3}
	res = gs.processPlayerBet(PlayerMove{MoveType: BET, Value: bet})
	util.Assert(t, res == true)

	bet = Bet{1, 4}
	res = gs.processPlayerBet(PlayerMove{MoveType: BET, Value: bet})
	util.Assert(t, res == true)

	bet = Bet{1, 5}
	res = gs.processPlayerBet(PlayerMove{MoveType: BET, Value: bet})
	util.Assert(t, res == true)

	bet = Bet{1, 6}
	res = gs.processPlayerBet(PlayerMove{MoveType: BET, Value: bet})
	util.Assert(t, res == true)

}

func Test_gameRoundStartSetsIsPalacifoRoundOnlyFirstTime(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1}, {3, 3, 3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3)), PalacifoablePlayers: []bool{false, true, true}}
	util.ChanSink(gs.PlayerChannels)
	gs.IsPalacifoRound = false
	gs.startNewRound()
	util.Assert(t, gs.IsPalacifoRound == false)
}

func Test_startNewGameResetsPalaficoablePlayers(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1}, {3, 3, 3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3)), PalacifoablePlayers: []bool{false, false, true}}
	util.ChanSink(gs.PlayerChannels)
	gs.IsPalacifoRound = false
	gs.StartNewGame()

	util.Assert(t, gs.IsPalacifoRound == false)
	util.Assert(t, slices.Equal(gs.PalacifoablePlayers, []bool{true, true, true}))
}

func Test_byDefaultRoundAfterPalacifoIsNotAlsoPalacifo(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1}, {3, 3, 3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3)), PalacifoablePlayers: []bool{false, false, true}, RoundMoveHistory: []PlayerMove{{MoveType: BET, Value: Bet{FaceVal: 3, NumDice: 3}}}}
	gs.CurrentPlayerIndex = 2
	util.ChanSink(gs.PlayerChannels)
	gs.IsPalacifoRound = true

	gs.processPlayerDudo() // will end the round

	util.Assert(t, gs.IsPalacifoRound == false)
	util.Assert(t, slices.Equal(gs.PalacifoablePlayers, []bool{false, false, true}))
}
func Test_byDefaultRoundAfterPalacifoIsNotAlsoPalacifoUnlessFollowingPlayerGoesToTheirPalacifoRound(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1}, {3, 3}, {4, 4, 4}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 3)), PalacifoablePlayers: []bool{false, true, true}, RoundMoveHistory: []PlayerMove{{MoveType: BET, Value: Bet{FaceVal: 3, NumDice: 3}}}}
	gs.CurrentPlayerIndex = 2
	util.ChanSink(gs.PlayerChannels)
	gs.IsPalacifoRound = true

	gs.processPlayerDudo() // will end the round

	// util.Assert(t, gs.IsPalacifoRound == true)
	// util.Assert(t, slices.Equal(gs.PalacifoablePlayers, []bool{false, false, true}))
	util.Assert(t, len(gs.PlayerHands[1]) == 1)
}

func Test_PalacifoRoundDudoNotCountWildcards(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1}, {2}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 2)), PalacifoablePlayers: []bool{false, false}, RoundMoveHistory: []PlayerMove{{MoveType: BET, Value: Bet{FaceVal: 2, NumDice: 2}}}}
	gs.CurrentPlayerIndex = 0
	util.ChanSink(gs.PlayerChannels)
	gs.IsPalacifoRound = true

	gs.processPlayerDudo()
	t.Log(gs.PlayerHands)
	util.Assert(t, len(gs.PlayerHands[1]) == 0)
}

func Test_NormalRoundDudoDoCountWildcards(t *testing.T) {
	gs := GameState{PlayerHands: []PlayerHand{{1}, {2}}, PlayerChannels: util.InitialiseChans(make([]chan []byte, 2)), PalacifoablePlayers: []bool{false, false}, RoundMoveHistory: []PlayerMove{{MoveType: BET, Value: Bet{FaceVal: 2, NumDice: 2}}}}
	gs.CurrentPlayerIndex = 0
	util.ChanSink(gs.PlayerChannels)
	gs.IsPalacifoRound = false

	gs.processPlayerDudo()
	t.Log(gs.PlayerHands)
	util.Assert(t, len(gs.PlayerHands[1]) == 1)
}

func Test_validBetGetsAddedToRoundMoveHistory(t *testing.T) {
	gameState := GameState{PlayerHands: []PlayerHand{{3, 3}, {5, 5}}}
	gameState.InitialiseSlicesWithDefaults()

	moveToBeMade := PlayerMove{MoveType: BET, Value: Bet{NumDice: 3, FaceVal: 4}, PlayerIndex: 0}
	gameState.ProcessPlayerMove(moveToBeMade)
	t.Log(gameState.RoundMoveHistory)
	util.Assert(t, slices.Equal(gameState.RoundMoveHistory, []PlayerMove{moveToBeMade}))
}
