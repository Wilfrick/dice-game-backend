package message_handlers

import (
	"HigherLevelPerudoServer/message_handlers/message_handler_interface"
	"HigherLevelPerudoServer/messages"
	"HigherLevelPerudoServer/util"
	"slices"
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

func Test_lobbyToGameLocation(t *testing.T) {
	channelLocations := message_handler_interface.ChannelLocations{}
	unPH := UnassignedPlayerHandler{}
	unPH.SetChannelLocations(&channelLocations)
	lobbyMap := make(map[string]*LobbyHandler)
	unPH.LobbyMap = &lobbyMap
	lobbyChan := LobbyHandler{}
	lobbyChan.GlobalUnassignedPlayerHandler = &unPH
	lobbyChan.SetChannelLocations(&channelLocations)
	gameHandler := GameHandler{}
	gameHandler.SetChannelLocations(&channelLocations)
	playerChan := make(chan []byte)
	util.ChanSink([]chan []byte{playerChan})

	// allHandlers := message_handlers.MessageHandlers{}
	// allHandlers[&gameState] = struct{}{}
	// allHandlers[&lobbyChan] = struct{}{}

	lobbyChan.AddChannel(playerChan)
	util.Assert(t, lobbyChan.LobbyPlayerChannels[0] == playerChan)
	util.Assert(t, len(lobbyChan.LobbyPlayerChannels) == 1)

	lobbyChan.MoveChannel(playerChan, &gameHandler)
	util.Assert(t, len(lobbyChan.LobbyPlayerChannels) == 0)
	util.Assert(t, len(gameHandler.gameState.PlayerChannels) == 1)
	util.Assert(t, gameHandler.gameState.PlayerChannels[0] == playerChan)
}

func Test_createLobby(t *testing.T) {
	playerChan := make(chan []byte)
	util.ChanSink([]chan []byte{playerChan})
	unPH := UnassignedPlayerHandler{UnassignedPlayers: []chan []byte{playerChan}}
	msg := messages.Message{TypeDescriptor: "Create Lobby"}
	channelLocations := message_handler_interface.ChannelLocations{}
	lobbyMap := make(map[string]*LobbyHandler)
	channelLocations[playerChan] = &unPH
	// allHandlers := message_handlers.MessageHandlers{} // needed so tests compile atm, but will be removed in future
	unPH.SetChannelLocations(&channelLocations)
	unPH.LobbyMap = &lobbyMap

	unPH.ProcessUserMessage(msg, playerChan) //Stuff done here
	// util.Assert(t, len(allHandlers) == 2)
	lobby, ok := channelLocations[playerChan]
	t.Log(ok, unPH.UnassignedPlayers)
	util.Assert(t, ok)
	util.Assert(t, len(unPH.UnassignedPlayers) == 0)
	cast_lobby, ok := lobby.(*LobbyHandler)
	if !ok {
		t.FailNow()
	}
	util.Assert(t, ok)
	util.Assert(t, cast_lobby.LobbyPlayerChannels[0] == playerChan)
	util.Assert(t, len(cast_lobby.LobbyID) >= 11)
	// // Test should also check that a message is sent that contains a lobby ID corresponding to the new lobby created
	// t.Fail() // must implement the above comment
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

func Test_joiningLobbyWithHash(t *testing.T) {
	playerChan := make(chan []byte)
	unPH := UnassignedPlayerHandler{UnassignedPlayers: []chan []byte{playerChan}}
	lobby := LobbyHandler{LobbyID: "alex"}

	util.ChanSink([]chan []byte{playerChan})
	lobbyMap := make(map[string]*LobbyHandler)
	lobbyMap["alex"] = &lobby
	unPH.LobbyMap = &lobbyMap
	msg := messages.Message{TypeDescriptor: "Join Lobby", Contents: struct{ LobbyID string }{LobbyID: "alex"}}
	channelLocations := message_handler_interface.ChannelLocations{}
	unPH.SetChannelLocations(&channelLocations)
	lobby.SetChannelLocations(&channelLocations)
	channelLocations[playerChan] = &unPH
	// util.Assert(t, len(channelLocations) == 1)
	unPH.ProcessUserMessage(msg, playerChan)
	util.Assert(t, len(unPH.UnassignedPlayers) == 0)
	util.Assert(t, len(lobby.LobbyPlayerChannels) == 1)
	util.Assert(t, lobby.LobbyPlayerChannels[0] == playerChan)
	util.Assert(t, channelLocations[playerChan].(*LobbyHandler) == &lobby)

}

func Test_leavingLobby(t *testing.T) {
	playerChan := make(chan []byte)
	unPH := UnassignedPlayerHandler{}
	lobby := LobbyHandler{LobbyID: "alex", LobbyPlayerChannels: []chan []byte{playerChan}}
	lobby.GlobalUnassignedPlayerHandler = &unPH
	channelLocations := message_handler_interface.ChannelLocations{}
	channelLocations[playerChan] = &lobby
	unPH.SetChannelLocations(&channelLocations)
	lobby.SetChannelLocations(&channelLocations)
	util.ChanSink([]chan []byte{playerChan})
	lobbyMap := make(map[string]*LobbyHandler)
	lobbyMap["alex"] = &lobby
	unPH.LobbyMap = &lobbyMap
	// SETUP COMPLETED
	msg := messages.Message{TypeDescriptor: "Leave Lobby"}
	lobby.ProcessUserMessage(msg, playerChan)

	util.Assert(t, len(lobby.LobbyPlayerChannels) == 0)
	util.Assert(t, len(unPH.UnassignedPlayers) == 1)
	util.Assert(t, channelLocations[playerChan] == message_handler_interface.MessageHandler(&unPH))
}

func Test_joiningQuickplay(t *testing.T) {
	playerChan := make(chan []byte)
	unPH := UnassignedPlayerHandler{UnassignedPlayers: []chan []byte{playerChan}}
	quickplay_lobby := LobbyHandler{IsQuickplay: true}
	quickplay_lobby.GlobalUnassignedPlayerHandler = &unPH
	util.ChanSink([]chan []byte{playerChan})
	lobbyMap := make(map[string]*LobbyHandler)
	unPH.LobbyMap = &lobbyMap
	unPH.currentQuickplayLobby = &quickplay_lobby
	msg := messages.Message{TypeDescriptor: "JoinQuickplay"}
	channelLocations := message_handler_interface.ChannelLocations{}
	unPH.SetChannelLocations(&channelLocations)
	quickplay_lobby.SetChannelLocations(&channelLocations)
	channelLocations[playerChan] = &unPH
	// util.Assert(t, len(channelLocations) == 1)
	unPH.ProcessUserMessage(msg, playerChan)
	util.Assert(t, len(unPH.UnassignedPlayers) == 0)
	util.Assert(t, len(quickplay_lobby.LobbyPlayerChannels) == 1)
	util.Assert(t, quickplay_lobby.LobbyPlayerChannels[0] == playerChan)
	util.Assert(t, channelLocations[playerChan].(*LobbyHandler) == &quickplay_lobby)
}

func Test_joiningQuickplayCausingAGame(t *testing.T) {
	waitingChans := util.InitialiseChans(make([]chan []byte, 3))
	playerChan := make(chan []byte)
	unPH := UnassignedPlayerHandler{UnassignedPlayers: []chan []byte{playerChan}}
	quickplay_lobby := LobbyHandler{IsQuickplay: true, LobbyPlayerChannels: waitingChans}
	quickplay_lobby.GlobalUnassignedPlayerHandler = &unPH
	util.ChanSink(waitingChans)
	util.ChanSink([]chan []byte{playerChan})
	lobbyMap := make(map[string]*LobbyHandler)
	unPH.LobbyMap = &lobbyMap
	unPH.currentQuickplayLobby = &quickplay_lobby
	msg := messages.Message{TypeDescriptor: "JoinQuickplay"}
	channelLocations := message_handler_interface.ChannelLocations{}
	unPH.SetChannelLocations(&channelLocations)
	quickplay_lobby.SetChannelLocations(&channelLocations)
	channelLocations[playerChan] = &unPH
	for _, channel := range waitingChans {
		channelLocations[channel] = &quickplay_lobby
	}
	// FINISHED SETUP

	// util.Assert(t, len(channelLocations) == 1)
	unPH.ProcessUserMessage(msg, playerChan)
	util.Assert(t, len(unPH.UnassignedPlayers) == 0)
	util.Assert(t, len(quickplay_lobby.LobbyPlayerChannels) == 0)
	messageHandler, ok := channelLocations[playerChan]
	if !ok {
		t.Error("the player ended up without a location")
	}
	gameHandler, ok := messageHandler.(*GameHandler)
	if !ok {
		t.Error("the player ended up not belonging to a game")
	}
	util.Assert(t, len(gameHandler.gameState.PlayerChannels) == 4)
	util.Assert(t, gameHandler.gameState.GameInProgress)
	for _, channel := range waitingChans {
		gh, ok := channelLocations[channel].(*GameHandler)
		util.Assert(t, ok)
		util.Assert(t, gh == gameHandler)
		util.Assert(t, slices.Contains(gameHandler.gameState.PlayerChannels, channel))
	}
	new_quickplay_lobby := unPH.currentQuickplayLobby
	util.Assert(t, new_quickplay_lobby.IsQuickplay)
	util.Assert(t, len(new_quickplay_lobby.LobbyPlayerChannels) == 0)
}
