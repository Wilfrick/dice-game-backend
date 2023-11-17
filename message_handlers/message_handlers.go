package message_handlers

import (
	"HigherLevelPerudoServer/messages"
)

type ChannelLocations map[chan []byte]MessageHandler

type MessageHandlers map[MessageHandler]struct{}

type MessageHandler interface {
	ProcessUserMessage(message messages.Message, thisChan chan []byte, channelLocations *ChannelLocations, allHandlers *MessageHandlers)
}

type UnassignedPlayerHandler struct {
	UnassignedPlayers []chan []byte
}

func (unPH *UnassignedPlayerHandler) ProcessUserMessage(msg messages.Message, thisChan chan []byte, channelLocations *ChannelLocations, allHandlers *MessageHandlers) {
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

type LobbyHandler struct {
	LobbyPlayerChannels []chan []byte
	LobbyID             string
}

func (lobbyHandler *LobbyHandler) ProcessUserMessage(msg messages.Message, thisChan chan []byte, channelLocations *ChannelLocations, allHandlers *MessageHandlers) {
	switch msg.TypeDescriptor {
	case "StartNewGame":
		// Lots here
		// In some order
		// Change the lobby to a game
		// Update the location of the players
		// Communicate this information. best handled here or in Game Handler
	}
}
