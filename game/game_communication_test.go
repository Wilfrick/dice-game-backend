package game

import (
	"HigherLevelPerudoServer/messages"
	"HigherLevelPerudoServer/util"
	"bytes"
	"fmt"
	"sync"
	"testing"
)

func Test_channelExample(t *testing.T) {
	chans := make([]chan []byte, 1)
	// chans[0] = make(chan []byte)
	util.InitialiseChans(chans)
	go func(cs []chan []byte) { cs[0] <- []byte("Hi") }(chans)
	t.Log("Did go")
	// t.FailNow()
	res := <-chans[0]
	util.Assert(t, bytes.Equal(res, []byte("Hi")))
}

func Test_sendWithoutWaitGroup(t *testing.T) {
	var gameState GameState
	gameState.PlayerChannels = make([]chan []byte, 1)
	gameState.PlayerChannels[0] = make(chan []byte)
	PLAYER_INDEX := 0
	msg := messages.Message{TypeDescriptor: "Bananas"}
	gameState.send(PLAYER_INDEX, msg)
	recieve := <-gameState.PlayerChannels[0]
	encodedMsg := messages.CreateEncodedMessage(msg)
	util.Assert(t, bytes.Equal(recieve, encodedMsg))
}
func Test_sendWithWaitGroupSimple(t *testing.T) {
	var gameState GameState
	gameState.PlayerChannels = make([]chan []byte, 1)
	gameState.PlayerChannels[0] = make(chan []byte)
	PLAYER_INDEX := 0
	var wait_group sync.WaitGroup
	msg := messages.Message{TypeDescriptor: "Bananas"}
	gameState.send(PLAYER_INDEX, msg, &wait_group)
	recieve := <-gameState.PlayerChannels[0]
	wait_group.Wait()
	encodedMsg := messages.CreateEncodedMessage(msg)
	util.Assert(t, bytes.Equal(recieve, encodedMsg))
}

func Test_sendWithWaitGroupMultiParty(t *testing.T) {
	var gameState GameState
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	var wait_group sync.WaitGroup
	msg := messages.Message{TypeDescriptor: "Bananas"}
	gameState.send(0, msg, &wait_group)
	gameState.send(1, msg, &wait_group)
	gameState.send(2, msg, &wait_group)
	recieve1 := <-gameState.PlayerChannels[0]
	<-gameState.PlayerChannels[1]

	done := make(chan int)
	go func(d chan int) {
		wait_group.Wait()
		d <- 1
	}(done)

	go func(d chan int) {
		<-gameState.PlayerChannels[2]
		// will Go ever switch go routines here? If not then this works as a test, but if so then this could fail for no apparent reason
		d <- 2
	}(done)
	if <-done == 1 {
		t.Fail()
	}
	encodedMsg := messages.CreateEncodedMessage(msg)
	util.Assert(t, bytes.Equal(recieve1, encodedMsg))
}

func Test_distributeHands(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{5, 6}), PlayerHand([]int{2, 3}), PlayerHand([]int{1})}
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	// t.Fail()
	// gameState.distributeHands()

	go gameState.distributeHands()

	result3 := <-gameState.PlayerChannels[2]
	result2 := <-gameState.PlayerChannels[1]
	result1 := <-gameState.PlayerChannels[0]

	// results should be the bytes of the strings for a SinglePlayerHandContents

	// assert that the results are as expected

	true_result1 := messages.CreateEncodedMessage(messages.Message{TypeDescriptor: "SinglePlayerHandContents", Contents: SinglePlayerHandContents{gameState.PlayerHands[0], 0}})
	true_result2 := messages.CreateEncodedMessage(messages.Message{TypeDescriptor: "SinglePlayerHandContents", Contents: SinglePlayerHandContents{gameState.PlayerHands[1], 1}})
	true_result3 := messages.CreateEncodedMessage(messages.Message{TypeDescriptor: "SinglePlayerHandContents", Contents: SinglePlayerHandContents{gameState.PlayerHands[2], 2}})

	util.Assert(t, bytes.Equal(result1, true_result1))
	util.Assert(t, bytes.Equal(result2, true_result2))
	util.Assert(t, bytes.Equal(result3, true_result3))

	// Does the connection representing player 0 recieved a message that is the correct information that should have recieved
	// Namely {5,6}, [2,2,1]

}

func Test_distributeSingularHand(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{5, 6})}
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 1))
	go gameState.distributeHands()

	result := <-gameState.PlayerChannels[0]
	true_result := messages.CreateEncodedMessage(messages.Message{TypeDescriptor: "SinglePlayerHandContents", Contents: SinglePlayerHandContents{gameState.PlayerHands[0], 0}})
	util.Assert(t, bytes.Equal(result, true_result))
}

func Test_broadcastSimple(t *testing.T) {
	gs := GameState{PlayerChannels: util.InitialiseChans(make([]chan []byte, 1))}
	msg := messages.Message{TypeDescriptor: "Bananas"}
	go gs.Broadcast(msg)
	util.Assert(t, bytes.Equal(<-gs.PlayerChannels[0], messages.CreateEncodedMessage(msg)))
}

func Test_broadcastTwoPlayers(t *testing.T) {
	gs := GameState{PlayerChannels: util.InitialiseChans(make([]chan []byte, 2))}
	msg := messages.Message{TypeDescriptor: "Bananas"}
	go gs.Broadcast(msg)
	util.Assert(t, bytes.Equal(<-gs.PlayerChannels[0], messages.CreateEncodedMessage(msg)))
	util.Assert(t, bytes.Equal(<-gs.PlayerChannels[1], messages.CreateEncodedMessage(msg)))
}
func Test_broadcastWithWaitgroupSimple(t *testing.T) {
	gs := GameState{PlayerChannels: util.InitialiseChans(make([]chan []byte, 1))}
	msg := messages.Message{TypeDescriptor: "Bananas"}
	use_waitgroup := true
	go gs.Broadcast(msg, use_waitgroup)
	util.Assert(t, bytes.Equal(<-gs.PlayerChannels[0], messages.CreateEncodedMessage(msg)))
}

// THE BELOW TEST IS UNTESTED
// MAY PROVIDE FURTHER TEST GUIDANCE
// AND EDUCATIONAL CONTENT / CONTEXT

// func Test_broadcastWithWaitgroupTwoPlayers(t *testing.T) {
// 	gs := GameState{PlayerChannels: util.InitialiseChans(make([]chan []byte, 2))}
// 	msg := messages.Message{TypeDescriptor: "Bananas"}
// 	msg2 := messages.Message{TypeDescriptor: "Oranges"}
// 	counter := 0
// 	go func(gs GameState, counter *int) {
// 		*counter += 1
// 		use_waitgroup := false
// 		gs.broadcast(msg, use_waitgroup)
// 		gs.broadcast(msg2, use_waitgroup)
// 	}(gs, &counter)
// 	<-gs.PlayerChannels[0]

// 	go func(gs GameState, t *testing.T) {
// 		<-gs.PlayerChannels[0]
// 		res := <-gs.PlayerChannels[1]
// 		if bytes.Equal(res,  messages.CreateEncodedMessage(msg)) {
// 			t.Error("could extract values without waiting for all parties")
// 		}
// 		// fail?
// 	}(gs, t)

// 	go func(gs GameState) {
// 		<-gs.PlayerChannels[1]
// 		<-gs.PlayerChannels[0]
// 		// succeed
// 	}(gs)

// }

func Test_revealHandsBasic(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{5, 6})}
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 1))
	go gameState.revealHands()

	result := <-gameState.PlayerChannels[0]
	true_result := messages.CreateEncodedMessage(messages.Message{TypeDescriptor: "PlayerHandsContents", Contents: PlayerHandsContents{PlayerHands: gameState.PlayerHands}})
	util.Assert(t, bytes.Equal(result, true_result))
}

func Test_revealHandsTwoPlayers(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{5, 6}), PlayerHand([]int{4, 4, 5})}
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 2))
	go gameState.revealHands()

	res1, res2 := <-gameState.PlayerChannels[0], <-gameState.PlayerChannels[1]

	true_result := messages.CreateEncodedMessage(messages.Message{TypeDescriptor: "PlayerHandsContents", Contents: PlayerHandsContents{PlayerHands: gameState.PlayerHands}})
	util.Assert(t, bytes.Equal(res1, true_result))
	util.Assert(t, bytes.Equal(res2, true_result))
}

func Test_revealHandsSomeDeadPlayers(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{}), PlayerHand([]int{4, 4, 5, 1, 1}), PlayerHand([]int{}), PlayerHand([]int{4, 4, 5})}
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 4))
	go gameState.revealHands()

	res1, res2, res3, res4 := <-gameState.PlayerChannels[0], <-gameState.PlayerChannels[1], <-gameState.PlayerChannels[2], <-gameState.PlayerChannels[3]

	true_result := messages.CreateEncodedMessage(messages.Message{TypeDescriptor: "PlayerHandsContents", Contents: PlayerHandsContents{PlayerHands: gameState.PlayerHands}})
	util.Assert(t, bytes.Equal(res1, true_result))
	util.Assert(t, bytes.Equal(res2, true_result))
	util.Assert(t, bytes.Equal(res3, true_result))
	util.Assert(t, bytes.Equal(res4, true_result))
}

func Test_revealHandsRaceCondition(t *testing.T) {
	var gameState GameState
	gameState.PlayerHands = []PlayerHand{PlayerHand([]int{4, 4, 5, 1, 1})}
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 1))
	msg := messages.Message{TypeDescriptor: "PlayerHandsContents", Contents: PlayerHandsContents{PlayerHands: []PlayerHand{PlayerHand([]int{4, 4, 5, 1, 1})}}}
	var output []byte

	gameState.revealHands() //Go routine here

	gameState.randomiseCurrentHands()
	gameState.randomiseCurrentHands()
	// wait_group.Wait()
	// fmt.Println("Waited")
	output = <-gameState.PlayerChannels[0]
	fmt.Println(string(output))
	t.Log(string(output))
	t.Log(string(messages.CreateEncodedMessage(msg)))
	util.Assert(t, bytes.Equal(output, messages.CreateEncodedMessage(msg)))

}
