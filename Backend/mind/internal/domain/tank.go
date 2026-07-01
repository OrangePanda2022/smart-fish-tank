package domain

type Tank struct {
	TankID      string   `json:"tank_id"`
	TankName    string   `json:"tank_name"`
	TankSize    int      `json:"tank_size"`
	FishCount   int      `json:"fish_count"`
	FishStatus  string   `json:"fish_status"`
	FishSpecies []string `json:"fish_species"`
	CreatedAt   string   `json:"created_at"`
}
