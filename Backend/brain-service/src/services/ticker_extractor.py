import re
from typing import List

# Hardcoded map (in real system, load from DB or shared file)
TICKER_MAP = {
    "tesla": "TSLA",
    "apple": "AAPL",
    "amazon": "AMZN",
    "microsoft": "MSFT",
    "meta": "META",
    "google": "GOOGL",
    "alphabet": "GOOGL",
    "nvidia": "NVDA",
    "netflix": "NFLX",
    "bitcoin": "BTC",
    "ethereum": "ETH",
}


def extract_tickers_from_text(text: str) -> List[str]:
    """
    Extracts potential financial tickers from user input.
    Supports $TSLA format and fuzzy matches like "Tesla" → "TSLA".
    """
    tickers = set()

    # Regex for $TSLA, $AAPL, etc.
    matches = re.findall(r"\$[A-Z]{1,5}", text.upper())
    tickers.update([m[1:] for m in matches])  # Strip the dollar sign

    # Fuzzy match from known names
    lowered = text.lower()
    for name, ticker in TICKER_MAP.items():
        if name in lowered:
            tickers.add(ticker)

    return sorted(tickers)
