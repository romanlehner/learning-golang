package sum

import "testing"

func TestSum(t *testing.T) {
	t.Run("Sums of go routines", func(t *testing.T) {
		numbers := [][]int{
			{1, 1, 1, 1, 1, 1}, //6
			{1, 1, 1},          //3
			{1, 1, 1, 1, 1},    //5
		}

		ch := make(chan int)

		for _, digits := range numbers {
			go Sum(digits, ch)
		}

		totalSum := 0
		for range numbers {
			totalSum += <-ch
		}

		got := totalSum
		want := 14
		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})
}
