# Build stage
FROM golang:1.23-bookworm AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o guardian ./main.go

# Final stage - minimal image
FROM alpine:latest

# Install dumb-init in Alpine for signal handling
RUN apk --no-cache add dumb-init

# Copy the built binary from the build stage
COPY --from=build /build/guardian /app/guardian

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/app/guardian"]

EXPOSE 8080
