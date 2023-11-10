package game

import (
	"HigherLevelPerudoServer/util"
	"errors"
	"fmt"
	"slices"
	"sync"
)

type GameState struct {
	GameID             string
	PlayerHands        []PlayerHand
	CurrentPlayerIndex int // should be < len(PlayerHands)
	PlayerChannels     []chan []byte
	PrevMove           PlayerMove
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
	for i := range gameState.PlayerHands {
		gameState.PlayerHands[i] = RandomPlayerHand(5)
	}
	gameState.distributeHands()

	InitialPlayerHandLengths := PlayerHandLengthsUpdate{util.Map(func(x PlayerHand) int { return len(x) }, gameState.PlayerHands)}
	gameState.broadcast(Message{TypeDescriptor: "PlayerHandLengthsUpdate", Contents: InitialPlayerHandLengths})
	gameState.CurrentPlayerIndex = 0
	gameState.broadcast(Message{"RoundResult", RoundResult{gameState.CurrentPlayerIndex, "next"}})

	fmt.Println("Ended StartNewGame")
}

func (gameState *GameState) ProcessPlayerMove(playerMove PlayerMove) bool {
	switch playerMove.MoveType {
	case "Bet":
		fmt.Println("in ProcessPlayerMove, made into case 'Bet' ")
		fmt.Println(*gameState)
		fmt.Println(playerMove)
		// validate bet
		newBet := playerMove.Value
		betValid := newBet.isGreaterThan(gameState.PrevMove.Value)
		betValid = true
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
		fmt.Println("in ProcessPlayerMove, made into case 'Dudo' ")
		fmt.Println(*gameState)
		fmt.Println(playerMove)
		// validate bet
		// dudo should always be valid, as long as the current player
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
			}
		}

		gameState.broadcast(Message{"RoundResult", RoundResult{losing_player_index, "dec"}})

		player_died, err := gameState.RemoveDice(losing_player_index)
		if err != nil {
			fmt.Println(err.Error())
		}
		if player_died {
			gameState.broadcast(Message{"RoundResult", RoundResult{losing_player_index, "lose"}})
			gameState.send(losing_player_index, Message{"GameResult", GameResult{losing_player_index, "lose"}})
		}
		// 2) possible who has now lost ()
		gameState.updatePlayerIndex(playerMove)

		gameState.randomiseCurrentHands()

		gameState.distributeHands() // send hand messages to each player with personalised hand info
		// 3) Send new hands out to all players (likely different messages to different players)

		// 4) Inform players who's turn is next
		gameState.broadcast(Message{"RoundResult", RoundResult{gameState.CurrentPlayerIndex, "next"}})
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
	// encodedMessage := createEncodedMessage(message)
	// for _, channel := range gameState.PlayerChannels {
	// 	fmt.Println("Sending message")
	// 	go func(c chan []byte) { c <- encodedMessage }(channel)
	// 	// go func (){channel <- encodedMessage}() // likely to lead to untracable bugs, do not copy
	// }
	for player_index := range gameState.PlayerChannels {
		fmt.Println("Sending message")
		gameState.send(player_index, message)
		// go func (){channel <- encodedMessage}() // likely to lead to untracable bugs, do not copy
	}
	fmt.Println("Finished broadcasting message")
}

func (gameState GameState) isBetTrue() bool {
	fmt.Println("Called isBetTrue")
	// TODO, for Dudo
	bet := gameState.PrevMove.Value
	dice_face_counts := count_dice_faces_considering_wild_ones(gameState)

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

func (gameState GameState) PreviousAlivePlayer() (int, error) {
	fmt.Println("Called PreviousAlivePlayer")
	alive_player_indices := util.Filter(func(player_index int) bool { return player_index >= 0 },
		util.Mapi(func(p PlayerHand, index int) int {
			if len(p) > 0 {
				return len(p)
			} else {
				return -1
			}
		}, gameState.PlayerHands))

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

func (gameState *GameState) RemoveDice(player_index int) (bool, error) {
	fmt.Println("Called RemoveDice")
	// should do bounds checking for player_index
	if player_index < 0 || player_index >= len(gameState.PlayerHands) {
		return false, errors.New("bad player index")
	}
	if len(gameState.PlayerHands[player_index]) <= 0 {
		return true, errors.New("this player is already dead")
	}

	gameState.PlayerHands = append(make([]PlayerHand, 0, len(gameState.PlayerHands[player_index])), gameState.PlayerHands[:len(gameState.PlayerHands[player_index])-1]...)

	// returns true if the player that lost a dice is now dead
	return len(gameState.PlayerHands[player_index]) == 0, nil
}

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

func (gameState *GameState) randomiseCurrentHands() {
	fmt.Println("Called randomiseCurrentHands")
	for _, hand := range gameState.PlayerHands {
		hand.Randomise()
	}
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
