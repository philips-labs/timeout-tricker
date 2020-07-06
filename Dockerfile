FROM golang:1.14.4 as builder
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 go build -o timeout-tricker

FROM alpine:latest 
RUN apk update && apk add ca-certificates mailcap && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /build/timeout-tricker /app
EXPOSE 8080
CMD ["/app/timeout-tricker"]
