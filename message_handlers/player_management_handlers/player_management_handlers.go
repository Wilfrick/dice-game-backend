package player_management_handlers

import (
	"HigherLevelPerudoServer/game"
	"HigherLevelPerudoServer/message_handlers"
	"HigherLevelPerudoServer/messages"
	"fmt"

	"golang.org/x/exp/maps"
)

type UnassignedPlayerHandler struct {
	UnassignedPlayers []chan []byte
}

func (unPH *UnassignedPlayerHandler) ProcessUserMessage(msg messages.Message, thisChan chan []byte,
	channelLocations *message_handlers.ChannelLocations, allHandlers *message_handlers.MessageHandlers) {
	fmt.Println("Unassigned player gave a message")
	gs1 := maps.Keys(*allHandlers)[0]
	// gs := game.GameState(gs1)
	gs1.AddChannel(thisChan, channelLocations)
	// (*channelLocations)[thisChan] = &gs
	switch msg.TypeDescriptor {

	case "Quickplay":
		// do quickplay
	case "Create lobby":
		// Generate a random hash
		// Check for collision till new
		// Create lobby with hash
		// Add player to lobby
	case "Join Lobby":
		// Look in contents for target lobby
		// ADd player to target lobby
	}
}

func (unPH *UnassignedPlayerHandler) AddChannel(thisChan chan []byte, channelLocations *message_handlers.ChannelLocations) {
	unPH.UnassignedPlayers = append(unPH.UnassignedPlayers, thisChan)
	(*channelLocations)[thisChan] = unPH
}

type LobbyHandler struct {
	LobbyPlayerChannels []chan []byte
	LobbyID             string
}

func (lobbyHandler *LobbyHandler) ProcessUserMessage(msg messages.Message, thisChan chan []byte, channelLocations *message_handlers.ChannelLocations, allHandlers *message_handlers.MessageHandlers) {
	var _ game.GameState
	switch msg.TypeDescriptor {
	case "StartNewGame":
		// Lots here
		// In some order
		// Change the lobby to a game
		// Update the location of the players
		// Communicate this information. best handled here or in Game Handler

	}
}

func (lobbyHandler *LobbyHandler) AddChannel(thisChan chan []byte, channelLocations *message_handlers.ChannelLocations) {
	lobbyHandler.LobbyPlayerChannels = append(lobbyHandler.LobbyPlayerChannels, thisChan)
	(*channelLocations)[thisChan] = lobbyHandler
}
