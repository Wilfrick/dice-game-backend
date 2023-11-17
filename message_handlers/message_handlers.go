package message_handlers

import (
	"HigherLevelPerudoServer/messages"
)

type ChannelLocations map[chan []byte]MessageHandler

type MessageHandlers map[MessageHandler]struct{}

type MessageHandler interface {
	ProcessUserMessage(message messages.Message, thisChan chan []byte, channelLocations *ChannelLocations, allHandlers *MessageHandlers)
	AddChannel(thisChan chan []byte, channelLocations *ChannelLocations)
}
