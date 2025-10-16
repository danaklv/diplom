FROM latest AS builder

WORKDIR /app
COPY . .
RUN go mod tidy && go build -o app .

# Финальный образ
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .
CMD ["./app"]
