package handshake

import (
	"HigherLevelPerudoServer/message_handlers"
	"HigherLevelPerudoServer/message_handlers/player_management_handlers"
	"fmt"
	"math/rand"

	"github.com/gorilla/websocket"
)

func generate_random_chars() string {
	const NUMBER_OF_CHARS = 10
	hash_chars := "abcdefghijklmnopqrtsuvwxyz"
	_ = hash_chars[rand.Intn(len(hash_chars))] // https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go has some good content
	hash := []byte{}
	for i := 0; i < NUMBER_OF_CHARS; i++ {
		hash = append(hash, hash_chars[rand.Intn(len(hash_chars))])
	}
	random_chars := string(hash)
	return random_chars
}

func generate_new_client_ID(globalClientIDToChannels map[string]chan []byte) string {
	candidateClientID := generate_random_chars()
	i := 0
	for _, ok := globalClientIDToChannels[candidateClientID]; ok && i < 1000; _, ok = globalClientIDToChannels[candidateClientID] {
		candidateClientID = generate_random_chars()
		i++
	} // we only run for max 1000 iterations. Possibly have a collision if this is exceeded
	if i == 1000 {
		fmt.Println("Likely collision")
	}
	return candidateClientID
}

func findAndInitialiseClientID(clientID string, globalClientIDToChannels *map[string]chan []byte,
	globalUnassignedPlayersHandler *player_management_handlers.UnassignedPlayerHandler, channelLocations *message_handlers.ChannelLocations) (string, chan []byte) {
	thisChan, ok := (*globalClientIDToChannels)[clientID]
	if !ok {
		// This is a new ClientID and so we want to give the client a new ID

		clientID = generate_new_client_ID(*globalClientIDToChannels)
		thisChan = make(chan []byte)                     // Make a new channel
		(*globalClientIDToChannels)[clientID] = thisChan //Record this channel against their clientID
		globalUnassignedPlayersHandler.UnassignedPlayers = append(globalUnassignedPlayersHandler.UnassignedPlayers, thisChan)
		(*channelLocations)[thisChan] = globalUnassignedPlayersHandler
	}
	return clientID, thisChan
}

func HandleClientHandshake(ws *websocket.Conn, globalClientIDToChannels *map[string]chan []byte,
	globalUnassignedPlayersHandler *player_management_handlers.UnassignedPlayerHandler, channelLocations *message_handlers.ChannelLocations) chan []byte {
	// buff := make([]byte, 1024)
	_, buff, err := ws.ReadMessage()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	clientID := string(buff)
	fmt.Printf("ClientID: %s \n", clientID)
	clientID, thisChan := findAndInitialiseClientID(clientID, globalClientIDToChannels, globalUnassignedPlayersHandler, channelLocations)
	err = ws.WriteMessage(websocket.TextMessage, []byte(clientID)) // Echo back to the client their ID (new if default)
	if err != nil {
		fmt.Println(err.Error())
		// return nil
	}
	return thisChan
}
