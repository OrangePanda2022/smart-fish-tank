package domain

type Action struct {
	Device string `json:"device"`
	Action string `json:"action"`
}

type Decision struct {
	StatusScore int      `json:"status_score"`
	Summary     string   `json:"summary"`
	Actions     []Action `json:"actions"`
	Reasoning   string   `json:"reasoning"`
}

type AnalyseResponse struct {
	ReportID string   `json:"report_id"`
	TankID   string   `json:"tank_id"`
	Decision Decision `json:"decision"`
}
