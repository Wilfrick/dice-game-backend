package message_handlers

import (
	"HigherLevelPerudoServer/message_handlers/message_handler_interface"
	"HigherLevelPerudoServer/messages"
	"HigherLevelPerudoServer/util"
	"testing"
)

func Test_unPH(t *testing.T) {
	unPH := UnassignedPlayerHandler{}
	// lobbyChan := LobbyHandler{}
	playerChan := make(chan []byte)
	util.ChanSink([]chan []byte{playerChan})
	channelLocations := message_handler_interface.ChannelLocations{}
	// allHandlers := message_handler_interface.MessageHandlers{} // moved away from this a while ago
	// allHandlers[&unPH] = struct{}{}
	unPH.SetChannelLocations(&channelLocations)
	unPH.AddChannel(playerChan)
	t.Log(unPH.UnassignedPlayers, playerChan)
	util.Assert(t, unPH.UnassignedPlayers[0] == playerChan)
	util.Assert(t, len(unPH.UnassignedPlayers) == 1)
	unPH.MoveChannel(playerChan, nil)
	t.Log(unPH.UnassignedPlayers, playerChan)
	util.Assert(t, len(unPH.UnassignedPlayers) == 0)
}

func Test_newLocation(t *testing.T) {
	unPH := UnassignedPlayerHandler{}
	lobbyChan := LobbyHandler{}
	playerChan := make(chan []byte)
	util.ChanSink([]chan []byte{playerChan})
	channelLocations := message_handler_interface.ChannelLocations{}
	// allHandlers := message_handlers.MessageHandlers{}
	// allHandlers[&unPH] = struct{}{}
	// allHandlers[&lobbyChan] = struct{}{}
	unPH.SetChannelLocations(&channelLocations)
	util.Assert(t, unPH.channelLocations == &channelLocations)
	unPH.AddChannel(playerChan)
	util.Assert(t, unPH.UnassignedPlayers[0] == playerChan)
	util.Assert(t, len(unPH.UnassignedPlayers) == 1)
	lobbyChan.SetChannelLocations(&channelLocations)
	unPH.MoveChannel(playerChan, &lobbyChan)
	util.Assert(t, len(unPH.UnassignedPlayers) == 0)
	util.Assert(t, len(lobbyChan.LobbyPlayerChannels) == 1)
	util.Assert(t, lobbyChan.LobbyPlayerChannels[0] == playerChan)
}

func Test_createGame(t *testing.T) {
	playerChan := make(chan []byte)
	util.ChanSink([]chan []byte{playerChan})
	lobby := LobbyHandler{LobbyPlayerChannels: []chan []byte{playerChan}, LobbyID: "alex"}
	lobbyMap := make(map[string]*LobbyHandler)
	unPH := UnassignedPlayerHandler{}

	lobby.GlobalUnassignedPlayerHandler = &unPH
	msg := messages.Message{TypeDescriptor: "Start Game"}
	channelLocations := message_handler_interface.ChannelLocations{}
	channelLocations[playerChan] = &lobby
	unPH.SetChannelLocations(&channelLocations)
	unPH.LobbyMap = &lobbyMap

	// allHandlers := message_handlers.MessageHandlers{} // needed so tests compile atm, but will be removed in future
	lobby.SetChannelLocations(&channelLocations)
	lobby.ProcessUserMessage(msg, playerChan)
	// util.Assert(t, len(allHandlers) == 2)
	gh, ok := channelLocations[playerChan]
	t.Log(ok, lobby.LobbyPlayerChannels)
	util.Assert(t, ok)
	util.Assert(t, len(lobby.LobbyPlayerChannels) == 0)
	cast_game_handler, ok := gh.(*GameHandler)
	if !ok {
		t.FailNow()
	}
	util.Assert(t, ok)
	util.Assert(t, cast_game_handler.gameState.PlayerChannels[0] == playerChan)
	util.Assert(t, cast_game_handler.gameState.GameID == "alex")
}
