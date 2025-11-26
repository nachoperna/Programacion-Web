package structs

import "time"

type Movimientos struct {
	Type   string    `json:"type"`
	Amount string    `json:"amount"`
	Date   time.Time `json:"day"`
}
