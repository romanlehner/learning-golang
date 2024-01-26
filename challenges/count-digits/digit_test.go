package digit

import (
	"testing"
)

func TestDigitCount(t *testing.T) {
	t.Run("count digits", func(t *testing.T) {
		table := []struct {
			number int
			count  int
		}{
			{number: 11111, count: 5},
			{number: 11, count: 2},
			{number: 2, count: 1},
		}

		for _, tt := range table {
			got := counter(tt.number)
			want := tt.count

			if got != want {
				t.Errorf("Want %d digits, but got %d digits", want, got)
			}
		}
	})
}
