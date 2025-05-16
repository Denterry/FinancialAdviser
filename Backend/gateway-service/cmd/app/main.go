package main

import (
	"log"

	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/config"
	"github.com/Denterry/FinancialAdviser/Backend/gateway-service/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
