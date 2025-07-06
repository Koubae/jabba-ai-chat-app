FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chat-session ./cmd/chat-session/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chat-identity ./cmd/chat-identity/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chat-orchestrator ./cmd/chat-orchestrator/main.go


FROM alpine:latest AS chat-identity

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/chat-identity .
COPY --from=builder /app/conf ./conf

EXPOSE 8080

CMD ["./chat-identity"]

FROM alpine:latest AS chat-orchestrator

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/chat-orchestrator .
COPY --from=builder /app/conf ./conf

EXPOSE 8081

CMD ["./chat-orchestrator"]

FROM alpine:latest AS chat-session

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/chat-session .
COPY --from=builder /app/conf ./conf

EXPOSE 8082

CMD ["./chat-session"]
