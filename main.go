package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

func manageWsConn(ws *websocket.Conn, thisChan chan []byte, allChans *[]chan []byte) {

	externalData := make(chan []byte)
	go func() {
		buff := make([]byte, 1024)
		for {
			fmt.Println("Waiting for data")
			n, err := ws.Read(buff)
			if err != nil {
				fmt.Println(err.Error())
				// should also remove thisChan from allChans, so allChans should probably be a map rather than a slice
				return
			}
			fmt.Println("Sending data internally")
			externalData <- buff[:n]
		}
	}()
	for {
		select {
		case b := <-thisChan:
			ws.Write(b)
			fmt.Println("Free case 1")
		case b := <-externalData:
			fmt.Println("Received data from the outside")
			ws.Write(b)
			for _, c := range *allChans {
				if c != thisChan {
					c <- b
				}
			}
			fmt.Println("Free case 2")
		}
	}
}

func main() {
	connectionChannels := []chan []byte{}
	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		c := make(chan []byte)

		go func() {
			encodedPlayerHand, _ := json.Marshal(randomPlayerHand(5))
			c <- []byte(encodedPlayerHand)
			time.Sleep(time.Second * 4)
			encodedPlayerHand, _ = json.Marshal(randomPlayerHand(4))
			c <- []byte(encodedPlayerHand)
		}()
		connectionChannels = append(connectionChannels, c)
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

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}
