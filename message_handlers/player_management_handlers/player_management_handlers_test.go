package player_management_handlers

import (
	"HigherLevelPerudoServer/game"
	"HigherLevelPerudoServer/message_handlers"
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
	unPH.AddChannel(playerChan, &channelLocations)
	t.Log(unPH.UnassignedPlayers, playerChan)
	util.Assert(t, unPH.UnassignedPlayers[0] == playerChan)
	util.Assert(t, len(unPH.UnassignedPlayers) == 1)
	unPH.MoveChannel(playerChan, nil, &channelLocations, &allHandlers)
	t.Log(unPH.UnassignedPlayers, playerChan)
	util.Assert(t, len(unPH.UnassignedPlayers) == 0)
}

func Test_newLocation(t *testing.T) {
	unPH := UnassignedPlayerHandler{}
	lobbyChan := LobbyHandler{}
	playerChan := make(chan []byte)
	channelLocations := message_handlers.ChannelLocations{}
	allHandlers := message_handlers.MessageHandlers{}
	allHandlers[&unPH] = struct{}{}
	allHandlers[&lobbyChan] = struct{}{}
	unPH.AddChannel(playerChan, &channelLocations)
	util.Assert(t, unPH.UnassignedPlayers[0] == playerChan)
	util.Assert(t, len(unPH.UnassignedPlayers) == 1)
	unPH.MoveChannel(playerChan, &lobbyChan, &channelLocations, &allHandlers)
	util.Assert(t, len(unPH.UnassignedPlayers) == 0)
	util.Assert(t, len(lobbyChan.LobbyPlayerChannels) == 1)
	util.Assert(t, lobbyChan.LobbyPlayerChannels[0] == playerChan)
}

func Test_lobbyToGameLocation(t *testing.T) {
	lobbyChan := LobbyHandler{}
	gameState := game.GameState{}
	playerChan := make(chan []byte)
	channelLocations := message_handlers.ChannelLocations{}
	allHandlers := message_handlers.MessageHandlers{}
	allHandlers[&gameState] = struct{}{}
	allHandlers[&lobbyChan] = struct{}{}
	lobbyChan.AddChannel(playerChan, &channelLocations)
	util.Assert(t, lobbyChan.LobbyPlayerChannels[0] == playerChan)
	util.Assert(t, len(lobbyChan.LobbyPlayerChannels) == 1)
	lobbyChan.MoveChannel(playerChan, &gameState, &channelLocations, &allHandlers)
	util.Assert(t, len(lobbyChan.LobbyPlayerChannels) == 0)
	util.Assert(t, len(gameState.PlayerChannels) == 1)
	util.Assert(t, gameState.PlayerChannels[0] == playerChan)
}
