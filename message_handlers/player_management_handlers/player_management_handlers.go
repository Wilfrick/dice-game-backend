package player_management_handlers

import (
	"HigherLevelPerudoServer/game"
	"HigherLevelPerudoServer/message_handlers"
	"HigherLevelPerudoServer/messages"
	"encoding/json"
	"fmt"
	"math/rand"
)

type UnassignedPlayerHandler struct {
	UnassignedPlayers     []chan []byte
	currentQuickplayLobby *LobbyHandler
	channelLocations      *(message_handlers.ChannelLocations)
	lobbyMap              *(map[string]*LobbyHandler)
}

type JoinLobbyRequest struct {
	LobbyID string
}

func (unPH *UnassignedPlayerHandler) ProcessUserMessage(msg messages.Message, thisChan chan []byte) {
	fmt.Println("Unassigned player gave a message")
	// gs1 := maps.Keys(*allHandlers)[0] // will go in due course
	// gs := game.GameState(gs1)
	// gs1.AddChannel(thisChan, channelLocations) // will go in due course
	// (*channelLocations)[thisChan] = &gs
	switch msg.TypeDescriptor {

	case "Quickplay":
		// do quickplay
	case "Create Lobby":
		// Generate a random hash
		hash_chars := "abcdefghiklmnopqrtsuvwxyz"
		_ = hash_chars[rand.Intn(len(hash_chars))]
		// Check for collision till new // very low prob that this is the same as another hash, so we all good
		// Create lobby with hash
		newLobby := LobbyHandler{LobbyID: "abcdefghiklmnopqrtsuvwxyz"}
		newLobby.SetChannelLocations(unPH.channelLocations)
		// (*allHandlers)[&newLobby] = struct{}{} // removed as we are moving away from allHandlers
		unPH.MoveChannel(thisChan, &newLobby)
		// Add player to lobby
	case "Join Lobby":
		// parse message contents to extract LobbyID
		var joinRequest JoinLobbyRequest
		buff, err := json.Marshal(msg.Contents)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = json.Unmarshal(buff, &joinRequest)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		lobbyID := joinRequest.LobbyID
		// Look in contents for target lobby
		lobby, ok := (*unPH.lobbyMap)[lobbyID]
		// ADd player to target lobby
		if !ok {
			fmt.Println("Player tried to join a nonexistent lobby")
			thisChan <- messages.CreateEncodedMessage(messages.Message{TypeDescriptor: "Reject Channel Join", Contents: "No lobby with that ID found"})
			return
		}
		thisChan <- messages.CreateEncodedMessage(messages.Message{TypeDescriptor: "Accept Channel Join", Contents: "Joining lobby at given ID"})
		unPH.MoveChannel(thisChan, lobby)

	}
}

func (unPH *UnassignedPlayerHandler) AddChannel(thisChan chan []byte) {
	unPH.UnassignedPlayers = append(unPH.UnassignedPlayers, thisChan)
	(*unPH.channelLocations)[thisChan] = unPH
}

func (unPH *UnassignedPlayerHandler) MoveChannel(thisChan chan []byte, newLocation message_handlers.MessageHandler) {
	message_handlers.MoveChannelLogic(&unPH.UnassignedPlayers, thisChan, newLocation, unPH.channelLocations)
}

func (unPH *UnassignedPlayerHandler) SetChannelLocations(channelLocations *message_handlers.ChannelLocations) {
	unPH.channelLocations = channelLocations
}

type LobbyHandler struct {
	LobbyPlayerChannels []chan []byte
	LobbyID             string
	channelLocations    *message_handlers.ChannelLocations
}

func (lobbyHandler *LobbyHandler) ProcessUserMessage(msg messages.Message, thisChan chan []byte) {
	// var _ game.GameState // wtf is this
	switch msg.TypeDescriptor {
	case "Start Game":
		// Lots here mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm
		// In some order
		// Change the lobby to a game
		// Update the location of the players
		// Communicate this information. best handled here or in Game Handler

		// create a new game (including setting channellocations)
		gameState := game.GameState{}
		gameState.SetChannelLocations(lobbyHandler.channelLocations)

		// add the players that are currently in this lobby to the game
		for _, playerChan := range lobbyHandler.LobbyPlayerChannels {
			lobbyHandler.MoveChannel(playerChan, &gameState)
		}

		// give the game the correct ID
		gameState.GameID = lobbyHandler.LobbyID

		// update the channelLocations for these channels

		// done

	}
}

func (lobbyHandler *LobbyHandler) AddChannel(thisChan chan []byte) {
	lobbyHandler.LobbyPlayerChannels = append(lobbyHandler.LobbyPlayerChannels, thisChan)
	(*lobbyHandler.channelLocations)[thisChan] = lobbyHandler
}

func (lobbyHandler *LobbyHandler) MoveChannel(thisChan chan []byte, newLocation message_handlers.MessageHandler) {
	// thisChanIndex := slices.Index(lobbyHandler.LobbyPlayerChannels, thisChan)
	// lobbyHandler.LobbyPlayerChannels = slices.Delete(lobbyHandler.LobbyPlayerChannels, thisChanIndex, thisChanIndex+1) // might need a +1 to make a valid slice
	// if len(lobbyHandler.LobbyPlayerChannels) == 0 && message_handlers.MessageHandler(lobbyHandler) != *newLocation {
	// 	delete((*allHandlers), lobbyHandler)
	// }
	// (*newLocation).AddChannel(thisChan, channelLocations)
	message_handlers.MoveChannelLogic(&lobbyHandler.LobbyPlayerChannels, thisChan, newLocation, lobbyHandler.channelLocations)
}

func (lobbyHandler *LobbyHandler) SetChannelLocations(channelLocations *message_handlers.ChannelLocations) {
	lobbyHandler.channelLocations = channelLocations
}
