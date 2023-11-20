package game

import (
	"HigherLevelPerudoServer/message_handlers"
	"HigherLevelPerudoServer/messages"
	"HigherLevelPerudoServer/util"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"golang.org/x/exp/maps"
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

type GamesMap map[string]*GameState

type MoveType string

const ( // ✓
	BET   MoveType = "Bet"
	DUDO  MoveType = "Dudo"
	CALZA MoveType = "Calza"
)

type PlayerMove struct {
	MoveType MoveType // "Bet", "Dudo", "Calza"
	Value    Bet
}
type Result string

const (
	DEC  Result = "dec"
	INC  Result = "inc"
	LOSE Result = "lose"
	WIN  Result = "win"
	NEXT Result = "next"
)

type RoundResult struct { // the result of the round
	PlayerIndex int
	Result      Result // "dec", "inc", "lose", "win", "next"
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

type PlayerHandsContents struct {
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
	gameState.PrevMove = PlayerMove{} // 'zero out' the previous move

	InitialPlayerHandLengths := PlayerHandLengthsUpdate{util.Map(func(x PlayerHand) int { return len(x) }, gameState.PlayerHands)}
	gameState.broadcast(messages.Message{TypeDescriptor: "PlayerHandLengthsUpdate", Contents: InitialPlayerHandLengths})
	// gameState.CurrentPlayerIndex = 0 //EVIL SIN CRIME GUILT FILTH UNWASHED
	gameState.broadcast(messages.Message{TypeDescriptor: "RoundResult", Contents: RoundResult{gameState.CurrentPlayerIndex, "next"}})
}

// Processes new player Move
// Returns validity of move
// Will update AllowableChannelLock after running
// Validates, Broadcasts RoundUpdate, Processes, Updates Index, Broadcasts RoundResult
func (gameState *GameState) ProcessPlayerMove(playerMove PlayerMove) bool {
	defer func() { gameState.AllowableChannelLock = gameState.CurrentPlayerIndex }()
	fmt.Printf("CPI %d || ACL %d \n", gameState.CurrentPlayerIndex, gameState.AllowableChannelLock)
	switch playerMove.MoveType {
	case BET:
		return gameState.processPlayerBet(playerMove)
	case DUDO:
		return gameState.processPlayerDudo()
	case CALZA:
		return gameState.processPlayerCalza()
	default:
		return false
	}
}

func (gameState *GameState) updatePlayerIndex(moveType MoveType, optional_player_lose ...int) error {
	if len(gameState.PlayerHands) == 0 {
		err := errors.New("can't update a Game with no players")
		return err
	} else if gameState.PlayersAllDead() {
		return errors.New("all players are dead")
	}
	if moveType == DUDO || moveType == CALZA {
		if len(optional_player_lose) != 1 {
			err := errors.New("a dudo or a calza always causes a player to have their number of dice change")
			return err
		} else {
			gameState.CurrentPlayerIndex = optional_player_lose[0]
		}
	} else if moveType == BET {
		gameState.CurrentPlayerIndex += 1
		gameState.CurrentPlayerIndex %= len(gameState.PlayerHands)
	}
	err := gameState.findNextAlivePlayerInclusive()
	return err
}

// broadcast message function, to used as: gameState.broadcast(message)

func (gameState *GameState) ProcessUserMessage(userMessage messages.Message, thisChan chan []byte, channelLocations *message_handlers.ChannelLocations, allGames *message_handlers.MessageHandlers) {

	// fmt.Println("Printing gamestate", gameState)
	// not very efficient. Should work
	fmt.Printf("Message recieved from websocket associated to player index %d \n", slices.Index(gameState.PlayerChannels, thisChan))
	switch userMessage.TypeDescriptor {
	case "PlayerMove":
		fmt.Println("Made it into PlayerMove switch")
		// If PlayerMove need to ensure that userMessage.Contents is of type PlayerMove

		// could check here to make sure that this message is coming from the current player
		// To do as such, we need a pairing from thisChan to playerIDs
		// Then check equality against gameState.CurrentPlayerIndex
		// thisChanIndex := slices.Index[[]chan []byte, chan []byte](gameState.PlayerChannels,thisChan)

		thisChanIndex := slices.Index(gameState.PlayerChannels, thisChan)
		if thisChanIndex != gameState.AllowableChannelLock {
			thisChan <- messages.PackMessage("NOT YOUR TURN", nil)
			return
		}

		var playerMove PlayerMove
		// the below marshal and Unmarshal could be removed with a good understanding of the type of PlayerMove being sent
		// Below may work but may also be a breaking change
		// playerMove = PlayerMove{MoveType = userMessage.Contents.MoveType, Bet userMessage.Contents.Bet}

		buff, err := json.Marshal(userMessage.Contents)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = json.Unmarshal(buff, &playerMove)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Calling gamestate.processPlayerMove")
		couldProcessMove := gameState.ProcessPlayerMove(playerMove)
		if !couldProcessMove {
			thisChan <- messages.PackMessage("Could not process player move", nil)
			return
		}
		fmt.Println("Finished gamestate.processPlayerMove")
		// if !valid {
		// 	gameState.PlayerChannels[gameState.CurrentPlayerIndex] <- packMessage("Invalid Bet", "Invalid Bet selection. Please select a valid move")
		// 	return
		// }
		// move was valid, broadcast new state
	case "GameStart":
		gameState.PlayerChannels = maps.Keys(*channelLocations)
		fmt.Println("Case: GameStart")
		gameState.StartNewGame()
		// will need to let players know the result of updating the game state
	}

}

func (gameState *GameState) AddChannel(thisChan chan []byte, channelLocations *message_handlers.ChannelLocations) {
	gameState.PlayerChannels = append(gameState.PlayerChannels, thisChan)
	(*channelLocations)[thisChan] = gameState
}

func (gameState *GameState) MoveChannel(thisChan chan []byte, newLocation message_handlers.MessageHandler, channelLocations *message_handlers.ChannelLocations, allHandlers *message_handlers.MessageHandlers) {
	// thisChanIndex := slices.Index(gameState.PlayerChannels, thisChan)
	// gameState.PlayerChannels = slices.Delete(gameState.PlayerChannels, thisChanIndex, thisChanIndex) // might need a +1 to make a valid slice
	// if len(gameState.PlayerChannels) == 0 && message_handlers.MessageHandler(gameState) != *newLocation {
	// 	delete((*allHandlers), gameState)
	// }
	// (*newLocation).AddChannel(thisChan, channelLocations)
	message_handlers.MoveChannelLogic(&gameState.PlayerChannels, thisChan, newLocation, channelLocations, allHandlers)
}
