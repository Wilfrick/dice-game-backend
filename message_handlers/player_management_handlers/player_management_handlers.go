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
	LobbyMap              *(map[string]*LobbyHandler)
}

type JoinLobbyRequest struct {
	LobbyID string
}

type LobbyJoinResponse struct {
	userReadableResponse string // possibly remove in future, but nice to have
	LobbyID              string
	NumPlayers           int
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
		hash := "abcdefghiklmnopqrtsuvwxyz"
		// Check for collision till new // very low prob that this is the same as another hash, so we all good
		// Create lobby with hash
		newLobby := LobbyHandler{LobbyID: hash}
		newLobby.SetChannelLocations(unPH.channelLocations)
		(*unPH.LobbyMap)[hash] = &newLobby // this possibly overwrites a previous lobby with the same hash, but hopefully hashes will never be the same
		// (*allHandlers)[&newLobby] = struct{}{} // removed as we are moving away from allHandlers
		lobbyJoinResponse := LobbyJoinResponse{userReadableResponse: "Successfully joined lobby with given ID", LobbyID: hash, NumPlayers: 1}
		msg = messages.Message{TypeDescriptor: "Lobby Join Accepted", Contents: lobbyJoinResponse}
		thisChan <- messages.CreateEncodedMessage(msg)
		unPH.MoveChannel(thisChan, &newLobby)

		// thisChan <- messages.PackMessage("Lobby Created", CreatedLobby{LobbyID: hash}) // replace with Joined Lobby (or similar)
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
		lobby, ok := (*unPH.LobbyMap)[lobbyID]
		// ADd player to target lobby
		if !ok {
			fmt.Println("Player tried to join a nonexistent lobby")
			lobbyJoinResponse := LobbyJoinResponse{userReadableResponse: "Failed to join lobby with given ID", LobbyID: lobbyID}
			msg = messages.Message{TypeDescriptor: "Lobby Join Failed", Contents: lobbyJoinResponse}
			thisChan <- messages.CreateEncodedMessage(msg)
			return
		}
		unPH.MoveChannel(thisChan, lobby)
		numLobbyPlayers := len(lobby.LobbyPlayerChannels)
		lobbyJoinResponse := LobbyJoinResponse{userReadableResponse: "Successfully joined lobby with given ID", LobbyID: lobbyID, NumPlayers: numLobbyPlayers}
		msg = messages.Message{TypeDescriptor: "Lobby Join Accepted", Contents: lobbyJoinResponse}

		encoded_msg := messages.CreateEncodedMessage(msg) // this should really be lobby.Broadcast, probably using a go routine
		for _, lobbyChan := range lobby.LobbyPlayerChannels {
			lobbyChan <- encoded_msg
		}

		// thisChan <- messages.CreateEncodedMessage(msg)

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
func (unPH UnassignedPlayerHandler) Broadcast(message messages.Message, optional_use_wait_group ...bool) {
	message_handlers.BroadcastLogic(unPH.UnassignedPlayers, message, optional_use_wait_group...)
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
func (lobbyHandler LobbyHandler) Broadcast(message messages.Message, optional_use_wait_group ...bool) {
	message_handlers.BroadcastLogic(lobbyHandler.LobbyPlayerChannels, message, optional_use_wait_group...)
}
