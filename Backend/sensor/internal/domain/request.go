package domain

type Request struct {
	TankID string `json:"tank_id"`
	Limit  int    `json:"limit"`
	Type   string `json:"type"`
}
