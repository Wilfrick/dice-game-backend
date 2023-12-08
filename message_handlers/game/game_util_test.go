package game

import (
	"HigherLevelPerudoServer/messages"
	"HigherLevelPerudoServer/util"
	"slices"
	"testing"
)

func Test_RemoveDice(t *testing.T) {
	var gameState GameState
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
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
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
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
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{1, 2, 3}), PlayerHand([]int{5}), PlayerHand([]int{4, 5, 6})}
	gameState.CurrentPlayerIndex = 0
	err := gameState.FindNextAlivePlayerInclusive()
	if err != nil {
		t.Fail()
	}
	t.Log(gameState.CurrentPlayerIndex)
	util.Assert(t, gameState.CurrentPlayerIndex == 0)
}
func Test_nextPlayerAliveDeadPlayer(t *testing.T) {
	var gameState GameState
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{}), PlayerHand([]int{5}), PlayerHand([]int{4, 5, 6})}
	gameState.InitialiseSlicesWithDefaults()
	gameState.CurrentPlayerIndex = 0
	err := gameState.FindNextAlivePlayerInclusive()
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
	err := gameState.RemovePlayer(1)
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
	err := gameState.RemovePlayer(4)
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
	err := gameState.RemovePlayer(0)
	if err != nil {
		t.Fail()
	}
	util.Assert(t, gameState.CurrentPlayerIndex == 1)
	// util.Assert(t, len(gameState.PlayerChannels) == 2)
	util.Assert(t, len(gameState.PlayerHands) == 2)
	util.Assert(t, slices.Equal(gameState.PlayerHands[0], PlayerHand([]int{5})) && slices.Equal(gameState.PlayerHands[1], PlayerHand([]int{4, 5, 6})))
}

// This test can't be well administered at the game level
// see game_util for discussion

func Test_removePlayerCausingVictory(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{2}), PlayerHand([]int{5})}
	thisChan := make(chan []byte)
	otherChan := make(chan []byte)
	util.ChanSink([]chan []byte{thisChan})
	gameState.PlayerChannels = []chan []byte{thisChan, otherChan}
	gameState.CurrentPlayerIndex = 1
	gameState.GameInProgress = true
	go gameState.RemovePlayer(0)
	winningResult := <-otherChan
	expectedMessage := messages.CreateEncodedMessage(messages.Message{TypeDescriptor: "GameResult", Contents: GameResult{1, "win"}})
	util.Assert(t, gameState.GameInProgress == false)
	util.Assert(t, slices.Equal(winningResult, expectedMessage))
}

func Test_cleanUpInactivePlayers(t *testing.T) {
	var gameState GameState
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 5))
	gameState.InitialiseSlicesWithDefaults()
	t.Log(gameState.PlayerChannels)
	gameState.PlayerChannels[3] = nil

	gameState.CleanUpInactivePlayers()
	t.Log(len(gameState.PlayerChannels))
	util.Assert(t, len(gameState.PlayerChannels) == 4)
}

func Test_skipInactivePlayerInBettingRound(t *testing.T) { // Jim
	var gs GameState
	gs.PlayerHands = []PlayerHand{PlayerHand{2, 2}, PlayerHand{3, 3}, PlayerHand{5, 5}}
	gs.InitialiseSlicesWithDefaults()
	gs.PlayerChannels[1] = nil
	// gs.RoundMoveHistory = []PlayerMove{{MoveType: BET, Value: Bet{FaceVal: 4, NumDice: 1}, PlayerIndex: 0}}
	betToBeMade := PlayerMove{MoveType: BET, Value: Bet{FaceVal: 4, NumDice: 1}, PlayerIndex: 0}
	gs.CurrentPlayerIndex = 0

	gs.ProcessPlayerMove(betToBeMade)

	t.Log(gs.CurrentPlayerIndex, gs.RoundMoveHistory)
	util.Assert(t, gs.CurrentPlayerIndex == 2)

	util.Assert(t, slices.Equal(gs.RoundMoveHistory, []PlayerMove{betToBeMade}))
	// t.FailNow()
}

func Test_previousInactivePlayerSkippedOnDudo(t *testing.T) { // Alex
	var gs GameState
	gs.PlayerHands = []PlayerHand{PlayerHand{2, 2}, PlayerHand{3, 3}, PlayerHand{5, 5}}
	gs.InitialiseSlicesWithDefaults()
	gs.PlayerChannels[1] = nil
	gs.RoundMoveHistory = []PlayerMove{{MoveType: BET, Value: Bet{FaceVal: 4, NumDice: 1}, PlayerIndex: 0}}
	gs.CurrentPlayerIndex = 2

	gs.ProcessPlayerMove(PlayerMove{MoveType: DUDO, PlayerIndex: 2})

	util.Assert(t, gs.CurrentPlayerIndex == 0)
}

func Test_previousInactivePlayerSkippedOnCalza(t *testing.T) { // Jim
	var gs GameState
	gs.PlayerHands = []PlayerHand{PlayerHand{2, 2}, PlayerHand{3, 3}, PlayerHand{5, 5}}
	gs.InitialiseSlicesWithDefaults()
	gs.PlayerChannels[1] = nil
	gs.RoundMoveHistory = []PlayerMove{{MoveType: BET, Value: Bet{FaceVal: 4, NumDice: 1}, PlayerIndex: 0}}
	gs.CurrentPlayerIndex = 2

	gs.ProcessPlayerMove(PlayerMove{MoveType: CALZA, PlayerIndex: 2})

	util.Assert(t, gs.CurrentPlayerIndex == 2)
}

func Test_StartNewRoundDeletingInactivePlayers(t *testing.T) { // Jim
	var gs GameState
	gs.PlayerHands = []PlayerHand{PlayerHand{2, 2}, PlayerHand{3, 3}, PlayerHand{5, 5}}
	gs.InitialiseSlicesWithDefaults()
	gs.PlayerChannels[1] = nil

	gs.startNewRound()

	util.Assert(t, len(gs.PlayerChannels) == 2)
	util.Assert(t, len(gs.PlayerHands) == 2)
}

func Test_CleanInactivePlayersCurrentPlayerIndexChangedByInactivePlayers(t *testing.T) {
	var gs GameState
	gs.PlayerHands = []PlayerHand{{2}, {3}, {5}}
	gs.InitialiseSlicesWithDefaults()

	gs.CurrentPlayerIndex = 2
	gs.PlayerChannels[1] = nil

	gs.CleanUpInactivePlayers()

	t.Log(gs.CurrentPlayerIndex)
	util.Assert(t, len(gs.PlayerChannels) == 2)
	util.Assert(t, gs.CurrentPlayerIndex == 1)
}

func Test_startNewRoundCurrentPlayerIndexChangedByInactivePlayers(t *testing.T) {
	var gs GameState
	gs.PlayerHands = []PlayerHand{{2}, {3}, {5}}
	gs.InitialiseSlicesWithDefaults()

	gs.CurrentPlayerIndex = 2
	gs.PlayerChannels[1] = nil

	gs.startNewRound()

	t.Log(gs.CurrentPlayerIndex)
	util.Assert(t, len(gs.PlayerChannels) == 2)
	util.Assert(t, gs.CurrentPlayerIndex == 1)
}
