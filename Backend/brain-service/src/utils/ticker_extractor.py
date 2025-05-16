import re
from typing import List


def extract_tickers(text: str) -> List[str]:
    """
    Extract stock ticker symbols from text.

    Args:
        text: The input text to extract tickers from

    Returns:
        List of found ticker symbols
    """
    # Common stock ticker pattern: 2-6 uppercase letters
    ticker_pattern = r'\b[A-Z]{2,6}\b'

    # Find all matches
    matches = re.findall(ticker_pattern, text)

    # Filter out common words that might match the pattern
    common_words = {
        'A', 'I', 'THE', 'IN', 'ON', 'AT', 'TO', 'FOR', 'AND', 'OR',
        'BUT', 'IF', 'OF', 'BY', 'AS', 'IS', 'ARE', 'WAS', 'WERE'
    }

    # Return unique tickers that aren't common words
    return list(set(match for match in matches if match not in common_words))
