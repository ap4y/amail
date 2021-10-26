FROM docker.io/library/golang:alpine AS builder

ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -tags netgo -ldflags '-extldflags "-static" -s -w' -o /cloud-mail ./cmd/cloud-mail

FROM docker.io/library/alpine

RUN apk add --no-cache "notmuch" "w3m"
COPY --from=builder /cloud-mail /cloud-mail

EXPOSE 8000

ENTRYPOINT ["/cloud-mail"]
