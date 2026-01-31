# 1. ビルド用ステージ
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 依存関係をコピーしてインストール
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド
RUN go build -o main main.go

# 2. 実行用ステージ（軽量化のため）
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .

# 実行
CMD ["./main"]