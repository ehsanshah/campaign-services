# =========================================================
# Stage 1: Builder
# =========================================================
# 🔴 تغییر مهم: استفاده از نسخه 1.24
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# بیلد کردن
RUN CGO_ENABLED=0 GOOS=linux go build -o campaign-server src/cmd/server/main.go

# =========================================================
# Stage 2: Runner
# =========================================================
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/campaign-server .

EXPOSE 50054

CMD ["./campaign-server"]