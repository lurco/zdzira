FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o zdzira ./cmd/zdzira

FROM alpine:3.21
RUN mkdir -p /data
COPY --from=builder /app/zdzira /usr/local/bin/zdzira
VOLUME /data
EXPOSE 8080
ENTRYPOINT ["zdzira", "-db", "/data/zdzira.db", "-addr", ":8080"]
