package message_handlers

import (
	"HigherLevelPerudoServer/message_handlers/game"
	"HigherLevelPerudoServer/message_handlers/message_handler_interface"
	"HigherLevelPerudoServer/messages"
	"encoding/json"
	"fmt"
	"slices"
)

type GameHandler struct {
	gameState                     game.GameState
	channelLocations              *message_handler_interface.ChannelLocations
	GlobalUnassignedPlayerHandler *UnassignedPlayerHandler
}

func (gameHandler *GameHandler) SetChannelLocations(channelLocations *message_handler_interface.ChannelLocations) {
	gameHandler.channelLocations = channelLocations
}

func (gameHandler *GameHandler) AddChannel(thisChan chan []byte) {
	gameHandler.gameState.PlayerChannels = append(gameHandler.gameState.PlayerChannels, thisChan)
	(*gameHandler.channelLocations)[thisChan] = gameHandler
	playerLocationMessage := messages.Message{TypeDescriptor: "PlayerLocation", Contents: "/game"}
	thisChan <- messages.CreateEncodedMessage(playerLocationMessage)
}

func (gameHandler *GameHandler) MoveChannel(thisChan chan []byte, newLocation message_handler_interface.MessageHandler) {
	// thisChanIndex := slices.Index(gameState.PlayerChannels, thisChan)
	// gameState.PlayerChannels = slices.Delete(gameState.PlayerChannels, thisChanIndex, thisChanIndex) // might need a +1 to make a valid slice
	// if len(gameState.PlayerChannels) == 0 && message_handlers.MessageHandler(gameState) != *newLocation {
	// 	delete((*allHandlers), gameState)
	// }
	// (*newLocation).AddChannel(thisChan, channelLocations)
	message_handler_interface.MoveChannelLogic(&gameHandler.gameState.PlayerChannels, thisChan, newLocation, gameHandler.channelLocations)
}

func (gameHandler *GameHandler) ProcessUserMessage(userMessage messages.Message, thisChan chan []byte) {

	// fmt.Println("Printing gamestate", gameState)
	// not very efficient. Should work
	fmt.Printf("Message recieved from websocket associated to player index %d \n", slices.Index[[]chan []byte](gameHandler.gameState.PlayerChannels, thisChan))
	switch userMessage.TypeDescriptor {
	case "PlayerMove":
		fmt.Println("Made it into PlayerMove switch")
		// If PlayerMove need to ensure that userMessage.Contents is of type PlayerMove

		// could check here to make sure that this message is coming from the current player
		// To do as such, we need a pairing from thisChan to playerIDs
		// Then check equality against gameState.CurrentPlayerIndex
		// thisChanIndex := slices.Index[[]chan []byte, chan []byte](gameState.PlayerChannels,thisChan)

		thisChanIndex := slices.Index[[]chan []byte](gameHandler.gameState.PlayerChannels, thisChan)
		if thisChanIndex != gameHandler.gameState.AllowableChannelLock {
			thisChan <- messages.PackMessage("NOT YOUR TURN", nil)
			return
		}

		var playerMove game.PlayerMove
		// the below marshal and Unmarshal could be removed with a good understanding of the type of PlayerMove being sent
		// Below may work but may also be a breaking change
		// playerMove = PlayerMove{MoveType = userMessage.Contents.MoveType, Bet userMessage.Contents.Bet}

		buff, err := json.Marshal(userMessage.Contents)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = json.Unmarshal(buff, &playerMove)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Calling gamestate.processPlayerMove")
		couldProcessMove := gameHandler.gameState.ProcessPlayerMove(playerMove)
		if !couldProcessMove {
			thisChan <- messages.PackMessage("Could not process player move", nil)
			return
		}
		fmt.Println("Finished gamestate.processPlayerMove")
		// if !valid {
		// 	gameState.PlayerChannels[gameState.CurrentPlayerIndex] <- packMessage("Invalid Bet", "Invalid Bet selection. Please select a valid move")
		// 	return
		// }
		// move was valid, broadcast new state

		// will need to let players know the result of updating the game state
	case "LeaveGame":
		gameHandler.processLeaveGame(thisChan)

	case "ReturnAllToLobby":
		fmt.Println("A player tried to return all to the lobby")
		if gameHandler.gameState.GameInProgress {
			fmt.Println("A player tried to return all players during the game")
			thisChan <- messages.PackMessage("You cannot return all players to the lobby whilst the game is in progress", nil)
			return
		}
		new_lobby := LobbyHandler{LobbyID: gameHandler.gameState.GameID, GlobalUnassignedPlayerHandler: gameHandler.GlobalUnassignedPlayerHandler}
		new_lobby.SetChannelLocations(gameHandler.channelLocations)
		for remaining := len(gameHandler.gameState.PlayerChannels); remaining > 0; remaining = len(gameHandler.gameState.PlayerChannels) {
			gameHandler.MoveChannel(gameHandler.gameState.PlayerChannels[remaining-1], &new_lobby)
		}

		(*gameHandler.GlobalUnassignedPlayerHandler.LobbyMap)[new_lobby.LobbyID] = &new_lobby
		numLobbyPlayers := len(new_lobby.LobbyPlayerChannels) // similar code in lines 85-88 of unPHandler.go, possible refactoring in future
		lobbyJoinResponse := LobbyJoinResponse{userReadableResponse: "Successfully joined lobby with given ID", LobbyID: new_lobby.LobbyID, NumPlayers: numLobbyPlayers}
		msg := messages.Message{TypeDescriptor: "Lobby Join Accepted", Contents: lobbyJoinResponse}
		new_lobby.Broadcast(msg)
	}

}

func (gameHandler *GameHandler) processLeaveGame(thisChan chan []byte) {
	fmt.Printf("Player %d tried to leave the game \n", slices.Index[[]chan []byte](gameHandler.gameState.PlayerChannels, thisChan))
	playerLocationMessage := messages.Message{TypeDescriptor: "PlayerLocation", Contents: "/"}
	if gameHandler.gameState.GameInProgress {
		// If the game is in progress, remove the relevant player
		thisChanIndex := slices.Index[[]chan []byte](gameHandler.gameState.PlayerChannels, thisChan)
		gameHandler.gameState.RemovePlayer(thisChanIndex)
		// Then once we have moved the channel, we should inform the other players that they have left
		msg := messages.Message{TypeDescriptor: "PlayerLeft", Contents: thisChanIndex}
		defer gameHandler.Broadcast(msg)
	}
	message_handler_interface.Send(thisChan, playerLocationMessage)
	gameHandler.MoveChannel(thisChan, gameHandler.GlobalUnassignedPlayerHandler)

}

func (gameHandler GameHandler) Broadcast(message messages.Message, optional_use_wait_group ...bool) {
	message_handler_interface.BroadcastLogic(gameHandler.gameState.PlayerChannels, message, optional_use_wait_group...)
}
