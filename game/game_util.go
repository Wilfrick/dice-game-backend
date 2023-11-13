package game

import (
	"errors"
	"fmt"
)

func (gameState *GameState) removeDice(player_index int) (bool, error) {
	fmt.Println("Called RemoveDice")
	// N.B: First time player 1 dice?
	// should do bounds checking for player_index
	if player_index < 0 || player_index >= len(gameState.PlayerHands) {
		return false, errors.New("bad player index")
	}
	if len(gameState.PlayerHands[player_index]) <= 0 {
		return true, errors.New("this player is already dead")
	}
	//Trim this players hands
	gameState.PlayerHands[player_index] = gameState.PlayerHands[player_index][:len(gameState.PlayerHands[player_index])-1]

	// returns true if the player that lost a dice is now dead
	return len(gameState.PlayerHands[player_index]) == 0, nil
}

func (gameState *GameState) randomiseCurrentHands() {
	fmt.Println("Called randomiseCurrentHands")
	for _, hand := range gameState.PlayerHands {
		hand.Randomise()
	}
}

func (gameState *GameState) findNextAlivePlayerInclusive() error {
	startingIndex := gameState.CurrentPlayerIndex

	playerDead := len(gameState.PlayerHands[gameState.CurrentPlayerIndex]) == 0
	for playerDead {
		gameState.CurrentPlayerIndex += 1
		gameState.CurrentPlayerIndex %= len(gameState.PlayerHands)
		if gameState.CurrentPlayerIndex == startingIndex {
			err := errors.New("passed a game that has already been won")
			return err
		}
		playerDead = len(gameState.PlayerHands[gameState.CurrentPlayerIndex]) == 0
	}
	return nil
}
