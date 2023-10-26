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

type PlayerMove struct {
	PlayerID int
	GameID   int
	MoveSTR  string
}

type SimpleObject struct {
	Name             string
	Age              int
	Favourite_Colour string
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

func BroadcastServer(ws *websocket.Conn) {
	// Read a JSON object from the websocket

	// Unmarshal that object into a PlayerMove

	// For the future:
	// Will send the MoveSTR to all other players in the same Game.

	for {
		fmt.Println("In Broadcast Server")
		buff := make([]byte, 1024)
		len_read, err := ws.Read(buff)
		fmt.Println(buff[:len_read])
		fmt.Println(string(buff))
		if err != nil {
			return
		}
		var playerMove PlayerMove
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

		ws.Write(buff[:len_read])
		fmt.Println("Finished copying")
	}
}

func exampleJSON() {
	// Take this string and print it out
	stringRepresentationOfJsonObject := "{\"Name\":\"Alex\",\"Age\": 23,\"Favourite_Colour\": \"Blue\"}"
	fmt.Println(stringRepresentationOfJsonObject)
	// create a []byte array from the string that looks like a json object
	// TODO
	jsonObjectPreUnmarshal := []byte(stringRepresentationOfJsonObject)
	fmt.Println(jsonObjectPreUnmarshal)
	// Unmarshal that object into a struct
	// TODO a = ["\x67", "\x83", "\x39"]
	var simpleobject SimpleObject
	err := json.Unmarshal(jsonObjectPreUnmarshal, &simpleobject)
	if err != nil {
		fmt.Println("Failed to Unmarshal")
		return
	}
	fmt.Println("Unmarshalled")
	fmt.Println(simpleobject)
	// Marshal the populated struct back into a []bytes
	// TODO
	byteMarshal, _ := json.Marshal(simpleobject)

	// Print the reconstructed []bytes (by using string())
	// TODO
	fmt.Println(string(byteMarshal))
}

// This example demonstrates a trivial echo server.
func main() {
	// http.Handle("/echo", websocket.Handler(EchoServer))
	// 	err := http.ListenAndServe(":12345", nil)
	// 	if err != nil {
	// 		panic("ListenAndServe: " + err.Error())
	// 	}
	// 	fmt.Println("Keep running")
	exampleJSON()
	http.Handle("/echo", websocket.Handler(BroadcastServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
	fmt.Println("Keep running")

}
