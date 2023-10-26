package main

import (
	"fmt"
	"net/http"

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
		connectionChannels = append(connectionChannels, c)
		manageWsConn(ws, c, &connectionChannels)
	}))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}
