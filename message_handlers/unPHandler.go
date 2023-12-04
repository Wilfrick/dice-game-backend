package message_handlers

import (
	"HigherLevelPerudoServer/message_handlers/message_handler_interface"
	"HigherLevelPerudoServer/messages"
	"encoding/json"
	"fmt"
	"math/rand"
)

type UnassignedPlayerHandler struct {
	UnassignedPlayers     []chan []byte
	currentQuickplayLobby *LobbyHandler
	channelLocations      *(message_handler_interface.ChannelLocations)
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
		hash_chars := "abcdefghijklmnopqrtsuvwxyz"
		_ = hash_chars[rand.Intn(len(hash_chars))] // https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go has some good content
		hash := []byte{}
		for i := 0; i <= 10; i++ {
			hash = append(hash, hash_chars[rand.Intn(len(hash_chars))])
		}
		lobbyID := string(hash)
		// Check for collision till new // very low prob that this is the same as another hash, so we all good
		// Create lobby with hash
		newLobby := LobbyHandler{LobbyID: lobbyID, GlobalUnassignedPlayerHandler: unPH}
		newLobby.SetChannelLocations(unPH.channelLocations)
		(*unPH.LobbyMap)[lobbyID] = &newLobby // this possibly overwrites a previous lobby with the same hash, but hopefully hashes will never be the same
		// (*allHandlers)[&newLobby] = struct{}{} // removed as we are moving away from allHandlers
		lobbyJoinResponse := LobbyJoinResponse{userReadableResponse: "Successfully joined lobby with given ID", LobbyID: lobbyID, NumPlayers: 1}
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

		// encoded_msg := messages.CreateEncodedMessage(msg) // this should really be lobby.Broadcast, probably using a go routine
		lobby.Broadcast(msg)
		// for _, lobbyChan := range lobby.LobbyPlayerChannels {
		// 	lobbyChan <- encoded_msg
		// }

		// thisChan <- messages.CreateEncodedMessage(msg)

	}
}

func (unPH *UnassignedPlayerHandler) AddChannel(thisChan chan []byte) {
	unPH.UnassignedPlayers = append(unPH.UnassignedPlayers, thisChan)
	(*unPH.channelLocations)[thisChan] = unPH
	playerLocationMessage := messages.Message{TypeDescriptor: "PlayerLocation", Contents: "/"}
	thisChan <- messages.CreateEncodedMessage(playerLocationMessage)
}

func (unPH *UnassignedPlayerHandler) MoveChannel(thisChan chan []byte, newLocation message_handler_interface.MessageHandler) {
	message_handler_interface.MoveChannelLogic(&unPH.UnassignedPlayers, thisChan, newLocation, unPH.channelLocations)
}

func (unPH *UnassignedPlayerHandler) SetChannelLocations(channelLocations *message_handler_interface.ChannelLocations) {
	unPH.channelLocations = channelLocations
}
func (unPH UnassignedPlayerHandler) Broadcast(message messages.Message, optional_use_wait_group ...bool) {
	message_handler_interface.BroadcastLogic(unPH.UnassignedPlayers, message, optional_use_wait_group...)
}
func (unPH *UnassignedPlayerHandler) RemoveChannel(thisChan chan []byte) {
	message_handler_interface.RemoveChannelLogic(&unPH.UnassignedPlayers, thisChan, unPH.channelLocations)
}
