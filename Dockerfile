# Build stage
FROM golang:1.23-alpine AS build

RUN http_proxy=http://127.0.0.1:8086 apk --no-cache add git

ENV GOPROXY=https://proxy.golang.org,direct

WORKDIR /app

COPY go.mod go.sum ./
RUN GOPROXY_REQUEST_TIMEOUT=30s go mod download -x

COPY . .

RUN CGO_ENABLED=0 go build -o guardian ./main.go

# Final stage
FROM alpine:3.20.3

# Install dumb-init in Alpine for signal handling
RUN apk --no-cache add dumb-init

# Copy the built binary from the build stage
COPY --from=build /app/guardian /app/guardian
COPY --from=build /app/.env.yaml /app/.env.yaml

EXPOSE 8080

ENTRYPOINT ["/usr/bin/dumb-init", "--", "/app/guardian"]

CMD ["/app/guardian"]
