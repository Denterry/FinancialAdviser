FROM python:3.12-slim AS builder

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

WORKDIR /app

COPY requirements.txt .
RUN --mount=type=cache,target=/root/.cache/pip \
    pip install --no-cache-dir --upgrade pip \
 && pip install --no-cache-dir -r requirements.txt

COPY . .

RUN python -m grpc_tools.protoc \
        -I src/proto \
        --python_out=src/proto/brain \
        --grpc_python_out=src/proto/brain \
        src/proto/brain/brain.proto

RUN adduser --disabled-password --gecos "" appuser

# Final Stage
FROM python:3.12-slim

WORKDIR /app

COPY --from=builder /usr/local/lib/python*/site-packages /usr/local/lib/python*/site-packages
COPY --from=builder /app /app

USER appuser

EXPOSE 9090 8080

CMD ["python", "-m", "src.main"]
