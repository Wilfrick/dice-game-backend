package main

type GameState struct {
	GameID         string
	PlayerHands    []PlayerHand
	PlayerChannels []chan []byte
	PrevBet        Bet
}
