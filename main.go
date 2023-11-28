package main

import (
	"HigherLevelPerudoServer/message_handlers"
	"HigherLevelPerudoServer/message_handlers/player_management_handlers"
	"HigherLevelPerudoServer/messages"
	"HigherLevelPerudoServer/util"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

func manageWsConn(ws *websocket.Conn, thisChan chan []byte, channelLocations *message_handlers.ChannelLocations,
	globalUnassignedPlayersHandler *player_management_handlers.UnassignedPlayerHandler) {

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
					(*channelLocations)[thisChan].MoveChannel(thisChan, nil)
					delete(*channelLocations, thisChan)
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
					(*channelLocations)[thisChan].MoveChannel(thisChan, nil)
					delete(*channelLocations, thisChan)
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
func generate_random_chars() string {
	const NUMBER_OF_CHARS = 10
	hash_chars := "abcdefghijklmnopqrtsuvwxyz"
	_ = hash_chars[rand.Intn(len(hash_chars))] // https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go has some good content
	hash := []byte{}
	for i := 0; i < NUMBER_OF_CHARS; i++ {
		hash = append(hash, hash_chars[rand.Intn(len(hash_chars))])
	}
	random_chars := string(hash)
	return random_chars
}

func generate_new_client_ID(globalClientIDToChannels map[string]chan []byte) string {
	candidateClientID := generate_random_chars()
	i := 0
	for _, ok := globalClientIDToChannels[candidateClientID]; ok && i < 1000; _, ok = globalClientIDToChannels[candidateClientID] {
		candidateClientID = generate_random_chars()
		i++
	} // we only run for max 1000 iterations. Possibly have a collision if this is exceeded
	if i == 1000 {
		fmt.Println("Likely collision")
	}
	return candidateClientID
}

func handleClientHandshake(ws *websocket.Conn, globalClientIDToChannels *map[string]chan []byte,
	globalUnassignedPlayersHandler *player_management_handlers.UnassignedPlayerHandler, channelLocations *message_handlers.ChannelLocations) chan []byte {
	// buff := make([]byte, 1024)
	_, buff, err := ws.ReadMessage()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	clientID := string(buff)
	fmt.Printf("ClientID: %s \n", clientID)
	thisChan, ok := (*globalClientIDToChannels)[clientID]
	if !ok {
		// This is a new ClientID and so we want to give the client a new ID

		clientID = generate_new_client_ID(*globalClientIDToChannels)
		thisChan = make(chan []byte)                     // Make a new channel
		(*globalClientIDToChannels)[clientID] = thisChan //Record this channel against their clientID
		globalUnassignedPlayersHandler.UnassignedPlayers = append(globalUnassignedPlayersHandler.UnassignedPlayers, thisChan)
		(*channelLocations)[thisChan] = globalUnassignedPlayersHandler
	}
	// newClientID := "AlexanderWarr"
	ws.WriteMessage(websocket.TextMessage, []byte(clientID)) // Echo back to the client their ID (new if default)
	return thisChan
}

func main() {
	// connectionChannels := make(map[chan []byte]int)
	channelLocations := message_handlers.ChannelLocations{}
	globalClientIDToChannels := make(map[string]chan []byte)
	globalLobbyMap := make(map[string]*player_management_handlers.LobbyHandler)
	globalUnassignedPlayersHandler := player_management_handlers.UnassignedPlayerHandler{}
	globalUnassignedPlayersHandler.LobbyMap = &globalLobbyMap
	globalUnassignedPlayersHandler.SetChannelLocations(&channelLocations)
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
		thisChan := handleClientHandshake(c, &globalClientIDToChannels, &globalUnassignedPlayersHandler, &channelLocations)
		manageWsConn(c, thisChan, &channelLocations, &globalUnassignedPlayersHandler)
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
