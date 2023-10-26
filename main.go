package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"golang.org/x/net/websocket"
)

type NewGame struct {
	GameId  int
	Players []int
}

type PlayerMoveMsg struct { // client telling server what move they have made
	PlayerID int
	GameID   int
	MoveSTR  string
}

type PlayerHandMsg struct { // server telling client what hand they have
	PlayerID int
	GameID   int
	PlayerHand
}

type Message struct {
	MessageType string
	Contents    string // should then be further unpacked with JSON
}

// {"NewHand": {"PlayerID":234, "GameID":345, "Hand":[3,4,3,4,4]}}

// {"${MessageType}: ${Contents}"}

type PlayerHand []int

func randomPlayerHand(length int) PlayerHand {
	hand := make([]int, length)
	for i := range hand {
		hand[i] = rand.Intn(6) + 1
	}
	return hand
}

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {
	// fmt.Println("Thinking about responding")
	// time.Sleep(time.Second * 3)

	// "while" in go is written as "for"
	for {
		fmt.Println("Actually responding")
		fmt.Println("Hello from Jim ")
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

func SendHand(playerhand PlayerHand, ws *websocket.Conn) {
	// ws.Write()
}

func BroadcastServer(ws *websocket.Conn) {
	// Read a JSON object from the websocket

	// Unmarshal that object into a PlayerMove

	// For the future:
	// Will send the MoveSTR to all other players in the same Game.

	for {
		buff := make([]byte, 1024)
		len_read, err := ws.Read(buff)

		if err != nil {
			return
		}
		var playerMove PlayerMoveMsg
		err = json.Unmarshal(buff[:len_read], &playerMove)
		if err != nil {
			fmt.Println("Couldn't unpack JSON")
			fmt.Println(err.Error())
		} else {
			fmt.Println("Successfully unpacked JSON")
			fmt.Println(playerMove)
			// encoded_bytes, _ := json.Marshal(playerMove)
			// fmt.Println(encoded_bytes)
		}
		// Test generating and sending a random player hand
		random_hand := randomPlayerHand(5)

		encoded_bytes, _ := json.Marshal(random_hand)
		ws.Write(encoded_bytes)

		ws.Write(buff[:len_read])
		fmt.Println("Finished action response")
	}
}

// This example demonstrates a trivial echo server.
func main() {
	// http.Handle("/echo", websocket.Handler(EchoServer))
	// 	err := http.ListenAndServe(":12345", nil)
	// 	if err != nil {
	// 		panic("ListenAndServe: " + err.Error())
	// 	}
	// 	fmt.Println("Keep running")
	rndm_hand := randomPlayerHand(5)
	fmt.Println(rndm_hand)
	http.Handle("/echo", websocket.Handler(BroadcastServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
	fmt.Println("Keep running")

}
