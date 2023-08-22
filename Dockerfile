
FROM golang:1.21-alpine as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /usr/bin/server .

FROM gcr.io/distroless/static-debian11

COPY --from=builder /usr/bin/ /app/bin/
COPY --from=builder /usr/src/app/migrations /app/migrations

ENTRYPOINT [ "/app/bin/server" ]