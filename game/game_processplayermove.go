package game

import (
	"fmt"
)

func (gameState GameState) checkNewBetValid(newBet Bet) bool {
	betValid := newBet.isGreaterThan(gameState.PrevMove.Value)
	betValid = true // ONLY For testing, not for production
	return betValid
}

func (gameState GameState) broadcastPlayerMove(playerMove PlayerMove) {
	gameState.broadcast(Message{"RoundUpdate", RoundUpdate{MoveMade: playerMove, PlayerIndex: gameState.CurrentPlayerIndex}})
}

func (gameState GameState) broadcastNextPlayer() {
	gameState.broadcast(Message{"RoundResult", RoundResult{PlayerIndex: gameState.CurrentPlayerIndex, Result: "next"}})
}

func (gameState GameState) broadcastDiceDec(playerIndex int) {
	gameState.broadcast(Message{"RoundResult", RoundResult{PlayerIndex: playerIndex, Result: DEC}})
}

func (gameState GameState) broadcastDiceInc(playerIndex int) {
	gameState.broadcast(Message{"RoundResult", RoundResult{PlayerIndex: playerIndex, Result: INC}})
}

func (gameState *GameState) processPlayerBet(playerMove PlayerMove) bool {
	fmt.Println("in ProcessPlayerMove, made into case 'Bet' ")
	betValid := gameState.checkNewBetValid(playerMove.Value)
	if !betValid {
		return false
	}
	gameState.broadcastPlayerMove(playerMove)

	gameState.PrevMove = playerMove

	// Update current player
	err := gameState.updatePlayerIndex(BET)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	gameState.broadcastNextPlayer()
	fmt.Println("Bet processed and broadcast RoundResult next")
	return true
}

func (gameState GameState) DudoIdentifyLosersWinners() (int, int, error) {
	var losing_player_index, candidate_victor int
	// var candidate_victor int
	previousAlivePlayer, err := gameState.PreviousAlivePlayer()
	bet_true := gameState.isBetTrue()
	if bet_true {
		losing_player_index = gameState.CurrentPlayerIndex
		candidate_victor = previousAlivePlayer
	} else {
		losing_player_index = previousAlivePlayer
		candidate_victor = gameState.CurrentPlayerIndex
	}
	return losing_player_index, candidate_victor, err
}

// ProcessPlayerDeath returns outcome of whether candidate_victor has won
func (gameState GameState) processPlayerDeath(losing_player_index, candidate_victor int) bool {
	fmt.Println("Called processPlayerDeath")
	gameState.broadcast(Message{"RoundResult", RoundResult{losing_player_index, "lose"}})
	gameState.send(losing_player_index, Message{"GameResult", GameResult{losing_player_index, "lose"}})
	if gameState.checkPlayerWin(candidate_victor) {
		gameState.broadcast(Message{"GameResult", GameResult{candidate_victor, "win"}})
		gameState.GameInProgress = false
		fmt.Println("A player has won and the game is no longer in progress")
		return true
	}
	return false
}

func (gameState *GameState) updatePlayerIndexFinalBet(dice_change_player_index, other_player_involved_in_call_index int) {
	if len(gameState.PlayerHands[dice_change_player_index]) > 0 {
		gameState.CurrentPlayerIndex = dice_change_player_index
	} else {
		gameState.CurrentPlayerIndex = other_player_involved_in_call_index
	}
}

func (gameState *GameState) processPlayerDudo() bool {
	fmt.Println("in ProcessPlayerMove, made into case 'Dudo' ")

	// dudo should always be valid, as long as the current player // checked earlier
	gameState.broadcastPlayerMove(PlayerMove{MoveType: DUDO})
	gameState.revealHands()
	// calculate the result of the call:
	// 1) who loses a dice (always happens)
	losing_player_index, candidate_victor, err := gameState.DudoIdentifyLosersWinners()
	fmt.Println(losing_player_index, candidate_victor)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	gameState.broadcastDiceDec(losing_player_index)

	player_died, err := gameState.removeDice(losing_player_index)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if player_died { // 2) possible who has now lost ()
		if gameState.processPlayerDeath(losing_player_index, candidate_victor) {
			return true
		}
	}
	gameState.updatePlayerIndexFinalBet(losing_player_index, candidate_victor)

	// err = gameState.updatePlayerIndex(DUDO, losing_player_index) // probably needs to take candidate_victor as well
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return false
	// }
	gameState.startNewRound()
	return true
}

func (gameState *GameState) processPlayerCalza() bool {
	fmt.Println("Made into Case Calza")
	//Input already valid
	gameState.broadcastPlayerMove(PlayerMove{MoveType: CALZA})
	gameState.revealHands()
	bet_true := gameState.isBetExactlyTrue()
	// Not sure if the following code deserves a function
	if bet_true {
		// try to increment this player's hand
		gameState.addDice(gameState.CurrentPlayerIndex)
		gameState.broadcastDiceInc(gameState.CurrentPlayerIndex)

		// gameState.CurrentPlayerIndex = gameState.CurrentPlayerIndex

		// don't need to run because the current player hasn't changed

		gameState.startNewRound()
		return true
	}
	// bet was not true, so the player losing a dice is

	dice_total_changing_player_index := gameState.CurrentPlayerIndex
	other_player_index, err := gameState.PreviousAlivePlayer()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	// broadcast round result (dec, next)
	player_died, err := gameState.removeDice(dice_total_changing_player_index)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if player_died { // 2) possible who has now lost ()
		if gameState.processPlayerDeath(dice_total_changing_player_index, other_player_index) {
			return true
		}
	}
	gameState.updatePlayerIndexFinalBet(dice_total_changing_player_index, other_player_index)

	// err = gameState.updatePlayerIndex(DUDO, losing_player_index) // probably needs to take candidate_victor as well
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return false
	// }
	gameState.startNewRound()
	return true

	// var dice_total_changing_player_index int
	// var candidate_victor int
	// previousAlivePlayer, err := gameState.PreviousAlivePlayer()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return false
	// }
	// if bet_true {
	// 	dice_total_changing_player_index = previousAlivePlayer
	// 	candidate_victor = gameState.CurrentPlayerIndex
	// } else {
	// 	dice_total_changing_player_index = gameState.CurrentPlayerIndex
	// 	candidate_victor = previousAlivePlayer
	// }
	// gameState.broadcast(Message{"RoundResult", RoundResult{dice_total_changing_player_index, "dec"}})
	// player_died, err := gameState.removeDice(dice_total_changing_player_index)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return false
	// }
	// if player_died {
	// 	gameState.broadcast(Message{"RoundResult", RoundResult{dice_total_changing_player_index, "lose"}})
	// 	gameState.send(dice_total_changing_player_index, Message{"GameResult", GameResult{dice_total_changing_player_index, "lose"}})
	// 	if gameState.checkPlayerWin(candidate_victor) {
	// 		gameState.broadcast(Message{"GameResult", GameResult{candidate_victor, "win"}})
	// 		gameState.GameInProgress = false
	// 		fmt.Println("A player has won and the game is no longer in progress")
	// 		return true
	// 	}
	// }
	// err = gameState.updatePlayerIndex(CALZA, dice_total_changing_player_index)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return false
	// }
	// gameState.startNewRound()
	// return true
}
