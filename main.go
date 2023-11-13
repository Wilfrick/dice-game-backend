package main

import (
	"HigherLevelPerudoServer/game"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"golang.org/x/exp/maps"
	"golang.org/x/net/websocket"
)

func processUserMessage(userMessage game.Message, thisChan chan []byte, allChannels *map[chan []byte]int, allGames *map[int]*game.GameState) {
	gameState := (*allGames)[0]
	// fmt.Println("Printing gamestate", gameState)
	gameState.PlayerChannels = maps.Keys(*allChannels) // not very efficient. Should work
	switch userMessage.TypeDescriptor {
	case "PlayerMove":
		fmt.Println("Made it into PlayerMove switch")
		// If PlayerMove need to ensure that userMessage.Contents is of type PlayerMove

		// could check here to make sure that this message is coming from the current player
		// To do as such, we need a pairing from thisChan to playerIDs
		// Then check equality against gameState.CurrentPlayerIndex
		// thisChanIndex := slices.Index[[]chan []byte, chan []byte](gameState.PlayerChannels,thisChan)

		thisChanIndex := slices.Index(gameState.PlayerChannels, thisChan)
		if thisChanIndex != gameState.AllowableChannelLock {
			thisChan <- game.PackMessage("NOT YOUR TURN", nil)
			return
		}

		var playerMove game.PlayerMove
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
		couldProcessMove := gameState.ProcessPlayerMove(playerMove)
		if !couldProcessMove {
			thisChan <- game.PackMessage("Could not process player move", nil)
			return
		}
		fmt.Println("Finished gamestate.processPlayerMove")
		// if !valid {
		// 	gameState.PlayerChannels[gameState.CurrentPlayerIndex] <- packMessage("Invalid Bet", "Invalid Bet selection. Please select a valid move")
		// 	return
		// }
		// move was valid, broadcast new state
	case "GameStart":
		fmt.Println("Case: GameStart")
		gameState.StartNewGame()
		// will need to let players know the result of updating the game state
	}

}

func manageWsConn(ws *websocket.Conn, thisChan chan []byte, allChans *map[chan []byte]int, allGames *map[int]*game.GameState) {

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
			// fmt.Println("Sending data internally: ", buff[:n])
			externalData <- buff[:n]
		}
	}()
	for {
		select {
		case b := <-thisChan:
			_, err := ws.Write(b)
			if err != nil {
				fmt.Println(err.Error())
				delete(*allChans, thisChan)
				continue
			}
			fmt.Println("Data written out to a websocket")
		case b := <-externalData:
			fmt.Println("Data read from a websocket")
			// fmt.Println("Received data from the outside")
			// fmt.Println(string(b))
			var message game.Message
			e := json.Unmarshal(b, &message)
			if e != nil {
				fmt.Println(e.Error())
				continue
			}
			go processUserMessage(message, thisChan, allChans, allGames)

			// old but useful code, echo + broadcast, will be removed in the future
			// ws.Write(b)
			// for c := range *allChans {
			// 	if c != thisChan {
			// 		c <- b
			// 	}
			// }

			fmt.Println("Number of Connections", len(*allChans))
		}
	}
}

func main() {
	connectionChannels := make(map[chan []byte]int)
	activeGames := make(map[int]*game.GameState)
	activeGames[0] = &game.GameState{}

	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		c := make(chan []byte)

		// go func() {
		// 	c <- game.RandomPlayerHand(5).AssembleHandMessage()
		// 	time.Sleep(time.Second * 4)
		// 	c <- game.RandomPlayerHand(4).AssembleHandMessage()
		// }()
		connectionChannels[c] = 0
		manageWsConn(ws, c, &connectionChannels, &activeGames)
	}))

	// // I think we can write code down here.
	// playerHand := game.RandomPlayerHand(5)
	// // buff := make([]byte, 1024)
	// encodedPlayerHand, e := json.Marshal(playerHand)
	// if e != nil {
	// 	fmt.Println(e.Error())
	// }
	// fmt.Println(encodedPlayerHand)

	fmt.Println("Server running")

	err := http.ListenAndServe("localhost:12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}
