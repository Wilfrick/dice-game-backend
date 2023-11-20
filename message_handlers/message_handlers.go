package message_handlers

import (
	"HigherLevelPerudoServer/messages"
	"slices"
)

type ChannelLocations map[chan []byte]MessageHandler

type MessageHandlers map[MessageHandler]struct{}

type MessageHandler interface {
	ProcessUserMessage(message messages.Message, thisChan chan []byte)
	AddChannel(thisChan chan []byte)
	MoveChannel(thisChan chan []byte, newLocation MessageHandler)
	SetChannelLocations(*ChannelLocations) // A message handler exists in the scope of a channelLocation
}

func MoveChannelLogic(sliceOfPlayerChannels *[]chan []byte, thisChan chan []byte, newLocation MessageHandler, channelLocations *ChannelLocations) {
	// thisHandler := (*channelLocations)[thisChan]
	thisChanIndex := slices.Index(*sliceOfPlayerChannels, thisChan)
	if thisChanIndex != -1 {
		*sliceOfPlayerChannels = slices.Delete(*sliceOfPlayerChannels, thisChanIndex, thisChanIndex+1) // might need a +1 to make a valid slice
	}
	if newLocation == nil {
		// if len(*sliceOfPlayerChannels) == 0 {
		// 	// delete((*allHandlers), thisHandler)
		// }
		return
	}
	// if len(*sliceOfPlayerChannels) == 0 && thisHandler != newLocation {
	// 	// delete((*allHandlers), thisHandler)
	// }
	(newLocation).AddChannel(thisChan)
}

// func (lobbyHandler *T) AddChannel[T](thisChan chan []byte) { // WE WANT THIS, but it doesn't exist yet
// 	lobbyHandler.LobbyPlayerChannels = append(lobbyHandler.LobbyPlayerChannels, thisChan)
// 	(*lobbyHandler.channelLocations)[thisChan] = lobbyHandler
// }
