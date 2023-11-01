package game

import "testing"

func Test_isGreaterThanDoesntModifyParameters(t *testing.T) {
	b1, b2 := Bet{5, 5}, Bet{1, 1}
	b1.isGreaterThan(b2)
	if b2.FaceVal != 1 {
		t.Fail()
	}
	if b2.NumDice != 1 {
		t.Fail()
	}
	b2.isGreaterThan(b1)
	if b2.FaceVal != 1 {
		t.Fail()
	}
	if b2.NumDice != 1 {
		t.Fail()
	}
}

func Test_firstBetComparisons(t *testing.T) {
	b1, b2, b3, b4, b5, b6, b7, b8, b9, b10, b11, b12 := Bet{1, 2}, Bet{1, 3}, Bet{1, 4}, Bet{1, 5}, Bet{1, 6}, Bet{2, 2}, Bet{2, 3}, Bet{2, 4}, Bet{2, 5}, Bet{2, 6}, Bet{1, 1}, Bet{3, 2}
	if !b12.isGreaterThan(b11) {
		t.Fail()
	}
	if !b11.isGreaterThan(b10) {
		t.Fail()
	}
	if !b10.isGreaterThan(b9) {
		t.Fail()
	}
	if !b9.isGreaterThan(b8) {
		t.Fail()
	}
	if !b8.isGreaterThan(b7) {
		t.Fail()
	}
	if !b7.isGreaterThan(b6) {
		t.Fail()
	}
	if !b6.isGreaterThan(b5) {
		t.Fail()
	}
	if !b5.isGreaterThan(b4) {
		t.Fail()
	}
	if !b4.isGreaterThan(b3) {
		t.Fail()
	}
	if !b3.isGreaterThan(b2) {
		t.Fail()
	}
	if !b2.isGreaterThan(b1) {
		t.Fail()
	}

}
