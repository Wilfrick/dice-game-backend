package main

import (
	"HigherLevelPerudoServer/game"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/exp/maps"
	"golang.org/x/net/websocket"
)

func processUserMessage(userMessage game.Message, channels *map[chan []byte]int) {
	var gameState game.GameState
	gameState.PlayerChannels = maps.Keys(*channels)
	switch userMessage.TypeDescriptor {
	case "PlayerMove":
		// If PlayerMove need to ensure that userMessage.Contents is of type PlayerMove
		fmt.Println("Made it into PlayerMove switch")
		var playerMove game.PlayerMove
		buff, _ := json.Marshal(userMessage.Contents)
		_ = json.Unmarshal(buff, &playerMove)
		fmt.Println("Calling gamestate.processPlayerMove")
		gameState.ProcessPlayerMove(playerMove)
		fmt.Println("Finished gamestate.processPlayerMove")
		// if !valid {
		// 	gameState.PlayerChannels[gameState.CurrentPlayerIndex] <- packMessage("Invalid Bet", "Invalid Bet selection. Please select a valid move")
		// 	return
		// }
		// move was valid, broadcast new state

		// will need to let players know the result of updating the game state
	}

}

func manageWsConn(ws *websocket.Conn, thisChan chan []byte, allChans *map[chan []byte]int) {

	externalData := make(chan []byte)
	go func() {
		buff := make([]byte, 1024)
		for {
			fmt.Println("Waiting for data")
			n, err := ws.Read(buff)
			if err != nil {
				fmt.Println(err.Error())
				delete(*allChans, thisChan)
				// should also remove thisChan from allChans, so allChans should probably be a map rather than a slice
				// er.Error() == 'EOF' represents the connection closing
				return
			}
			fmt.Println("Sending data internally: ", buff[:n])
			externalData <- buff[:n]
		}
	}()
	for {
		select {
		case b := <-thisChan:
			ws.Write(b)
			fmt.Println("Wrote ", b)
			fmt.Println("Free case 1")
		case b := <-externalData:
			fmt.Println("Received data from the outside")
			fmt.Println(string(b))
			var message game.Message
			e := json.Unmarshal(b, &message)
			if e != nil {
				fmt.Println(e.Error())
			}
			processUserMessage(message, allChans)

			// old but useful code, echo + broadcast, will be removed in the future
			ws.Write(b)
			for c := range *allChans {
				if c != thisChan {
					c <- b
				}
			}
			fmt.Println("Free case 2")
			fmt.Println("Number of Connections", len(*allChans))
		}
	}
}

func main() {
	connectionChannels := make(map[chan []byte]int)
	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		c := make(chan []byte)

		go func() {
			c <- game.RandomPlayerHand(5).AssembleHandMessage()
			time.Sleep(time.Second * 4)
			c <- game.RandomPlayerHand(4).AssembleHandMessage()
		}()
		connectionChannels[c] = 0
		manageWsConn(ws, c, &connectionChannels)
	}))
	// I think we can write code down here.
	playerHand := game.RandomPlayerHand(5)
	// buff := make([]byte, 1024)
	encodedPlayerHand, e := json.Marshal(playerHand)
	if e != nil {
		fmt.Println(e.Error())
	}
	fmt.Println(encodedPlayerHand)

	fmt.Println("Yes I can keep running")

	err := http.ListenAndServe("localhost:12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}
