package maxprice

import "encoding/json"

type Item struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func maxPrice(input string) string {
	var items []Item

	json.Unmarshal([]byte(input), &items)

	var maxPriceItem Item
	for _, item := range items {
		if item.Price > maxPriceItem.Price {
			maxPriceItem = item
		}
	}

	maxPriceItemJSON, _ := json.Marshal(maxPriceItem)

	return string(maxPriceItemJSON)
}
