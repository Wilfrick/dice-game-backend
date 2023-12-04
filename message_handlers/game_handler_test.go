package message_handlers

import (
	"HigherLevelPerudoServer/message_handlers/game"
	"HigherLevelPerudoServer/message_handlers/message_handler_interface"
	"HigherLevelPerudoServer/messages"
	"HigherLevelPerudoServer/util"
	"slices"
	"testing"
)

func Test_singlePlayerLeavesGameNotInProgress(t *testing.T) {
	gh := GameHandler{}
	gh.gameState.GameInProgress = false
	channelLocations := message_handler_interface.ChannelLocations{}
	unPH := UnassignedPlayerHandler{}
	unPH.SetChannelLocations(&channelLocations)
	gh.SetChannelLocations(&channelLocations)
	gh.GlobalUnassignedPlayerHandler = &unPH
	thisChan := make(chan []byte)
	util.ChanSink([]chan []byte{thisChan})
	gh.gameState.PlayerChannels = []chan []byte{thisChan}
	channelLocations[thisChan] = &gh
	// Finished setup
	msg := messages.Message{TypeDescriptor: "LeaveGame"}
	gh.ProcessUserMessage(msg, thisChan)
	util.Assert(t, len(gh.gameState.PlayerChannels) == 0)
	util.Assert(t, len(unPH.UnassignedPlayers) == 1)
	util.Assert(t, unPH.UnassignedPlayers[0] == thisChan)
}

func Test_singlePlayerLeavesGameWithOtherPlayerInNotInProgress(t *testing.T) {
	gh := GameHandler{}
	gh.gameState.GameInProgress = false
	channelLocations := message_handler_interface.ChannelLocations{}
	unPH := UnassignedPlayerHandler{}
	unPH.SetChannelLocations(&channelLocations)
	gh.SetChannelLocations(&channelLocations)
	gh.GlobalUnassignedPlayerHandler = &unPH
	thisChan := make(chan []byte)
	otherChan := make(chan []byte)
	util.ChanSink([]chan []byte{thisChan, otherChan})
	gh.gameState.PlayerChannels = []chan []byte{thisChan, otherChan}
	channelLocations[thisChan] = &gh
	channelLocations[otherChan] = &gh
	// Finished setup
	t.Log(gh.gameState.PlayerChannels)
	msg := messages.Message{TypeDescriptor: "LeaveGame"}
	gh.ProcessUserMessage(msg, thisChan) //thisChan leaves
	util.Assert(t, len(gh.gameState.PlayerChannels) == 1)
	util.Assert(t, len(unPH.UnassignedPlayers) == 1)
	util.Assert(t, unPH.UnassignedPlayers[0] == thisChan)
}

func Test_singlePlayerMovesMultiplePlayersToLobbyGameFinished(t *testing.T) {
	gh := GameHandler{}
	gh.gameState.GameInProgress = false //not required but nice to be explicit
	GAMEID := "GAMEID"
	gh.gameState.GameID = GAMEID
	channelLocations := message_handler_interface.ChannelLocations{}
	unPH := UnassignedPlayerHandler{}
	lobbyMap := make(map[string]*LobbyHandler)
	unPH.LobbyMap = &lobbyMap
	unPH.SetChannelLocations(&channelLocations)
	gh.SetChannelLocations(&channelLocations)
	gh.GlobalUnassignedPlayerHandler = &unPH
	thisChan := make(chan []byte)
	otherChan := make(chan []byte)
	util.ChanSink([]chan []byte{thisChan, otherChan})
	gh.gameState.PlayerChannels = []chan []byte{thisChan, otherChan}
	channelLocations[thisChan] = &gh
	channelLocations[otherChan] = &gh
	// Finished setup
	msg := messages.Message{TypeDescriptor: "ReturnAllToLobby"}
	gh.ProcessUserMessage(msg, thisChan) //thisChan leaves
	util.Assert(t, len(gh.gameState.PlayerChannels) == 0)
	util.Assert(t, len(unPH.UnassignedPlayers) == 0)
	util.Assert(t, len(lobbyMap) == 1)
	lobby, ok := lobbyMap[GAMEID]
	if !ok {
		t.Error("Did not create lobby with the correct ID")
	}
	util.Assert(t, channelLocations[thisChan] == lobby)
	util.Assert(t, channelLocations[otherChan] == lobby)
	util.Assert(t, lobby.LobbyID == GAMEID)
}

func Test_singlePlayerMovesMultiplePlayersToLobbyGameInProgress(t *testing.T) {
	gh := GameHandler{}
	gh.gameState.GameInProgress = true //not required but nice to be explicit
	GAMEID := "GAMEID"
	gh.gameState.GameID = GAMEID
	channelLocations := message_handler_interface.ChannelLocations{}
	unPH := UnassignedPlayerHandler{}
	lobbyMap := make(map[string]*LobbyHandler)
	unPH.LobbyMap = &lobbyMap
	unPH.SetChannelLocations(&channelLocations)
	gh.SetChannelLocations(&channelLocations)
	gh.GlobalUnassignedPlayerHandler = &unPH
	thisChan := make(chan []byte)
	otherChan := make(chan []byte)
	util.ChanSink([]chan []byte{thisChan, otherChan})
	gh.gameState.PlayerChannels = []chan []byte{thisChan, otherChan}
	channelLocations[thisChan] = &gh
	channelLocations[otherChan] = &gh
	// Finished setup
	msg := messages.Message{TypeDescriptor: "ReturnAllToLobby"}
	gh.ProcessUserMessage(msg, thisChan) //thisChan leaves
	util.Assert(t, len(gh.gameState.PlayerChannels) == 2)
	util.Assert(t, len(unPH.UnassignedPlayers) == 0)
	util.Assert(t, len(lobbyMap) == 0)
	// lobby, ok := lobbyMap[GAMEID]
	// if !ok {
	// 	t.Error("Did not create lobby with the correct ID")
	// }
	// util.Assert(t, channelLocations[thisChan] == lobby)
	// util.Assert(t, channelLocations[otherChan] == lobby)
	// util.Assert(t, lobby.LobbyID == GAMEID)
}

func Test_singlePlayerLeavesGameInProgress(t *testing.T) {
	gh := GameHandler{}
	gh.gameState.GameInProgress = true
	gh.gameState.PlayerHands = []game.PlayerHand{game.PlayerHand([]int{2, 2, 2})} //Give it a PlayerHands
	channelLocations := message_handler_interface.ChannelLocations{}
	unPH := UnassignedPlayerHandler{}
	unPH.SetChannelLocations(&channelLocations)
	gh.SetChannelLocations(&channelLocations)
	gh.GlobalUnassignedPlayerHandler = &unPH
	thisChan := make(chan []byte)
	util.ChanSink([]chan []byte{thisChan})
	gh.gameState.PlayerChannels = []chan []byte{thisChan}
	channelLocations[thisChan] = &gh
	// Finished setup
	msg := messages.Message{TypeDescriptor: "LeaveGame"}
	gh.ProcessUserMessage(msg, thisChan)
	util.Assert(t, len(gh.gameState.PlayerChannels) == 0)
	util.Assert(t, len(gh.gameState.PlayerChannels) == 0)
	util.Assert(t, len(unPH.UnassignedPlayers) == 1)
	util.Assert(t, unPH.UnassignedPlayers[0] == thisChan)
}
func Test_singlePlayerLeavesGameInProgressWithOtherPlayer(t *testing.T) {
	gh := GameHandler{}
	gh.gameState.GameInProgress = true
	gh.gameState.PlayerHands = []game.PlayerHand{game.PlayerHand([]int{2, 2, 2}), game.PlayerHand([]int{3, 3, 3})} //Give it a PlayerHands
	channelLocations := message_handler_interface.ChannelLocations{}
	unPH := UnassignedPlayerHandler{}
	unPH.SetChannelLocations(&channelLocations)
	gh.SetChannelLocations(&channelLocations)
	gh.GlobalUnassignedPlayerHandler = &unPH
	thisChan := make(chan []byte)
	otherChan := make(chan []byte)
	util.ChanSink([]chan []byte{thisChan}) //Don't sink otherChan
	gh.gameState.PlayerChannels = []chan []byte{thisChan, otherChan}
	channelLocations[thisChan] = &gh
	channelLocations[otherChan] = &gh
	// Finished setup
	msg := messages.Message{TypeDescriptor: "LeaveGame"}
	go gh.ProcessUserMessage(msg, thisChan)
	otherChansMessage := <-otherChan
	target_msg := messages.CreateEncodedMessage(messages.Message{TypeDescriptor: "PlayerLeft", Contents: 0})
	// fmt.Println(otherChansMessage)
	// fmt.Println(target_msg)
	util.Assert(t, slices.Equal(otherChansMessage, target_msg))
	util.Assert(t, len(gh.gameState.PlayerChannels) == 1)
	util.Assert(t, gh.gameState.PlayerChannels[0] == otherChan)
	util.Assert(t, len(gh.gameState.PlayerHands) == 1)

	util.Assert(t, len(unPH.UnassignedPlayers) == 1)
	util.Assert(t, unPH.UnassignedPlayers[0] == thisChan)
}
