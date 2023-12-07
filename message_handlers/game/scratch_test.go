package game

// import "HigherLevelPerudoServer/game"

// // func Test_aosuenth(t *testing.T) {
// // 	var gameState GameState
// // 	gameState.CurrentPlayerIndex = 4
// // 	gameState.AllowableChannelLock = 4
// // 	func() {
// // 		defer func() { gameState.AllowableChannelLock = gameState.CurrentPlayerIndex }()
// // 		gameState.CurrentPlayerIndex = 5
// // 	}()

// // 	util.Assert(t, gameState.AllowableChannelLock == 5)
// // }

// // func Test_gun(t *testing.T) {
// // 	_ = PlayerMove{MoveType: "Rubbish"}
// // 	var moveType MoveType
// // 	moveType = "Jim"
// // 	_ = moveType
// // }

// // const (
// // 	GAMELIKE  LobbyLike
// // 	LOBBYLIKE LobbyLike
// // )

// //b := map[chan []byte]Channel_Manager
// type Channel_Manager struct {
// 	Handler      string // "game", "lobby"
// 	gameHandler  func(Message)
// 	lobbyHandler func(Message)
// 	game         *GameState
// 	lobby        *Lobby
// }

// func (cm Channel_Manager) Handle(message Message) {
// 	switch cm.Handler {
// 	case "game":
// 		cm.gameHandler(message)
// 	case "lobby":
// 		cm.lobbyHandler(message)
// 	}

// }

// type Handler interface {
// 	Handle(Message)
// }
// type GameHandler struct {
// 	gameState *GameState
// }

// func (gameHandler GameHandler) Handle(mesage Message) {
// 	// processPlayerMove
// }

// type LobbyHandler struct {

// }

// func foo2(){
// var _ map[chan[]byte]Handler
// }
// // b[thisChan].Handle(message)

// processUserMessage
// func (gameState *GameState) Handle
