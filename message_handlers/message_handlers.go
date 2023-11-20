package message_handlers

import (
	"HigherLevelPerudoServer/messages"
	"slices"
)

type ChannelLocations map[chan []byte]MessageHandler

type MessageHandlers map[MessageHandler]struct{}

type MessageHandler interface {
	ProcessUserMessage(message messages.Message, thisChan chan []byte, channelLocations *ChannelLocations, allHandlers *MessageHandlers)
	AddChannel(thisChan chan []byte, channelLocations *ChannelLocations)
	MoveChannel(thisChan chan []byte, newLocation MessageHandler, channelLocations *ChannelLocations, allHandlers *MessageHandlers)
}

func MoveChannelLogic(sliceOfPlayerChannels *[]chan []byte, thisChan chan []byte, newLocation MessageHandler, channelLocations *ChannelLocations, allHandlers *MessageHandlers) {
	thisHandler := (*channelLocations)[thisChan]
	thisChanIndex := slices.Index(*sliceOfPlayerChannels, thisChan)
	*sliceOfPlayerChannels = slices.Delete(*sliceOfPlayerChannels, thisChanIndex, thisChanIndex+1) // might need a +1 to make a valid slice
	if newLocation == nil {
		if len(*sliceOfPlayerChannels) == 0 {
			delete((*allHandlers), thisHandler)
		}
		return
	}
	if len(*sliceOfPlayerChannels) == 0 && thisHandler != newLocation {
		delete((*allHandlers), thisHandler)
	}
	(newLocation).AddChannel(thisChan, channelLocations)
}
