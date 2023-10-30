package main

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
