package handshake

import (
	"HigherLevelPerudoServer/message_handlers"
	"HigherLevelPerudoServer/message_handlers/player_management_handlers"
	"HigherLevelPerudoServer/util"
	"testing"
)

func Test_findAndInitialiseClientIDNewClient(t *testing.T) {
	//SETUP
	clientIDtoChannels := make(map[string]chan []byte)
	unPH := player_management_handlers.UnassignedPlayerHandler{}
	channelLocations := message_handlers.ChannelLocations{}
	//SETUP OVER
	initial_clientID := " "
	new_clientID, thisChan := findAndInitialiseClientID(initial_clientID, &clientIDtoChannels, &unPH, &channelLocations)

	// ASSERT
	util.Assert(t, len(new_clientID) > 2)
	util.Assert(t, clientIDtoChannels[new_clientID] == thisChan)
	util.Assert(t, channelLocations[thisChan] == &unPH)
}

func Test_findAndInitialiseClientIDOldClientInvalid(t *testing.T) {
	//SETUP
	clientIDtoChannels := make(map[string]chan []byte)
	unPH := player_management_handlers.UnassignedPlayerHandler{}
	channelLocations := message_handlers.ChannelLocations{}
	//SETUP OVER
	initial_clientID := "arccoeuhtdco" // e.g. server restarted
	new_clientID, thisChan := findAndInitialiseClientID(initial_clientID, &clientIDtoChannels, &unPH, &channelLocations)

	// ASSERT
	util.Assert(t, initial_clientID != new_clientID) //has a small chance to fail even when working correctly
	util.Assert(t, clientIDtoChannels[new_clientID] == thisChan)
	util.Assert(t, channelLocations[thisChan] == &unPH)
}

func Test_findAndInitialiseClientIDPreviousClientID(t *testing.T) {
	//SETUP
	clientIDtoChannels := make(map[string]chan []byte)
	thisChan := make(chan []byte)
	initial_clientID := "arccoeuhtdco"
	clientIDtoChannels[initial_clientID] = thisChan
	// unPH.UnassignedPlayers =
	unPH := player_management_handlers.UnassignedPlayerHandler{UnassignedPlayers: []chan []byte{thisChan}}
	channelLocations := message_handlers.ChannelLocations{thisChan: &unPH}
	//SETUP OVER
	// e.g. server restarted
	new_clientID, new_thisChan := findAndInitialiseClientID(initial_clientID, &clientIDtoChannels, &unPH, &channelLocations)

	// ASSERT
	util.Assert(t, initial_clientID == new_clientID)
	util.Assert(t, clientIDtoChannels[new_clientID] == thisChan)
	util.Assert(t, channelLocations[new_thisChan] == &unPH)
	util.Assert(t, unPH.UnassignedPlayers[0] == thisChan)
	util.Assert(t, len(unPH.UnassignedPlayers) == 1)
	util.Assert(t, new_thisChan == thisChan)

}

func Test_findAndInitialiseClientIDFirstTimeThenAgain(t *testing.T) {
	//SETUP
	clientIDtoChannels := make(map[string]chan []byte)
	initial_clientID := " "

	unPH := player_management_handlers.UnassignedPlayerHandler{}
	channelLocations := message_handlers.ChannelLocations{}
	//SETUP OVER
	// e.g. server restarted
	new_clientID, thisChan := findAndInitialiseClientID(initial_clientID, &clientIDtoChannels, &unPH, &channelLocations)

	// ASSERT
	util.Assert(t, initial_clientID != new_clientID)
	util.Assert(t, clientIDtoChannels[new_clientID] == thisChan)
	util.Assert(t, channelLocations[thisChan] == &unPH)

	newer_clientID, otherChan := findAndInitialiseClientID(new_clientID, &clientIDtoChannels, &unPH, &channelLocations)

	util.Assert(t, newer_clientID == new_clientID)
	util.Assert(t, thisChan == otherChan)
	util.Assert(t, len(unPH.UnassignedPlayers) == 1)
}
