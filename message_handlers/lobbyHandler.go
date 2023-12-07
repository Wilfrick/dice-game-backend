package message_handlers

import (
	"HigherLevelPerudoServer/message_handlers/message_handler_interface"
	"HigherLevelPerudoServer/messages"
	"fmt"
	"slices"
)

type LobbyHandler struct {
	LobbyPlayerChannels           []chan []byte
	LobbyID                       string
	channelLocations              *message_handler_interface.ChannelLocations
	GlobalUnassignedPlayerHandler *UnassignedPlayerHandler
	IsQuickplay                   bool
}

const QUICKPLAY_LOBBY_SIZE int = 4

func (lobbyHandler *LobbyHandler) ProcessUserMessage(msg messages.Message, thisChan chan []byte) {
	// var _ game.GameState // wtf is this
	switch msg.TypeDescriptor {
	case "Start Game":
		if lobbyHandler.IsQuickplay {
			fmt.Println("A player tried to start a new game")
			msg := messages.Message{TypeDescriptor: "RejectBadMessage", Contents: "You can't force a quickplay game to start"}
			thisChan <- messages.CreateEncodedMessage(msg)
			return
		}
		lobbyHandler.StartGame()

	case "Leave Lobby":
		fmt.Println("Player trying to leave lobby")

		whichPlayer := slices.Index(lobbyHandler.LobbyPlayerChannels, thisChan)

		lobbyHandler.MoveChannel(thisChan, lobbyHandler.GlobalUnassignedPlayerHandler)
		// Not real code
		// Tell everyone who left
		playerLeftMessage := messages.Message{TypeDescriptor: "Player Left Lobby", Contents: struct {
			PlayerIndex    int
			NewPlayerCount int
		}{PlayerIndex: whichPlayer, NewPlayerCount: len(lobbyHandler.LobbyPlayerChannels)}}
		// fmt.Println(lobbyHandler)
		lobbyHandler.Broadcast(playerLeftMessage)
		// successfulLeaveMessage := Message{msg.TypeDescriptor: ""}
	}
}

func (lobbyHandler *LobbyHandler) StartGame() {
	// Lots here mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm
	// In some order
	// Change the lobby to a game
	// Update the location of the players
	// Communicate this information. best handled here or in Game Handler

	// create a new game (including setting channellocations)
	gameHandler := GameHandler{}
	gameHandler.SetChannelLocations(lobbyHandler.channelLocations)
	// numPlayers := len(lobbyHandler.LobbyPlayerChannels)
	// // add the players that are currently in this lobby to the game
	// for _, playerChan := range lobbyHandler.LobbyPlayerChannels {
	// 	lobbyHandler.MoveChannel(playerChan, &gameState)
	// }
	fmt.Println("Here is fine")
	for remaining := len(lobbyHandler.LobbyPlayerChannels); remaining > 0; remaining = len(lobbyHandler.LobbyPlayerChannels) {
		lobbyHandler.MoveChannel(lobbyHandler.LobbyPlayerChannels[remaining-1], &gameHandler)
	}
	// give the game the correct ID
	gameHandler.gameState.GameID = lobbyHandler.LobbyID
	gameHandler.GlobalUnassignedPlayerHandler = lobbyHandler.GlobalUnassignedPlayerHandler

	// update the channelLocations for these channels
	delete(*lobbyHandler.GlobalUnassignedPlayerHandler.LobbyMap, lobbyHandler.LobbyID)
	// done
	msg := messages.Message{TypeDescriptor: "Game Started", Contents: struct{ GameID string }{GameID: gameHandler.gameState.GameID}}
	gameHandler.gameState.Broadcast(msg)

	fmt.Println("Case: GameStart")
	gameHandler.gameState.StartNewGame()
}

func (lobbyHandler *LobbyHandler) AddChannel(thisChan chan []byte) {
	lobbyHandler.LobbyPlayerChannels = append(lobbyHandler.LobbyPlayerChannels, thisChan)
	(*lobbyHandler.channelLocations)[thisChan] = lobbyHandler
	var lobbyNavigation string
	if lobbyHandler.IsQuickplay {
		lobbyNavigation = "/quickplay"
	} else {
		lobbyNavigation = "/lobby"
	}
	playerLocationMessage := messages.Message{TypeDescriptor: "PlayerLocation", Contents: lobbyNavigation}
	thisChan <- messages.CreateEncodedMessage(playerLocationMessage)
	// Tell everyone new lobby player count
	numLobbyPlayers := len(lobbyHandler.LobbyPlayerChannels)
	lobbyJoinResponse := LobbyJoinResponse{userReadableResponse: "Successfully joined lobby with given ID", LobbyID: lobbyHandler.LobbyID, NumPlayers: numLobbyPlayers}
	msg := messages.Message{TypeDescriptor: "Lobby Join Accepted", Contents: lobbyJoinResponse}
	lobbyHandler.Broadcast(msg)
	if lobbyHandler.IsQuickplay && len(lobbyHandler.LobbyPlayerChannels) == QUICKPLAY_LOBBY_SIZE {
		// defer lobbyHandler.GlobalUnassignedPlayerHandler.CreateNewQuickPlay()
		fmt.Println("Commencing quickplay game")
		lobbyHandler.StartGame()
	}
}

func (lobbyHandler *LobbyHandler) MoveChannel(thisChan chan []byte, newLocation message_handler_interface.MessageHandler) {
	// thisChanIndex := slices.Index(lobbyHandler.LobbyPlayerChannels, thisChan)
	// lobbyHandler.LobbyPlayerChannels = slices.Delete(lobbyHandler.LobbyPlayerChannels, thisChanIndex, thisChanIndex+1) // might need a +1 to make a valid slice
	// if len(lobbyHandler.LobbyPlayerChannels) == 0 && message_handlers.MessageHandler(lobbyHandler) != *newLocation {
	// 	delete((*allHandlers), lobbyHandler)
	// }
	// (*newLocation).AddChannel(thisChan, channelLocations)
	message_handler_interface.MoveChannelLogic(&lobbyHandler.LobbyPlayerChannels, thisChan, newLocation, lobbyHandler.channelLocations)
	if len(lobbyHandler.LobbyPlayerChannels) == 0 {
		delete((*lobbyHandler.GlobalUnassignedPlayerHandler.LobbyMap), lobbyHandler.LobbyID)
	}
}

func (lobbyHandler *LobbyHandler) SetChannelLocations(channelLocations *message_handler_interface.ChannelLocations) {
	lobbyHandler.channelLocations = channelLocations
}
func (lobbyHandler LobbyHandler) Broadcast(message messages.Message, optional_use_wait_group ...bool) {
	message_handler_interface.BroadcastLogic(lobbyHandler.LobbyPlayerChannels, message, optional_use_wait_group...)
}

func (lobbyHandler *LobbyHandler) RemoveChannel(thisChan chan []byte) {
	message_handler_interface.RemoveChannelLogic(&lobbyHandler.LobbyPlayerChannels, thisChan, lobbyHandler.channelLocations)
}
