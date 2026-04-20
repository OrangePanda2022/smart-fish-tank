package domain

type Tank struct {
	TankID      string   `json:"tank_id"`
	FishSpecies []string `json:"fish_species"`
	CreatedAt   string   `json:"created_at"`
}
