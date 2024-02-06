package maxprice

import "testing"

func TestMaxPrice(t *testing.T) {
	t.Run("when data is perfect", func(t *testing.T) {
		input := `[
			{"name":"item1","price":100},
			{"name":"item2","price":300},
			{"name":"item3","price":200}
		]`

		want := `{"name":"item2","price":300}`
		got := maxPrice(input)

		if got != want {
			t.Errorf("got %s, but want %s", got, want)
		}
	})

	t.Run("when data is incomplete", func(t *testing.T) {
		input := `[
			{"name":"item1","price":100},
			{"name":"item2","price":300},
			{"name":"item3"}
		]`

		want := `{"name":"item2","price":300}`
		got := maxPrice(input)

		if got != want {
			t.Errorf("got %s, but want %s", got, want)
		}
	})
}
