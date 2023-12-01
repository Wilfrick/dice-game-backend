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
		fmt.Printf("Player %d tried to leave the game \n", slices.Index[[]chan []byte](gameHandler.gameState.PlayerChannels, thisChan))
		playerLocationMessage := messages.Message{TypeDescriptor: "PlayerLocation", Contents: "/"}
		message_handler_interface.Send(thisChan, playerLocationMessage)
		gameHandler.MoveChannel(thisChan, gameHandler.GlobalUnassignedPlayerHandler)
	case "ReturnAllToLobby":
		fmt.Println("A player tried to return all to the lobby")
		if gameHandler.gameState.GameInProgress {
			fmt.Println("A player tried to return all players during the game")
			thisChan <- messages.PackMessage("You cannot return all players to the lobby whilst the game is in progress", nil)
			return
		}
		// new_lobby := player_management_handlers.LobbyHandler{LobbyID: gameState.GameID}
		// gameState
		returnToLobbyMessage := messages.Message{TypeDescriptor: "PlayerLocation", Contents: "/lobby"}
		gameHandler.gameState.Broadcast(returnToLobbyMessage)
	}

}

func (gameHandler GameHandler) Broadcast(message messages.Message, optional_use_wait_group ...bool) {
	message_handler_interface.BroadcastLogic(gameHandler.gameState.PlayerChannels, message, optional_use_wait_group...)
}
