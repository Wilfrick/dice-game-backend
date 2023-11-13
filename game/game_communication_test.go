package game

import (
	"HigherLevelPerudoServer/util"
	"bytes"
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
	msg := Message{TypeDescriptor: "Bananas"}
	gameState.send(PLAYER_INDEX, msg)
	recieve := <-gameState.PlayerChannels[0]
	encodedMsg := createEncodedMessage(msg)
	util.Assert(t, bytes.Equal(recieve, encodedMsg))
}
func Test_sendWithWaitGroupSimple(t *testing.T) {
	var gameState GameState
	gameState.PlayerChannels = make([]chan []byte, 1)
	gameState.PlayerChannels[0] = make(chan []byte)
	PLAYER_INDEX := 0
	var wait_group sync.WaitGroup
	msg := Message{TypeDescriptor: "Bananas"}
	gameState.send(PLAYER_INDEX, msg, &wait_group)
	recieve := <-gameState.PlayerChannels[0]
	wait_group.Wait()
	encodedMsg := createEncodedMessage(msg)
	util.Assert(t, bytes.Equal(recieve, encodedMsg))
}

func Test_sendWithWaitGroupMultiParty(t *testing.T) {
	var gameState GameState
	gameState.PlayerChannels = util.InitialiseChans(make([]chan []byte, 3))
	var wait_group sync.WaitGroup
	msg := Message{TypeDescriptor: "Bananas"}
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
	encodedMsg := createEncodedMessage(msg)
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

	true_result1 := createEncodedMessage(Message{TypeDescriptor: "SinglePlayerHandContents", Contents: SinglePlayerHandContents{gameState.PlayerHands[0], 0}})
	true_result2 := createEncodedMessage(Message{TypeDescriptor: "SinglePlayerHandContents", Contents: SinglePlayerHandContents{gameState.PlayerHands[1], 1}})
	true_result3 := createEncodedMessage(Message{TypeDescriptor: "SinglePlayerHandContents", Contents: SinglePlayerHandContents{gameState.PlayerHands[2], 2}})

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
	true_result := createEncodedMessage(Message{TypeDescriptor: "SinglePlayerHandContents", Contents: SinglePlayerHandContents{gameState.PlayerHands[0], 0}})
	util.Assert(t, bytes.Equal(result, true_result))
}

func Test_broadcastSimple(t *testing.T) {
	gs := GameState{PlayerChannels: util.InitialiseChans(make([]chan []byte, 1))}
	msg := Message{TypeDescriptor: "Bananas"}
	go gs.broadcast(msg)
	util.Assert(t, bytes.Equal(<-gs.PlayerChannels[0], createEncodedMessage(msg)))
}

func Test_broadcastTwoPlayers(t *testing.T) {
	gs := GameState{PlayerChannels: util.InitialiseChans(make([]chan []byte, 2))}
	msg := Message{TypeDescriptor: "Bananas"}
	go gs.broadcast(msg)
	util.Assert(t, bytes.Equal(<-gs.PlayerChannels[0], createEncodedMessage(msg)))
	util.Assert(t, bytes.Equal(<-gs.PlayerChannels[1], createEncodedMessage(msg)))
}
func Test_broadcastWithWaitgroupSimple(t *testing.T) {
	gs := GameState{PlayerChannels: util.InitialiseChans(make([]chan []byte, 1))}
	msg := Message{TypeDescriptor: "Bananas"}
	use_waitgroup := true
	go gs.broadcast(msg, use_waitgroup)
	util.Assert(t, bytes.Equal(<-gs.PlayerChannels[0], createEncodedMessage(msg)))
}

func xTest_broadcastWithWaitgroupTwoPlayers(t *testing.T) {
	gs := GameState{PlayerChannels: util.InitialiseChans(make([]chan []byte, 2))}
	msg := Message{TypeDescriptor: "Bananas"}
	msg2 := Message{TypeDescriptor: "Oranges"}
	counter := 0
	go func(gs GameState, counter *int) {
		*counter += 1
		use_waitgroup := false
		gs.broadcast(msg, use_waitgroup)
		gs.broadcast(msg2, use_waitgroup)
	}(gs, &counter)
	<-gs.PlayerChannels[0]

	go func(gs GameState, t *testing.T) {
		<-gs.PlayerChannels[0]
		res := <-gs.PlayerChannels[1]
		if bytes.Equal(res, createEncodedMessage(msg)) {
			t.Error("could extract values without waiting for all parties")
		}
		// fail?
	}(gs, t)

	go func(gs GameState) {
		<-gs.PlayerChannels[1]
		<-gs.PlayerChannels[0]
		// succeed
	}(gs)

}

func xTest_consistent(t *testing.T) {
	c := make(chan int)
	//
	l := []int{0}
	go func() { <-c }()
	go func() { c <- 0 }()
	go func() { _ = l[<-c] }()
	go func() { c <- 1 }()
	_ = l[0]
}
