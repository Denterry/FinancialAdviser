import re


_WORD_RE = re.compile(r"\w+|[^\s\w]", re.UNICODE)


def rough_token_count(text: str) -> int:
    """
    Грубая оценка количества токенов (≈ слов + знаки).
    Для биллинга этого достаточно; при желании
    можно подключить tiktoken
    """
    return len(_WORD_RE.findall(text))
