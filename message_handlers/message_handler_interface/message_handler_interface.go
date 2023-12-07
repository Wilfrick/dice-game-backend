package message_handler_interface

import (
	"HigherLevelPerudoServer/messages"
	"slices"
	"sync"
)

type ChannelLocations map[chan []byte]MessageHandler

type MessageHandlers map[MessageHandler]struct{}

type MessageHandler interface {
	ProcessUserMessage(message messages.Message, thisChan chan []byte)
	Broadcast(message messages.Message, optional_wait_group ...bool)
	AddChannel(thisChan chan []byte)
	MoveChannel(thisChan chan []byte, newLocation MessageHandler)
	RemoveChannel(thisChan chan []byte)
	SetChannelLocations(*ChannelLocations) // A message handler exists in the scope of a channelLocation
}

func RemoveChannelLogic(sliceOfPlayerChannels *[]chan []byte, thisChan chan []byte, channelLocations *ChannelLocations) {
	MoveChannelLogic(sliceOfPlayerChannels, thisChan, nil, channelLocations)
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
	(newLocation).AddChannel(thisChan)
}

// func (lobbyHandler *T) AddChannel[T](thisChan chan []byte) { // WE WANT THIS, but it doesn't exist yet
//
//		lobbyHandler.LobbyPlayerChannels = append(lobbyHandler.LobbyPlayerChannels, thisChan)
//		(*lobbyHandler.channelLocations)[thisChan] = lobbyHandler
//	}
func BroadcastLogic(sliceOfPlayerChannels []chan []byte, message messages.Message, optional_use_wait_group ...bool) {
	var wait_group sync.WaitGroup

	use_wait_group := len(optional_use_wait_group) == 1 && optional_use_wait_group[0]
	encodedMessage := messages.CreateEncodedMessage(message)
	for _, playerChan := range sliceOfPlayerChannels {
		// fmt.Println("Sending message")
		if use_wait_group {
			sendBytes(playerChan, encodedMessage, &wait_group)
		} else {
			sendBytes(playerChan, encodedMessage)
		}
	}
	if use_wait_group {
		wait_group.Wait()
	}
}

func sendBytes(thisChan chan []byte, bytesContents []byte, optional_wait_group ...*sync.WaitGroup) {

	var wait_group sync.WaitGroup
	if len(optional_wait_group) == 1 {
		wait_group := (optional_wait_group[0])
		wait_group.Add(1)
	}
	// fmt.Println("going sending", string(bytesContents))
	go func(thisChan chan []byte, encodedMsg []byte) {
		thisChan <- encodedMsg
		if len(optional_wait_group) == 1 {
			wait_group.Done()
		}

	}(thisChan, bytesContents)

}

func Send(thisChan chan []byte, message messages.Message, optional_wait_group ...*sync.WaitGroup) {
	encodedMessage := messages.CreateEncodedMessage(message)
	sendBytes(thisChan, encodedMessage, optional_wait_group...)
}
