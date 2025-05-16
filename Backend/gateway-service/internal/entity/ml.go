package entity

type Analysis struct {
	Score    float64  `json:"score"`
	Insights []string `json:"insights"`
	Risk     string   `json:"risk"`
}

type Recommendation struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
}
