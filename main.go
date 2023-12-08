package main

import (
	"HigherLevelPerudoServer/handshake"
	"HigherLevelPerudoServer/message_handlers"
	"HigherLevelPerudoServer/message_handlers/message_handler_interface"
	"HigherLevelPerudoServer/messages"
	"HigherLevelPerudoServer/util"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func tidyConnection(thisChan chan []byte, thisClientID string, channelLocations *message_handler_interface.ChannelLocations, globalClientIDsToChannel *map[string]chan []byte) {
	channelHandler, ok := (*channelLocations)[thisChan]
	if channelHandler != nil {
		channelHandler.RemoveChannel(thisChan)
	} else {
		fmt.Printf("Channel Handler was nil and ok was %t \n", ok)
	}
	// (*channelLocations)[thisChan].MoveChannel(thisChan, nil)
	delete(*channelLocations, thisChan)
	delete(*globalClientIDsToChannel, thisClientID)
}

func manageWsConn(ws *websocket.Conn, thisChan chan []byte, thisClientID string, channelLocations *message_handler_interface.ChannelLocations,
	globalUnassignedPlayersHandler *message_handlers.UnassignedPlayerHandler, globalClientIDsToChannel *map[string]chan []byte) {

	externalData := make(chan []byte)
	go func() {
		// buff := make([]byte, 1024)
		for {
			fmt.Println("Waiting for data")
			message_type, buff, err := ws.ReadMessage()
			// fmt.Println("Message Type", message_type)
			_ = message_type
			if err != nil {
				fmt.Printf("Websocket closed with error: %s \n", err.Error())
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					fmt.Println("There was an abnormal closure")
					// We should not delete the channel and this players stuff

				} else {
					tidyConnection(thisChan, thisClientID, channelLocations, globalClientIDsToChannel)

				}
				// (*channelLocations)[thisChan].MoveChannel(thisChan, nil)
				// delete(*channelLocations, thisChan)
				// should also remove thisChan from allChans, so allChans should probably be a map rather than a slice
				// er.Error() == 'EOF' represents the connection closing
				return
			}
			// fmt.Println("Sending data internally: ", buff[:n])
			externalData <- buff
		}
	}()
	for {
		select {
		case b := <-thisChan:
			// fmt.Println("This channel just got", string(b))
			err := ws.WriteMessage(websocket.TextMessage, b)

			if err != nil {
				fmt.Printf("Websocket couldn't write with error: %s \n", err.Error())
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					fmt.Println("There was an abnormal closure")
					// We should not delete the channel and this players stuff

				} else {
					tidyConnection(thisChan, thisClientID, channelLocations, globalClientIDsToChannel)
				}
				continue // very questionable. Should probably return
			}
			fmt.Println("Data written out to a websocket")
		case b := <-externalData:
			fmt.Println("Data read from a websocket")
			// fmt.Println("Received data from the outside")
			// fmt.Println(string(b))
			var message messages.Message
			e := json.Unmarshal(b, &message)
			if e != nil {
				fmt.Println(e.Error())
				continue
			}
			if _, ok := ((*channelLocations)[thisChan]); !ok {
				fmt.Println("The current channel does not have an associated location!")
				// Should be impossible
				(*channelLocations)[thisChan] = globalUnassignedPlayersHandler
			}
			handler := ((*channelLocations)[thisChan])
			go (handler).ProcessUserMessage(message, thisChan)

			// old but useful code, echo + broadcast, will be removed in the future
			// ws.Write(b)
			// for c := range *allChans {
			// 	if c != thisChan {
			// 		c <- b
			// 	}
			// }

			fmt.Println("Number of Connections", len(*channelLocations), "with ", len(globalUnassignedPlayersHandler.UnassignedPlayers), "players unassigned")
		}
	}
}

func main() {
	// connectionChannels := make(map[chan []byte]int)
	channelLocations := message_handler_interface.ChannelLocations{}
	globalClientIDToChannels := make(map[string]chan []byte)
	globalLobbyMap := make(map[string]*message_handlers.LobbyHandler)
	globalUnassignedPlayersHandler := message_handlers.UnassignedPlayerHandler{}
	globalUnassignedPlayersHandler.LobbyMap = &globalLobbyMap
	globalUnassignedPlayersHandler.SetChannelLocations(&channelLocations)
	globalUnassignedPlayersHandler.CreateNewQuickPlay()
	// activeHandlers[&globalUnassignedPlayersHandler] = struct{}{}
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // less than Zero CORS security.
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("Upgrade:", err)
			return
		}
		defer c.Close()
		thisChan, clientID := handshake.HandleClientHandshake(c, &globalClientIDToChannels, &globalUnassignedPlayersHandler, &channelLocations)
		manageWsConn(c, thisChan, clientID, &channelLocations, &globalUnassignedPlayersHandler, &globalClientIDToChannels)
	})

	// http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
	// 	thisChan := make(chan []byte)
	// 	globalUnassignedPlayersHandler.UnassignedPlayers = append(globalUnassignedPlayersHandler.UnassignedPlayers, thisChan)
	// 	channelLocations[thisChan] = &globalUnassignedPlayersHandler

	// 	manageWsConn(ws, thisChan, &channelLocations, &globalUnassignedPlayersHandler)
	// }))

	// // I think we can write code down here.
	// playerHand := game.RandomPlayerHand(5)
	// // buff := make([]byte, 1024)
	// encodedPlayerHand, e := json.Marshal(playerHand)
	// if e != nil {
	// 	fmt.Println(e.Error())
	// }
	// fmt.Println(encodedPlayerHand)
	path := "config.json"
	config, err := util.ReadConfigFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// fmt.Println(config)
	address := fmt.Sprintf("%s:%d", config.WebsocketHost, config.WebsocketPort)
	fmt.Println("Server running")

	err = http.ListenAndServe(address, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
