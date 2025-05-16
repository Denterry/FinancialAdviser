package webapi

import (
	"regexp"
	"strings"
	"time"

	"github.com/Denterry/FinancialAdviser/Backend/x-service/internal/entity"
	"github.com/google/uuid"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

var (
	// Matches $AAPL, $BTC, etc
	dollarSymbolPattern = regexp.MustCompile(`\$([A-Z]{1,5})`)
	// Matches #AAPL, #BTC, etc
	hashtagSymbolPattern = regexp.MustCompile(`#([A-Z]{1,5})`)
	// Matches standalone symbols like AAPL, BTC, etc
	standaloneSymbolPattern = regexp.MustCompile(`\b([A-Z]{1,5})\b`)
)

// extractSymbols finds financial symbols in text and hashtags
func extractSymbols(text string, hashtags []string) []string {
	// use a map to deduplicate symbols
	symbols := make(map[string]struct{})

	// extract from text
	extractFromText(text, symbols)

	// extract from hashtags
	for _, hashtag := range hashtags {
		if sym := strings.TrimPrefix(hashtag, "#"); len(sym) > 0 {
			symbols[strings.ToUpper(sym)] = struct{}{}
		}
	}

	result := make([]string, 0, len(symbols))
	for sym := range symbols {
		// symbol should be 1-5 uppercase letters
		if len(sym) >= 1 && len(sym) <= 5 && isUpperCase(sym) {
			result = append(result, sym)
		}
	}

	return result
}

// extractFromText extracts symbols from text using various patterns
func extractFromText(text string, symbols map[string]struct{}) {
	// extract $SYMBOL patterns
	matches := dollarSymbolPattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) > 1 {
			symbols[strings.ToUpper(match[1])] = struct{}{}
		}
	}

	// extract #SYMBOL patterns
	matches = hashtagSymbolPattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) > 1 {
			symbols[strings.ToUpper(match[1])] = struct{}{}
		}
	}

	// extract standalone symbols
	matches = standaloneSymbolPattern.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) > 1 {
			sym := strings.ToUpper(match[1])
			// only add if it looks like a financial symbol
			if isLikelySymbol(sym) {
				symbols[sym] = struct{}{}
			}
		}
	}
}

// isUpperCase checks if a string contains only uppercase letters
func isUpperCase(s string) bool {
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}

	return true
}

// isLikelySymbol checks if a string looks like a financial symbol
func isLikelySymbol(s string) bool {
	// common crypto prefixes
	cryptoPrefixes := []string{"BTC", "ETH", "USDT", "USDC", "BNB", "XRP", "ADA", "SOL", "DOT", "DOGE"}
	for _, prefix := range cryptoPrefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}

	// common stock suffixes
	stockSuffixes := []string{"INC", "CORP", "LTD", "CO", "HOLDINGS", "GROUP"}
	for _, suffix := range stockSuffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}

	// common ETF suffixes
	etfSuffixes := []string{"ETF", "FUND", "TRUST"}
	for _, suffix := range etfSuffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}

	return false
}

func mapScrapedTweet(t *twitterscraper.Tweet) *entity.Tweet {
	now := time.Now().UTC()

	return &entity.Tweet{
		ID:             uuid.MustParse(t.ID),
		Text:           t.Text,
		AuthorID:       t.UserID,
		UserName:       t.Username,
		CreatedAt:      t.TimeParsed,
		FetchedAt:      now,
		UpdatedAt:      now,
		Likes:          t.Likes,
		Replies:        t.Replies,
		Retweets:       t.Retweets,
		Views:          t.Views,
		URLs:           append(t.URLs, t.PermanentURL),
		Photos:         toURLs(t.Photos, func(p twitterscraper.Photo) string { return p.URL }),
		Videos:         toURLs(t.Videos, func(v twitterscraper.Video) string { return v.URL }),
		IsFinancial:    false,
		Symbols:        extractSymbols(t.Text, t.Hashtags),
		SentimentScore: 0,
		SentimentLabel: "",
	}
}

func toURLs[T any](in []T, getURL func(T) string) []string {
	out := make([]string, 0, len(in))
	for _, v := range in {
		out = append(out, getURL(v))
	}
	return out
}
