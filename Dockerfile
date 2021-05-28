FROM    golang:1.16.4-alpine3.13 AS builder

WORKDIR /app

COPY    go.mod go.sum ./
RUN     go mod download

COPY    *.go ./
RUN     CGO_ENABLED=0 go build -o api .

FROM    scratch

COPY    --from=builder /app/api .

EXPOSE  8080

ENTRYPOINT ["./api"]
