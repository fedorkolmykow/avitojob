FROM golang:alpine as builder

ENV GO111MODULE=on
WORKDIR /app/
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build ./cmd/main.go

######## Start a new stage from scratch #######
FROM alpine:latest as app

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"]