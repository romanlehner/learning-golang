package main

import "testing"

func TestArmstrong(t *testing.T) {
	table := []struct {
		number      int
		isArmstrong bool
	}{
		{number: 371, isArmstrong: true},
		{number: 407, isArmstrong: true},
		{number: 200, isArmstrong: false},
	}

	for _, tt := range table {
		got := isArmstrong(tt.number)
		want := tt.isArmstrong

		if got != want {
			t.Errorf("want %v, but got %v with number %d", want, got, tt.number)
		}
	}
}
