package entity

// MessageRole represents the role of a message in a conversation
type MessageRole string

const (
	SystemRole    MessageRole = "system"
	UserRole      MessageRole = "user"
	AssistantRole MessageRole = "assistant"
)

// MessageStatus represents the processing status of a message
type MessageStatus string

const (
	MessageStatusPending   MessageStatus = "pending"
	MessageStatusComplete MessageStatus = "complete"
	MessageStatusFailed   MessageStatus = "failed"
)

// ModelProvider represents the LLM provider
type ModelProvider string

const (
	OpenAIProvider ModelProvider = "openai"
	AnthropicProvider ModelProvider = "anthropic"
	LocalProvider ModelProvider = "local"
)

// Message represents a single message in a conversation
type Message struct {
	ID        string       `json:"id"`
	Role      MessageRole  `json:"role"`
	Content   string       `json:"content"`
	Status    MessageStatus `json:"status"`
	CreatedAt int64        `json:"created_at"`
}

// Conversation represents a chat session
type Conversation struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	Title        string     `json:"title"`
	Messages     []Message  `json:"messages"`
	CreatedAt    int64      `json:"created_at"`
	LastActivity int64      `json:"last_activity"`
}

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
