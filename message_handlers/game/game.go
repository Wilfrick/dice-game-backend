package game

import (
	"HigherLevelPerudoServer/messages"
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
	// channelLocations              *message_handlers.ChannelLocations
	PalacifoablePlayers []bool
	IsPalacifoRound     bool
	// GlobalUnassignedPlayerHandler message_handlers.MessageHandler
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
	FinalBet    PlayerMove
}

type PalacifoRoundMessage bool

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

	// gameState.PlayerChannels = maps.Keys(*gameState.channelLocations) // PlayerChannels were already set from the lobby

	gameState.PlayerHands = make([]PlayerHand, len(gameState.PlayerChannels))
	gameState.PalacifoablePlayers = make([]bool, len(gameState.PlayerHands))
	for i := range gameState.PalacifoablePlayers {
		gameState.PalacifoablePlayers[i] = true
	}
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
	// Zero out previous
	gameState.PrevMove = PlayerMove{} // 'zero out' the previous move
	gameState.IsPalacifoRound = false

	InitialPlayerHandLengths := PlayerHandLengthsUpdate{util.Map(func(x PlayerHand) int { return len(x) }, gameState.PlayerHands)}

	for hand_index, hand_length := range InitialPlayerHandLengths.PlayerHandLengths {
		if hand_length == 1 && gameState.PalacifoablePlayers[hand_index] {
			gameState.IsPalacifoRound = true
			gameState.PalacifoablePlayers[hand_index] = false
			// could have a break here, but don't need to in this specific case
		}
	}

	gameState.Broadcast(messages.Message{TypeDescriptor: "PlayerHandLengthsUpdate", Contents: InitialPlayerHandLengths})
	// gameState.CurrentPlayerIndex = 0 //EVIL SIN CRIME GUILT FILTH UNWASHED
	gameState.Broadcast(messages.Message{TypeDescriptor: "RoundResult", Contents: RoundResult{gameState.CurrentPlayerIndex, "next"}})
	gameState.Broadcast(messages.Message{TypeDescriptor: "PalacifoRound", Contents: gameState.IsPalacifoRound})
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
