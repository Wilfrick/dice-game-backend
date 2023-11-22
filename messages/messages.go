package messages

import "encoding/json"

// {"type": "playerHand", "Contents":[4,5,3,4,5]}
type Message struct {
	TypeDescriptor string
	Contents       interface{}
}

func PackMessage(TypeDescriptor string, Contents interface{}) []byte {
	message := Message{TypeDescriptor, Contents}
	encodedMessage := CreateEncodedMessage(message)
	return encodedMessage
}

func CreateEncodedMessage(message Message) []byte {
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
