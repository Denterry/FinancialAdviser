package entity

type Plan struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Features    []string `json:"features"`
}

type Status struct {
	Active    bool   `json:"active"`
	PlanID    string `json:"plan_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	AutoRenew bool   `json:"auto_renew"`
}
