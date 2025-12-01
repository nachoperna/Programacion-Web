package structs

type Movimientos struct {
	Type   string `json:"type"`
	Amount string `json:"amount"`
	Day    string `json:"day"`
	Time   string `json:"time"`
}
