# ---- build stage ----
FROM golang:1.25-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Use modules; copy go.* first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# build static binary without CGO
ENV CGO_ENABLED=0
RUN GOOS=linux GOARCH=amd64 \
    go build -o CouponAPI .

# ---- runtime stage ----
FROM alpine:3.22

# CA needed for HTTP requests
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/CouponAPI .
COPY --from=builder /app/.env .

EXPOSE 1323

CMD ["./CouponAPI"]
