package game

import (
	"HigherLevelPerudoServer/messages"
	"HigherLevelPerudoServer/util"
	"errors"
	"fmt"
	"slices"
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

func (gameState *GameState) addDice(player_index int) {
	if len(gameState.PlayerHands[player_index]) < 5 {
		gameState.PlayerHands[player_index] = append(gameState.PlayerHands[player_index], 0)
	}
}

func (gameState *GameState) randomiseCurrentHands() {
	fmt.Println("Called randomiseCurrentHands")
	for _, hand := range gameState.PlayerHands {
		hand.Randomise()
	}
}

func (gameState *GameState) FindNextAlivePlayerInclusive() error {
	// startingIndex := gameState.CurrentPlayerIndex
	alivePlayerIndices := gameState.alivePlayerIndices()
	if len(alivePlayerIndices) < 1 {
		err := errors.New("passed a game that has already been won")
		return err
	}
	indexInAlivePlayerIndices, found := slices.BinarySearch(alivePlayerIndices, gameState.CurrentPlayerIndex)
	if found {
		return nil
	}
	// fmt.Println(indexInAlivePlayerIndices, alivePlayerIndices)
	gameState.CurrentPlayerIndex = alivePlayerIndices[(indexInAlivePlayerIndices)%len(alivePlayerIndices)]
	return nil
	// playerDead := len(gameState.PlayerHands[startingIndex]) == 0 && slices.Index(alivePlayerIndices, startingIndex) != -1
	// for playerDead {
	// 	gameState.CurrentPlayerIndex += 1
	// 	gameState.CurrentPlayerIndex %= len(gameState.PlayerHands)
	// 	if gameState.CurrentPlayerIndex == startingIndex {
	// 		err := errors.New("passed a game that has already been won")
	// 		return err
	// 	}
	// 	playerDead = len(gameState.PlayerHands[gameState.CurrentPlayerIndex]) == 0 && slices.Index(alivePlayerIndices, gameState.CurrentPlayerIndex) != -1
	// }
	// return nil
}

func (gameState *GameState) RemovePlayer(playerIndex int) error {
	if playerIndex >= len(gameState.PlayerChannels) {
		err := errors.New("attempted to remove a player lying beyond the channels")
		return err
	}
	gameState.PlayerHands[playerIndex] = PlayerHand{}
	if gameState.CurrentPlayerIndex == playerIndex {
		err := gameState.FindNextAlivePlayerInclusive()
		if err != nil {
			return err
		}
	}
	alivePlayerIndices := gameState.alivePlayerIndices()
	if len(alivePlayerIndices) == 1 {
		gameState.GameInProgress = false
		victor := alivePlayerIndices[0]
		gameState.send(victor, messages.Message{TypeDescriptor: "GameResult", Contents: GameResult{victor, "win"}})
	}
	gameState.PlayerHands = slices.Delete(gameState.PlayerHands, playerIndex, playerIndex+1)

	return nil
}
func (gameState *GameState) CleanUpInactivePlayers() {
	for i := len(gameState.PlayerChannels) - 1; i >= 0; i-- {
		if gameState.PlayerChannels[i] == nil {
			if i < gameState.CurrentPlayerIndex {
				gameState.CurrentPlayerIndex--
			}
			gameState.PlayerHands = slices.Delete(gameState.PlayerHands, i, i+1)
			gameState.PlayerChannels = slices.Delete(gameState.PlayerChannels, i, i+1)
			gameState.PalacifoablePlayers = slices.Delete(gameState.PalacifoablePlayers, i, i+1)
		}
	}
}

func (gameState *GameState) InitialiseSlicesWithDefaults() {
	max_slice_length := max(len(gameState.PlayerChannels), len(gameState.PlayerHands), len(gameState.PalacifoablePlayers))
	for len(gameState.PlayerChannels) < max_slice_length {
		c := make(chan []byte)
		util.ChanSink([]chan []byte{c})
		gameState.PlayerChannels = append(gameState.PlayerChannels, c)
	}
	for len(gameState.PlayerHands) < max_slice_length {
		gameState.PlayerHands = append(gameState.PlayerHands, PlayerHand{})
	}
	for len(gameState.PalacifoablePlayers) < max_slice_length {
		gameState.PalacifoablePlayers = append(gameState.PalacifoablePlayers, false)
	}
}
