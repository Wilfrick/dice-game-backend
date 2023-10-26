package main

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
