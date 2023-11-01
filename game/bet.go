package game

type Bet struct {
	NumDice int
	FaceVal int
}

func (bet1 Bet) isGreaterThan(bet2 Bet) bool {
	//1's special
	//Otherwise 2<3<...<6

	// The following modifies a copy of bet1, bet2
	// So does not modify external state
	if bet1.FaceVal == 1 {
		bet1.FaceVal = 7
		bet1.NumDice = bet1.NumDice * 2
	}
	if bet2.FaceVal == 1 {
		bet2.FaceVal = 7
		bet2.NumDice = bet2.NumDice * 2
	}

	if bet1.NumDice != bet2.NumDice {
		return bet1.NumDice > bet2.NumDice
	}
	return bet1.FaceVal > bet2.FaceVal

}
