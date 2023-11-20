package player_management_handlers

import (
	"HigherLevelPerudoServer/message_handlers"
	"HigherLevelPerudoServer/messages"
	"fmt"

	"golang.org/x/exp/maps"
)

type UnassignedPlayerHandler struct {
	UnassignedPlayers     []chan []byte
	currentQuickplayLobby *LobbyHandler
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
		hash := "abcdefghiklmnopqrtsuvwxyz"
		// Check for collision till new // very low prob that this is the same as another hash, so we all good
		// Create lobby with hash
		newLobby := LobbyHandler{LobbyID: hash}
		(*allHandlers)[&newLobby] = struct{}{}
		unPH.MoveChannel(thisChan, &newLobby, channelLocations, allHandlers)
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

func (unPH *UnassignedPlayerHandler) MoveChannel(thisChan chan []byte, newLocation message_handlers.MessageHandler, channelLocations *message_handlers.ChannelLocations, allHandlers *message_handlers.MessageHandlers) {
	message_handlers.MoveChannelLogic(&unPH.UnassignedPlayers, thisChan, newLocation, channelLocations, allHandlers)
}

type LobbyHandler struct {
	LobbyPlayerChannels []chan []byte
	LobbyID             string
}

func (lobbyHandler *LobbyHandler) ProcessUserMessage(msg messages.Message, thisChan chan []byte, channelLocations *message_handlers.ChannelLocations, allHandlers *message_handlers.MessageHandlers) {
	// var _ game.GameState // wtf is this
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

func (lobbyHandler *LobbyHandler) MoveChannel(thisChan chan []byte, newLocation message_handlers.MessageHandler, channelLocations *message_handlers.ChannelLocations, allHandlers *message_handlers.MessageHandlers) {
	// thisChanIndex := slices.Index(lobbyHandler.LobbyPlayerChannels, thisChan)
	// lobbyHandler.LobbyPlayerChannels = slices.Delete(lobbyHandler.LobbyPlayerChannels, thisChanIndex, thisChanIndex+1) // might need a +1 to make a valid slice
	// if len(lobbyHandler.LobbyPlayerChannels) == 0 && message_handlers.MessageHandler(lobbyHandler) != *newLocation {
	// 	delete((*allHandlers), lobbyHandler)
	// }
	// (*newLocation).AddChannel(thisChan, channelLocations)
	message_handlers.MoveChannelLogic(&lobbyHandler.LobbyPlayerChannels, thisChan, newLocation, channelLocations, allHandlers)
}
