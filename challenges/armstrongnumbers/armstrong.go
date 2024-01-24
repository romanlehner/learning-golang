package main

import (
	"math"
)

func isArmstrong(number int) bool {

	unit := float64(number % 10)
	tens := float64((number / 10) % 10)
	hundreds := float64((number / 100) % 10)

	sumCubes := math.Pow(unit, 3) + math.Pow(tens, 3) + math.Pow(hundreds, 3)
	
	if sumCubes == float64(number) {
		return true
	} else {
		return false
	}
}