package game

import (
	"HigherLevelPerudoServer/util"
	"testing"
)

func Test_processPlayerMoveCalzaTrue(t *testing.T) {
	var gs GameState
	gs.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 2}), PlayerHand([]int{1, 1})}
	gs.PlayerChannels = util.InitialiseChans(make([]chan []byte, 2))
	go func() {
		for {
			<-gs.PlayerChannels[0]
			<-gs.PlayerChannels[1]
		}
	}()
	gs.PrevMove = PlayerMove{MoveType: "Bet", Value: Bet{5, 2}}
	gs.CurrentPlayerIndex = 1
	playerMove := PlayerMove{MoveType: "Calza"} // True
	validity := gs.ProcessPlayerMove(playerMove)
	if !validity {
		t.Logf("Validity: %t", validity)
		t.Fail()
	}
	t.Logf("CPI: %d", gs.CurrentPlayerIndex)
	util.Assert(t, gs.CurrentPlayerIndex == 1)  // Do we update player correctly when CALZA
	util.Assert(t, len(gs.PlayerHands[1]) == 3) // Do we update dice correctly when CALZA
}

func Test_processPlayerMoveCalzaFalse(t *testing.T) {
	var gs GameState
	gs.PlayerHands = []PlayerHand{PlayerHand([]int{2, 2, 3}), PlayerHand([]int{1, 1})}
	gs.PlayerChannels = util.InitialiseChans(make([]chan []byte, 2))
	go func() {
		for {
			<-gs.PlayerChannels[0]
			<-gs.PlayerChannels[1]
		}
	}()
	gs.PrevMove = PlayerMove{MoveType: "Bet", Value: Bet{5, 2}}
	gs.CurrentPlayerIndex = 1
	playerMove := PlayerMove{MoveType: "Calza"} // False only 4 2's
	validity := gs.ProcessPlayerMove(playerMove)
	if !validity {
		t.Logf("Validity %t", validity)
		t.FailNow()
	}
	t.Log(gs.CurrentPlayerIndex)
	util.Assert(t, gs.CurrentPlayerIndex == 1)
	util.Assert(t, len(gs.PlayerHands[1]) == 1)

}
