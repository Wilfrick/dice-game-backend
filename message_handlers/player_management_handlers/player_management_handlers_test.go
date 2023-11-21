package player_management_handlers

import (
	"HigherLevelPerudoServer/game"
	"HigherLevelPerudoServer/message_handlers"
	"HigherLevelPerudoServer/messages"
	"HigherLevelPerudoServer/util"
	"testing"
)

func Test_unPH(t *testing.T) {
	unPH := UnassignedPlayerHandler{}
	// lobbyChan := LobbyHandler{}
	playerChan := make(chan []byte)
	channelLocations := message_handlers.ChannelLocations{}
	allHandlers := message_handlers.MessageHandlers{}
	allHandlers[&unPH] = struct{}{}
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
	channelLocations := message_handlers.ChannelLocations{}
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
	channelLocations := message_handlers.ChannelLocations{}
	lobbyChan := LobbyHandler{}
	lobbyChan.SetChannelLocations(&channelLocations)
	gameState := game.GameState{}
	gameState.SetChannelLocations(&channelLocations)
	playerChan := make(chan []byte)

	// allHandlers := message_handlers.MessageHandlers{}
	// allHandlers[&gameState] = struct{}{}
	// allHandlers[&lobbyChan] = struct{}{}

	lobbyChan.AddChannel(playerChan)
	util.Assert(t, lobbyChan.LobbyPlayerChannels[0] == playerChan)
	util.Assert(t, len(lobbyChan.LobbyPlayerChannels) == 1)

	lobbyChan.MoveChannel(playerChan, &gameState)
	util.Assert(t, len(lobbyChan.LobbyPlayerChannels) == 0)
	util.Assert(t, len(gameState.PlayerChannels) == 1)
	util.Assert(t, gameState.PlayerChannels[0] == playerChan)
}

func Test_createLobby(t *testing.T) {
	playerChan := make(chan []byte)
	unPH := UnassignedPlayerHandler{UnassignedPlayers: []chan []byte{playerChan}}
	msg := messages.Message{TypeDescriptor: "Create Lobby"}
	channelLocations := message_handlers.ChannelLocations{}
	channelLocations[playerChan] = &unPH
	// allHandlers := message_handlers.MessageHandlers{} // needed so tests compile atm, but will be removed in future
	unPH.SetChannelLocations(&channelLocations)
	unPH.ProcessUserMessage(msg, playerChan)
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
	util.Assert(t, cast_lobby.LobbyID == "abcdefghiklmnopqrtsuvwxyz")
	// Test should also check that a message is sent that contains a lobby ID corresponding to the new lobby created
	t.Fail() // must implement the above comment
}
func Test_createGame(t *testing.T) {
	playerChan := make(chan []byte)
	lobby := LobbyHandler{LobbyPlayerChannels: []chan []byte{playerChan}, LobbyID: "alex"}
	msg := messages.Message{TypeDescriptor: "Start Game"}
	channelLocations := message_handlers.ChannelLocations{}
	channelLocations[playerChan] = &lobby
	// allHandlers := message_handlers.MessageHandlers{} // needed so tests compile atm, but will be removed in future
	lobby.SetChannelLocations(&channelLocations)
	lobby.ProcessUserMessage(msg, playerChan)
	// util.Assert(t, len(allHandlers) == 2)
	gs, ok := channelLocations[playerChan]
	t.Log(ok, lobby.LobbyPlayerChannels)
	util.Assert(t, ok)
	util.Assert(t, len(lobby.LobbyPlayerChannels) == 0)
	cast_game, ok := gs.(*game.GameState)
	if !ok {
		t.FailNow()
	}
	util.Assert(t, ok)
	util.Assert(t, cast_game.PlayerChannels[0] == playerChan)
	util.Assert(t, cast_game.GameID == "alex")
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
	channelLocations := message_handlers.ChannelLocations{}
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
