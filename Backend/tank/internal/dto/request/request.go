package request

type CreateTankRequest struct {
	TankName string `json:"tank_name"`
	TankSize int    `json:"tank_size"`
}

type UpdateTankRequest struct {
	TankName string `json:"tank_name"`
	TankSize int    `json:"tank_size"`
}

type GetTankRequest struct {
	TankID string `json:"tank_id"`
	Limit  int    `json:"limit"`
}
