package game

import "encoding/json"

// {"type": "playerHand", "Contents":[4,5,3,4,5]}
type Message struct {
	TypeDescriptor string
	Contents       interface{}
}

type RoundUpdateMessage struct { // lets the other players
	PrevMove PlayerMove
}

type RoundResultMessage struct {
	NewPlayerIndex int
}

func packMessage(TypeDescriptor string, Contents interface{}) []byte {
	message := Message{TypeDescriptor, Contents}
	encodedMessage := createEncodedMessage(message)
	return encodedMessage
}

func createEncodedMessage(message Message) []byte {
	encodedMessage, _ := json.Marshal(message)
	return encodedMessage
}

// {"TypeDescriptor":"PlayerMove", "Contents":{"PlayerID":234,"GameID":345,"MoveSTR":"Bet 3 five"}}

// {"TypeDescriptor":"NewLobby","Contents": {"LobbyName":"Fun game time", "GameID":876}}

// {"TypeDescriptor":"NewLobby","Contents":{"GameID":876}}

// type NewGame struct {
// 	GameId  int
// 	Players []int
// }

// type PlayerMoveMsg struct { // client telling server what move they have made
// 	PlayerID int
// 	GameID   int
// 	MoveSTR  string
// }

// type PlayerHandMsg struct { // server telling client what hand they have
// 	PlayerID int
// 	GameID   int
// 	PlayerHand
// }

// type Message struct {
// 	MessageType string
// 	Contents    string // should then be further unpacked with JSON
// }

// {"NewHand": {"PlayerID":234, "GameID":345, "Hand":[3,4,3,4,4]}}

// {"${MessageType}: ${Contents}"}
