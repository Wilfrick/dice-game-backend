package game

import (
	"errors"
	"fmt"
)

type GameState struct {
	GameID             string
	PlayerHands        []PlayerHand
	CurrentPlayerIndex int // should be < len(PlayerHands)
	PlayerChannels     []chan []byte
	PrevMove           PlayerMove
}

type PlayerMove struct {
	MoveType string // "Bet", "dudo", "calza"
	Value    Bet
}

type RoundResult struct { // the result of the round
	PlayerIndex int
	Result      string // "dec", "inc", "lose", "win", "next"
}

type RoundUpdate struct { // a player made a specific move
	MoveMade    PlayerMove
	PlayerIndex int
}

// 3 players left, A1, B1, C4
// A calls Dudo on C's bet and is wrong

// RoundUpdate{A called Dudo}
// RoundResult{A dec}
// RoundResult{A lose}
// RoundResult{C goes next}

// 2 players left, B1, C4
// C calls Dudo of B's bet and is correct

// RoundUpdate{C called Dudo}
// RoundResult{B dec}
// RoundResult{B lose}
// RoundResult{C win}

// GameResult{C win}

func (gameState *GameState) ProcessPlayerMove(playerMove PlayerMove) bool {
	switch playerMove.MoveType {
	case "Bet":
		fmt.Println("in ProcessPlayerMove, made into case 'Bet' ")
		fmt.Println(*gameState)
		fmt.Println(playerMove)
		// validate bet
		newBet := playerMove.Value
		betValid := newBet.isGreaterThan(gameState.PrevMove.Value)
		if !betValid {
			fmt.Println("Leaving processPlayerMove, case Bet, early")
			return false //Representing invalid move / bet could not be made
		}
		// should also check that the player making the bet is the current player

		gameState.broadcast(Message{"RoundUpdate", RoundUpdate{MoveMade: playerMove, PlayerIndex: gameState.CurrentPlayerIndex}})

		gameState.PrevMove = playerMove

		fmt.Println(playerMove)
		fmt.Println(gameState.PrevMove)

		// Update current player
		gameState.updatePlayerIndex(playerMove)

		gameState.broadcast(Message{"RoundResult", RoundResult{PlayerIndex: gameState.CurrentPlayerIndex, Result: "next"}})

		fmt.Println("Just broadcasted 'random values'")
		return true
	case "Dudo":

	}
	gameState.PrevMove = playerMove

	return true
}

// func getNewPlayerIndex () int {
func (gameState GameState) PlayersAllDead() bool {
	hand_sums := true
	for _, hand := range gameState.PlayerHands {
		if len(hand) > 0 {
			hand_sums = false
		}
	}
	return hand_sums
}

func (gameState *GameState) updatePlayerIndex(newbet PlayerMove) error {
	if len(gameState.PlayerHands) == 0 {
		err := errors.New("can't update a Game with no players")
		return err
	} else if gameState.PlayersAllDead() {
		return errors.New("all players are dead")
	}
	startingIndex := gameState.CurrentPlayerIndex
	gameState.CurrentPlayerIndex += 1
	gameState.CurrentPlayerIndex %= len(gameState.PlayerHands)
	newPlayerDead := len(gameState.PlayerHands[gameState.CurrentPlayerIndex]) == 0
	for newPlayerDead {
		gameState.CurrentPlayerIndex += 1
		gameState.CurrentPlayerIndex %= len(gameState.PlayerHands)
		if gameState.CurrentPlayerIndex == startingIndex {
			//Have done a loop no Bueno
			err := errors.New("looped around to our initial player. all other players dead")
			return err
		}
		newPlayerDead = len(gameState.PlayerHands[gameState.CurrentPlayerIndex]) == 0
	}
	return nil
}

// broadcast message function, to used as: gameState.broadcast(message)

func (gameState GameState) broadcast(message Message) {
	fmt.Println("Trying to broadcast message")
	encodedMessage := createEncodedMessage(message)
	for _, channel := range gameState.PlayerChannels {
		fmt.Println("Sending message")
		go func(c chan []byte) { c <- encodedMessage }(channel)
		// go func (){channel <- encodedMessage}() // likely to lead to untracable bugs, do not copy
	}
	fmt.Println("Finished broadcasting message")
}
