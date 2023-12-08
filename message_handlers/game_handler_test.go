package message_handlers

import (
	"HigherLevelPerudoServer/message_handlers/game"
	"HigherLevelPerudoServer/message_handlers/message_handler_interface"
	"HigherLevelPerudoServer/messages"
	"HigherLevelPerudoServer/util"
	"testing"
)

func Test_singlePlayerLeavesGameNotInProgress(t *testing.T) {
	gh := GameHandler{}

	channelLocations := message_handler_interface.ChannelLocations{}
	unPH := UnassignedPlayerHandler{}
	unPH.SetChannelLocations(&channelLocations)
	gh.SetChannelLocations(&channelLocations)
	gh.GlobalUnassignedPlayerHandler = &unPH
	thisChan := make(chan []byte)
	util.ChanSink([]chan []byte{thisChan})
	gh.gameState.PlayerChannels = []chan []byte{thisChan}
	gh.gameState.StartNewGame()
	gh.gameState.GameInProgress = false
	channelLocations[thisChan] = &gh
	// Finished setup
	msg := messages.Message{TypeDescriptor: "LeaveGame"}
	gh.ProcessUserMessage(msg, thisChan)
	t.Log(len(gh.gameState.PlayerChannels), len(unPH.UnassignedPlayers))
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
	gh.gameState.StartNewGame()
	gh.gameState.GameInProgress = false
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
	gh.gameState.StartNewGame()
	gh.gameState.GameInProgress = false
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

func Test_singlePlayerLeavesGameInProgressWithOtherPlayers(t *testing.T) {
	gh := GameHandler{}
	playerChan := make(chan []byte)
	util.ChanSink([]chan []byte{playerChan})
	gh.gameState.PlayerChannels = []chan []byte{playerChan}
	gh.gameState.GameInProgress = true
	gh.gameState.PlayerHands = []game.PlayerHand{game.PlayerHand{2}, game.PlayerHand{3}, game.PlayerHand{4}} //Give it a PlayerHands
	gh.gameState.InitialiseSlicesWithDefaults()

	channelLocations := message_handler_interface.ChannelLocations{}
	unPH := UnassignedPlayerHandler{}
	unPH.SetChannelLocations(&channelLocations)
	gh.SetChannelLocations(&channelLocations)
	gh.GlobalUnassignedPlayerHandler = &unPH
	gh.gameState.StartNewGame()
	gh.gameState.GameInProgress = true
	for _, thisChan := range gh.gameState.PlayerChannels {
		channelLocations[thisChan] = &gh
	}
	// Finished setup
	msg := messages.Message{TypeDescriptor: "LeaveGame"}
	gh.ProcessUserMessage(msg, playerChan)
	t.Log(len(gh.gameState.PlayerChannels))
	util.Assert(t, len(gh.gameState.PlayerChannels) == 3)
	util.Assert(t, gh.gameState.PlayerChannels[0] == nil)
	util.Assert(t, len(unPH.UnassignedPlayers) == 1)
	util.Assert(t, unPH.UnassignedPlayers[0] == playerChan)
}
func xTest_singlePlayerLeavesGameInProgressWithOtherPlayer(t *testing.T) {
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
	// util.ChanSink([]chan []byte{thisChan}) //Don't sink otherChan
	util.ChanSink([]chan []byte{thisChan, otherChan}) //Do sink otherChan
	gh.gameState.PlayerChannels = []chan []byte{thisChan, otherChan}
	gh.gameState.StartNewGame()
	gh.gameState.GameInProgress = true
	channelLocations[thisChan] = &gh
	channelLocations[otherChan] = &gh
	// Finished setup
	msg := messages.Message{TypeDescriptor: "LeaveGame"}
	// go gh.ProcessUserMessage(msg, thisChan)
	gh.ProcessUserMessage(msg, thisChan)

	// The following messages are sent in a random order which isn't ideal:
	// otherChansMessage1 := <-otherChan
	// _ = otherChansMessage1
	// otherChansMessage2 := <-otherChan
	// expectedMessage1 := messages.CreateEncodedMessage(messages.Message{TypeDescriptor: "GameResult", Contents: game.GameResult{1, "win"}})
	// expectedMessage2 := messages.CreateEncodedMessage(messages.Message{TypeDescriptor: "PlayerLeft", Contents: 0})
	// if slices.Equal(otherChansMessage1, expectedMessage1) {
	// 	util.Assert(t, slices.Equal(otherChansMessage2, expectedMessage2))
	// } else if slices.Equal(otherChansMessage1, expectedMessage2) {
	// 	util.Assert(t, slices.Equal(otherChansMessage2, expectedMessage1))
	// } else {
	// 	t.Error("failed at the randomness")
	// }
	t.Log(gh.gameState.PlayerChannels)
	util.Assert(t, len(gh.gameState.PlayerChannels) == 1)
	util.Assert(t, gh.gameState.PlayerChannels[0] == otherChan)
	util.Assert(t, len(gh.gameState.PlayerHands) == 1)

	util.Assert(t, len(unPH.UnassignedPlayers) == 1)
	util.Assert(t, unPH.UnassignedPlayers[0] == thisChan)
}

func Test_currentPlayerDisconnectsGameInProgress(t *testing.T) {
	var gh GameHandler
	gh.gameState.PlayerHands = []game.PlayerHand{{2}, {3}, {5}}
	gh.gameState.InitialiseSlicesWithDefaults()
	gh.gameState.GameInProgress = true
	gh.gameState.CurrentPlayerIndex = 0
	unPH := UnassignedPlayerHandler{}
	channelLocations := message_handler_interface.ChannelLocations{}
	unPH.channelLocations = &channelLocations
	gh.channelLocations = &channelLocations
	gh.GlobalUnassignedPlayerHandler = &unPH
	for _, thisChan := range gh.gameState.PlayerChannels {
		channelLocations[thisChan] = &gh
	}
	msg := messages.Message{TypeDescriptor: "LeaveGame"}
	gh.ProcessUserMessage(msg, gh.gameState.PlayerChannels[0])

	t.Log(len(gh.gameState.PlayerChannels))
	util.Assert(t, len(gh.gameState.PlayerChannels) == 3)
	util.Assert(t, gh.gameState.PlayerChannels[0] == nil)
	util.Assert(t, gh.gameState.CurrentPlayerIndex == 1)
}
