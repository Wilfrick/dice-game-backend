package main

import (
	"HigherLevelPerudoServer/game"
	"HigherLevelPerudoServer/message_handlers"
	"HigherLevelPerudoServer/message_handlers/player_management_handlers"
	"HigherLevelPerudoServer/messages"
	"HigherLevelPerudoServer/util"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

func manageWsConn(ws *websocket.Conn, thisChan chan []byte, channelLocations *message_handlers.ChannelLocations, allGames *message_handlers.MessageHandlers,
	globalUnassignedPlayersHandler *player_management_handlers.UnassignedPlayerHandler) {

	externalData := make(chan []byte)
	go func() {
		buff := make([]byte, 1024)
		for {
			fmt.Println("Waiting for data")
			n, err := ws.Read(buff)
			if err != nil {
				fmt.Println(err.Error())
				delete(*channelLocations, thisChan)
				// should also remove thisChan from allChans, so allChans should probably be a map rather than a slice
				// er.Error() == 'EOF' represents the connection closing
				return
			}
			// fmt.Println("Sending data internally: ", buff[:n])
			externalData <- buff[:n]
		}
	}()
	for {
		select {
		case b := <-thisChan:
			// fmt.Println("This channel just got", string(b))
			_, err := ws.Write(b)
			if err != nil {
				fmt.Println(err.Error())
				delete(*channelLocations, thisChan)
				continue
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
			go (handler).ProcessUserMessage(message, thisChan, channelLocations, allGames)

			// old but useful code, echo + broadcast, will be removed in the future
			// ws.Write(b)
			// for c := range *allChans {
			// 	if c != thisChan {
			// 		c <- b
			// 	}
			// }

			fmt.Println("Number of Connections", len(*channelLocations))
		}
	}
}

func main() {
	// connectionChannels := make(map[chan []byte]int)
	channelLocations := message_handlers.ChannelLocations{}
	activeGames := message_handlers.MessageHandlers{}
	activeGames[&game.GameState{}] = struct{}{}
	// activeGames["some_hash"] = &game.GameState{}
	globalUnassignedPlayersHandler := player_management_handlers.UnassignedPlayerHandler{}

	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		thisChan := make(chan []byte)
		globalUnassignedPlayersHandler.UnassignedPlayers = append(globalUnassignedPlayersHandler.UnassignedPlayers, thisChan)
		channelLocations[thisChan] = &globalUnassignedPlayersHandler

		manageWsConn(ws, thisChan, &channelLocations, &activeGames, &globalUnassignedPlayersHandler)
	}))

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
