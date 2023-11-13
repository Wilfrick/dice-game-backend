package game

import (
	"fmt"
	"sync"
)

func (gameState GameState) send(player_index int, msg Message, wait_groups ...*sync.WaitGroup) {
	if len(wait_groups) == 1 {
		wait_groups[0].Add(1)
	}

	fmt.Println("Called send")
	go func() {
		if len(wait_groups) == 1 {
			defer wait_groups[0].Done()
		}

		gameState.PlayerChannels[player_index] <- createEncodedMessage(msg)
	}()
}

func (gameState GameState) distributeHands() {
	fmt.Println("Called distributeHands")
	var distribute_hands_wait_group sync.WaitGroup
	for playerHandIndex, playerHand := range gameState.PlayerHands {
		gameState.send(playerHandIndex,
			Message{"SinglePlayerHandContents",
				SinglePlayerHandContents{PlayerHand: playerHand,
					PlayerIndex: playerHandIndex}},
			&distribute_hands_wait_group)
	}
	// Wait for all the hands to be sent
	fmt.Println("Waiting for completion")
	distribute_hands_wait_group.Wait()
	fmt.Println("Completed")
}

func (gameState GameState) broadcast(message Message, optional_use_wait_group ...bool) {
	fmt.Println("Trying to broadcast message")
	var wait_group sync.WaitGroup

	use_wait_group := len(optional_use_wait_group) == 1 && optional_use_wait_group[0]
	// encodedMessage := createEncodedMessage(message)
	// for _, channel := range gameState.PlayerChannels {
	// 	fmt.Println("Sending message")
	// 	go func(c chan []byte) { c <- encodedMessage }(channel)
	// 	// go func (){channel <- encodedMessage}() // likely to lead to untracable bugs, do not copy
	// }

	for player_index := range gameState.PlayerChannels {
		fmt.Println("Sending message")
		if use_wait_group {
			gameState.send(player_index, message, &wait_group)
		} else {
			gameState.send(player_index, message)
		}

		// go func (){channel <- encodedMessage}() // likely to lead to untracable bugs, do not copy
	}
	if use_wait_group {
		wait_group.Wait()
	}

	fmt.Println("Finished broadcasting message")
}
