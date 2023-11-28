package game

import (
	"errors"
	"fmt"
	"slices"
)

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

func (gameState GameState) isBetTrue() bool {
	fmt.Println("Called isBetTrue")
	// TODO, for Dudo
	bet := gameState.PrevMove.Value
	var dice_face_counts []int
	if !gameState.IsPalacifoRound {
		dice_face_counts = count_dice_faces_considering_wild_ones(gameState)
	} else {
		dice_face_counts = count_dice_faces(gameState)
	}

	return dice_face_counts[bet.FaceVal] >= bet.NumDice
}

func (gameState GameState) isBetExactlyTrue() bool {
	fmt.Println("Called isBetExactlyTrue")
	// TODO, for Calza
	bet := gameState.PrevMove.Value
	dice_face_counts := count_dice_faces_considering_wild_ones(gameState)

	return dice_face_counts[bet.FaceVal] == bet.NumDice
}

func count_dice_faces(gameState GameState) []int {
	dice_face_counts := make([]int, 7) // 7 is a magic number, it is one more than 6, which is the number of possible values that we allow for a dice. This could be refactored with a map to handle arbitrary dice face values.
	for _, player_hand := range gameState.PlayerHands {
		for _, dice_value := range player_hand {
			dice_face_counts[dice_value]++
		}
	}
	return dice_face_counts
}

func count_dice_faces_considering_wild_ones(gameState GameState) []int {
	WILD_DICE_FACE_VALUE := 1
	dice_face_counts := count_dice_faces(gameState)
	for face_value := range dice_face_counts {
		if face_value != WILD_DICE_FACE_VALUE {
			dice_face_counts[face_value] += dice_face_counts[WILD_DICE_FACE_VALUE] // ones are wild in this game
		}
	}
	return dice_face_counts
}

func (gameState GameState) alivePlayerIndices() []int {
	// WRONG
	// return util.Filter(func(player_index int) bool { return player_index > 0 },
	// 	util.Map(func(p PlayerHand) int { return len(p) }, gameState.PlayerHands))
	alivePlayerIndices := make([]int, 0, len(gameState.PlayerHands))
	for i, playerHand := range gameState.PlayerHands {
		if len(playerHand) > 0 {
			alivePlayerIndices = append(alivePlayerIndices, i)
		}
	}
	return alivePlayerIndices
}

func (gameState GameState) PreviousAlivePlayer() (int, error) {
	fmt.Println("Called PreviousAlivePlayer")
	alive_player_indices := gameState.alivePlayerIndices()

	if len(alive_player_indices) <= 1 {
		return -1, errors.New("not enough alive players") // the game should have already finished by now
	}
	current_player_relative_position := slices.Index(alive_player_indices, gameState.CurrentPlayerIndex)
	if current_player_relative_position == -1 {
		return -1, errors.New("couldn't find current player in the list of alive players")
	}

	if current_player_relative_position == 0 {
		return alive_player_indices[len(alive_player_indices)-1], nil
	}
	return alive_player_indices[current_player_relative_position-1], nil
}

func (gameState GameState) checkPlayerWin(candidate_victor int) bool {
	alivePlayers := gameState.alivePlayerIndices()
	return len(alivePlayers) == 1 && alivePlayers[0] == candidate_victor
}
