package game

import (
	"HigherLevelPerudoServer/util"
	"errors"
	"fmt"
)

type GameState struct {
	GameID               string
	PlayerHands          []PlayerHand
	CurrentPlayerIndex   int // should be < len(PlayerHands)
	PlayerChannels       []chan []byte
	PrevMove             PlayerMove
	GameInProgress       bool
	AllowableChannelLock int // should live with PlayerChannels, wherever that ends up
}

type PlayerMove struct {
	MoveType string // "Bet", "Dudo", "Calza"
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

type GameResult struct {
	PlayerIndex int
	Result      string // "win", "lose"
}

type SinglePlayerHandContents struct {
	PlayerHand  PlayerHand
	PlayerIndex int
}

type PlayerHandLengthsUpdate struct {
	PlayerHandLengths []int
}

type PlayerHandContents struct {
	PlayerHands []PlayerHand
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

func (gameState *GameState) StartNewGame() {
	fmt.Println("Called StartNewGame")
	// Things to do: Generate player hands
	// Distribute player hands
	// Broadcast PlayerHandsLengthsUpdate
	// Choose a first player
	// Broadcast first player?
	gameState.PlayerHands = make([]PlayerHand, len(gameState.PlayerChannels))
	gameState.GameInProgress = true
	for i := range gameState.PlayerHands {
		gameState.PlayerHands[i] = RandomPlayerHand(5)
	}

	gameState.startNewRound()
	fmt.Println("Ended StartNewGame")
}

func (gameState *GameState) startNewRound() {
	fmt.Println("Called StartNewRound")
	gameState.randomiseCurrentHands()
	gameState.distributeHands()

	InitialPlayerHandLengths := PlayerHandLengthsUpdate{util.Map(func(x PlayerHand) int { return len(x) }, gameState.PlayerHands)}
	gameState.broadcast(Message{TypeDescriptor: "PlayerHandLengthsUpdate", Contents: InitialPlayerHandLengths})
	gameState.CurrentPlayerIndex = 0
	gameState.broadcast(Message{"RoundResult", RoundResult{gameState.CurrentPlayerIndex, "next"}})
}

func (gameState *GameState) ProcessPlayerMove(playerMove PlayerMove) bool {
	defer func() { gameState.AllowableChannelLock = gameState.CurrentPlayerIndex }()
	switch playerMove.MoveType {
	case "Bet":
		fmt.Println("in ProcessPlayerMove, made into case 'Bet' ")
		// fmt.Println(*gameState)
		// fmt.Println(playerMove)
		// validate bet
		newBet := playerMove.Value
		betValid := newBet.isGreaterThan(gameState.PrevMove.Value)
		betValid = true // ONLY For testing, not for production
		if !betValid {
			fmt.Println("Leaving processPlayerMove, case Bet, early")
			// should tell this player that their bet was invalid, then return
			return false // Representing invalid move / bet could not be made
		}
		// should also check that the player making the bet is the current player DEALT with higher up

		gameState.broadcast(Message{"RoundUpdate", RoundUpdate{MoveMade: playerMove, PlayerIndex: gameState.CurrentPlayerIndex}})

		gameState.PrevMove = playerMove

		// fmt.Println(playerMove)
		// fmt.Println(gameState.PrevMove)

		// Update current player
		_, err := gameState.updatePlayerIndex(playerMove)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}

		// thread one

		gameState.broadcast(Message{"RoundResult", RoundResult{PlayerIndex: gameState.CurrentPlayerIndex, Result: "next"}})

		fmt.Println("Bet processed and broadcast RoundResult next")
		return true
	case "Dudo":
		fmt.Println("in ProcessPlayerMove, made into case 'Dudo' ")
		// fmt.Println(*gameState)
		// fmt.Println(playerMove)
		// validate bet
		// dudo should always be valid, as long as the current player // checked earlier
		gameState.broadcast(Message{"RoundUpdate", RoundUpdate{
			PlayerIndex: gameState.CurrentPlayerIndex,
			MoveMade:    PlayerMove{MoveType: "Dudo"},
		}})

		// calculate the result of the call:
		// 1) who loses a dice (always happens)
		bet_true := gameState.isBetTrue()

		var losing_player_index int
		if bet_true {
			losing_player_index = gameState.CurrentPlayerIndex
		} else {
			var err error
			losing_player_index, err = gameState.PreviousAlivePlayer()
			if err != nil {
				fmt.Println(err.Error())
				return false
			}
		}

		gameState.broadcast(Message{"RoundResult", RoundResult{losing_player_index, "dec"}})

		player_died, err := gameState.removeDice(losing_player_index)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		if player_died {
			gameState.broadcast(Message{"RoundResult", RoundResult{losing_player_index, "lose"}})
			gameState.send(losing_player_index, Message{"GameResult", GameResult{losing_player_index, "lose"}})
		}
		// 2) possible who has now lost ()
		did_current_player_win, err := gameState.updatePlayerIndex(playerMove)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		if did_current_player_win {
			gameState.broadcast(Message{"GameResult", GameResult{gameState.CurrentPlayerIndex, "win"}})
			gameState.GameInProgress = false
			fmt.Println("A player has won and the game is no longer in progress")
			return true
			// gameState.send(gameState.CurrentPlayerIndex, Message{"GameResult", GameResult{gameState.CurrentPlayerIndex, "win"}})
		}
		gameState.startNewRound()
	}
	gameState.PrevMove = playerMove
	return true
}

func (gameState *GameState) updatePlayerIndex(newbet PlayerMove) (bool, error) {
	if len(gameState.PlayerHands) == 0 {
		err := errors.New("can't update a Game with no players")
		return false, err
	} else if gameState.PlayersAllDead() {
		return false, errors.New("all players are dead")
	}
	startingIndex := gameState.CurrentPlayerIndex
	gameState.CurrentPlayerIndex += 1
	gameState.CurrentPlayerIndex %= len(gameState.PlayerHands)
	newPlayerDead := len(gameState.PlayerHands[gameState.CurrentPlayerIndex]) == 0
	for newPlayerDead {
		gameState.CurrentPlayerIndex += 1
		gameState.CurrentPlayerIndex %= len(gameState.PlayerHands)
		if gameState.CurrentPlayerIndex == startingIndex {
			// err := errors.New("looped around to our initial player. all other players dead")
			// looped around to our initial player. all other players dead

			// This actually is the winning condition
			// Winner stored as gameState.CurrentPlayerIndex
			return true, nil
		}
		newPlayerDead = len(gameState.PlayerHands[gameState.CurrentPlayerIndex]) == 0
	}
	return false, nil
}

// broadcast message function, to used as: gameState.broadcast(message)
