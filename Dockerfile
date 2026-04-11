# Stage 1 - Build
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api .


# Stage 2 - Run
FROM alpine:latest

WORKDIR /app

COPY --from=builder /bin/api /bin/api


EXPOSE 8080

CMD ["/bin/api"]