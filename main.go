package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

type NewGame struct {
	GameId  int
	Players []int
}

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {
	// fmt.Println("Thinking about responding")
	// time.Sleep(time.Second * 3)
	for {
		fmt.Println("Actually responding")
		fmt.Println("Hello from Jim")
		fmt.Println(("Further Debug Statement"))
		buff := make([]byte, 1024)

		// fmt.Println(string(buff))

		len_read, err := ws.Read(buff)
		if err != nil {
			return
		}
		var newGame NewGame
		// {"fname": "Alex", "age": 23}

		err = json.Unmarshal(buff, &newGame)
		if err != nil {
			fmt.Println("Couldn't unpack JSON")
		} else {
			fmt.Println("Successfully unpacked JSON")
			encoded_bytes, _ := json.Marshal(newGame)
			fmt.Println(encoded_bytes)
		}

		fmt.Println(string(buff))

		ws.Write(buff[:len_read])
		fmt.Println("Finished copying")
	}
}

// This example demonstrates a trivial echo server.
func main() {
	http.Handle("/echo", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
	fmt.Println("Keep running")
}
