from typing import Dict, List, Optional

SYSTEM_PROMPT = (
    "You are FS Adviser Agent, an AI-powered financial advisor. "
    "You provide well-reasoned, concise, and data-backed responses. "
    "If data is unavailable, indicate so clearly. "
    "Avoid speculation, and always include a disclaimer.\n"
)

DISCLAIMER = "\n\n_Disclaimer: This is an AI-generated response and not \
    personalized financial advice._"


def build_chat_prompt(
    user_question: str,
    context_lines: List[str],
) -> List[dict]:
    """
    Constructs a list of messages for chat-based LLM APIs (like OpenAI).

    :param user_question: The original user query
    :param context_lines: Relevant background information lines
    :return: List of chat messages
    """
    prompt = "\n".join(context_lines)

    return [
        {
            "role": "system",
            "content": SYSTEM_PROMPT,
        },
        {
            "role": "user",
            "content": f"{prompt}\n\nUser Question: {user_question}",
        },
    ]


def append_disclaimer(response: str) -> str:
    return response.strip() + DISCLAIMER


class PromptBuilder:
    @staticmethod
    def build_analysis_prompt(
        user_query: str,
        market_data: Dict[str, Dict],
        user_context: Optional[Dict] = None
    ) -> str:
        """
        Build a structured prompt for financial analysis.
        """
        prompt_parts = [
            "Based on the following information, \
                provide a detailed financial analysis:",
            "\nUser Query:",
            user_query,
            "\nMarket Data:"
        ]

        for ticker, data in market_data.items():
            prompt_parts.extend([
                f"\n{ticker}:",
                f"- Current Price: ${data.get('current_price', 'N/A')}",
                f"- 24h Change: {data.get('price_change_24h', 'N/A')}%",
                f"- Market Sentiment: {data.get('sentiment', 'N/A')}",
                f"- Price Prediction: {data.get('prediction', 'N/A')}"
            ])

        if user_context:
            prompt_parts.extend([
                "\nUser Context:",
                f"- Risk Tolerance: {
                    user_context.get('risk_tolerance', 'N/A'),
                }",
                f"- Investment Horizon: {
                    user_context.get('investment_horizon', 'N/A'),
                    }",
                f"- Preferred Assets: {
                    ', '.join(user_context.get('preferred_assets', [])),
                }"
            ])

        prompt_parts.extend([
            "\nPlease provide:",
            "1. A brief market analysis",
            "2. Specific recommendations based on the user's context",
            "3. Key risks and considerations",
            "4. Suggested next steps"
        ])

        return "\n".join(prompt_parts)

    @staticmethod
    def build_recommendation_prompt(
        user_query: str,
        market_data: Dict[str, Dict],
        user_context: Dict
    ) -> str:
        """
        Build a prompt specifically focused on investment recommendations.
        """
        prompt_parts = [
            "Based on the following information, \
                provide specific investment recommendations:",
            "\nUser Query:",
            user_query,
            "\nUser Profile:",
            f"- Risk Tolerance: {
                user_context.get('risk_tolerance', 'N/A'),
            }",
            f"- Investment Horizon: {
                user_context.get('investment_horizon', 'N/A'),
            }",
            f"- Preferred Assets: {
                ', '.join(user_context.get('preferred_assets', [])),
            }",
            "\nMarket Data:"
        ]

        for ticker, data in market_data.items():
            prompt_parts.extend([
                f"\n{ticker}:",
                f"- Current Price: ${data.get('current_price', 'N/A')}",
                f"- 24h Change: {data.get('price_change_24h', 'N/A')}%",
                f"- Market Sentiment: {data.get('sentiment', 'N/A')}",
                f"- Price Prediction: {data.get('prediction', 'N/A')}"
            ])

        prompt_parts.extend([
            "\nPlease provide:",
            "1. Specific investment recommendations",
            "2. Entry and exit points",
            "3. Position sizing suggestions",
            "4. Risk management strategies",
            "5. Monitoring recommendations"
        ])

        return "\n".join(prompt_parts)
