package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

func processUserMessage(userMessage Message) {
	switch userMessage.Contents {
	case "PlayerMove":

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
			var message Message
			e := json.Unmarshal(b, &message)
			if e != nil {
				fmt.Println(e.Error())
			}
			processUserMessage(message)

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
			c <- randomPlayerHand(5).assembleHandMessage()
			time.Sleep(time.Second * 4)
			c <- randomPlayerHand(4).assembleHandMessage()
		}()
		connectionChannels[c] = 0
		manageWsConn(ws, c, &connectionChannels)
	}))
	// I think we can write code down here.
	playerHand := randomPlayerHand(5)
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
