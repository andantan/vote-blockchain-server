FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/blockchain-node .

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /out/blockchain-node /app/blockchain-node
RUN mkdir -p /app/blocks

EXPOSE 9000 9001 9010 9100

CMD ["/app/blockchain-node"]
