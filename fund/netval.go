package fund

import (
	"encoding/json"
	"fmt"
)

// Netval is a price of a day
type Netval struct {
	Price Price `json:"price"`
	Date  Date  `json:"date"`
}

// MarshalJSON converts Netval event to json
func (n Netval) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Price string
		Date  string
	}{
		Price: fmt.Sprintf("%s", n.Price.String()),
		Date:  fmt.Sprintf("%s", n.Date.String()),
	})
}
