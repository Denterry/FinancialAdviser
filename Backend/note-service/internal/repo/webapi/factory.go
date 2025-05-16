package webapi

import (
	"fmt"

	"github.com/Denterry/FinancialAdviser/Backend/x-service/config"
	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/repo"
)

// NewSocialFetcher picks the implementation based on cfg.XProvider.Type
func NewSocialFetcher(cfg config.XProvider) (repo.SocialFetcher, error) {
	switch cfg.Type {
	case "api":
		return NewTwitterAPI(cfg)
	case "scraper", "":
		return NewTwitterScraper(cfg)
	default:
		return nil, fmt.Errorf("unknown X_PROVIDER_TYPE %q", cfg.Type)
	}
}
