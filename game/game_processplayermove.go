package game

import (
	"fmt"
)

func (gameState *GameState) processPlayerBet(playerMove PlayerMove) bool {
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
	err := gameState.updatePlayerIndex(BET)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	gameState.broadcast(Message{"RoundResult", RoundResult{PlayerIndex: gameState.CurrentPlayerIndex, Result: "next"}})

	fmt.Println("Bet processed and broadcast RoundResult next")
	return true
}

func (gameState *GameState) processPlayerDudo() bool {
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

	var losing_player_index, candidate_victor int
	// var candidate_victor int
	previousAlivePlayer, err := gameState.PreviousAlivePlayer()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if bet_true {
		losing_player_index = gameState.CurrentPlayerIndex
		candidate_victor = previousAlivePlayer
	} else {
		losing_player_index = previousAlivePlayer
		candidate_victor = gameState.CurrentPlayerIndex
	}

	gameState.broadcast(Message{"RoundResult", RoundResult{losing_player_index, "dec"}})

	player_died, err := gameState.removeDice(losing_player_index)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if player_died { // 2) possible who has now lost ()
		gameState.broadcast(Message{"RoundResult", RoundResult{losing_player_index, "lose"}})
		gameState.send(losing_player_index, Message{"GameResult", GameResult{losing_player_index, "lose"}})
		if gameState.checkPlayerWin(candidate_victor) {
			gameState.broadcast(Message{"GameResult", GameResult{candidate_victor, "win"}})
			gameState.GameInProgress = false
			fmt.Println("A player has won and the game is no longer in progress")
			return true
		}
	}

	err = gameState.updatePlayerIndex(DUDO, losing_player_index)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	gameState.startNewRound()
	return true
}

func (gameState *GameState) processPlayerCalza() bool {
	fmt.Println("Made into Case Calza")
	//Input already valid
	gameState.broadcast(Message{"RoundUpdate", RoundUpdate{
		PlayerIndex: gameState.CurrentPlayerIndex,
		MoveMade:    PlayerMove{MoveType: "Calza"},
	}})
	bet_true := gameState.isBetExactlyTrue()
	// Not sure if the following code deserves a function
	var dice_total_changing_player_index int
	var candidate_victor int
	previousAlivePlayer, err := gameState.PreviousAlivePlayer()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if bet_true {
		dice_total_changing_player_index = previousAlivePlayer
		candidate_victor = gameState.CurrentPlayerIndex
	} else {
		dice_total_changing_player_index = gameState.CurrentPlayerIndex
		candidate_victor = previousAlivePlayer
	}
	gameState.broadcast(Message{"RoundResult", RoundResult{dice_total_changing_player_index, "dec"}})
	player_died, err := gameState.removeDice(dice_total_changing_player_index)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if player_died {
		gameState.broadcast(Message{"RoundResult", RoundResult{dice_total_changing_player_index, "lose"}})
		gameState.send(dice_total_changing_player_index, Message{"GameResult", GameResult{dice_total_changing_player_index, "lose"}})
		if gameState.checkPlayerWin(candidate_victor) {
			gameState.broadcast(Message{"GameResult", GameResult{candidate_victor, "win"}})
			gameState.GameInProgress = false
			fmt.Println("A player has won and the game is no longer in progress")
			return true
		}
	}
	err = gameState.updatePlayerIndex(CALZA, dice_total_changing_player_index)
	if err != nil {
		fmt.Println(err)
		return false
	}
	gameState.startNewRound()
	return true
}
