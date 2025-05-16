package entity

// ProviderType represents the source of the content
type ProviderType string

const (
	ProviderTwitter     ProviderType = "twitter"
	ProviderTradingView ProviderType = "tradingview"
	ProviderReddit      ProviderType = "reddit"
	ProviderRSS         ProviderType = "rss"
)

// SymbolType represents the type of financial instrument
type SymbolType string

const (
	SymbolTypeEquity    SymbolType = "equity"
	SymbolTypeCrypto    SymbolType = "crypto"
	SymbolTypeETF       SymbolType = "etf"
	SymbolTypeForex     SymbolType = "forex"
	SymbolTypeCommodity SymbolType = "commodity"
)

// CrawlStatus represents the status of a crawl job
type CrawlStatus string

const (
	CrawlStatusRunning CrawlStatus = "running"
	CrawlStatusSuccess CrawlStatus = "success"
	CrawlStatusFailed  CrawlStatus = "failed"
)
