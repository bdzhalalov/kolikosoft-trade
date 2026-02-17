FROM golang:1.25-alpine AS builder

RUN apk --no-cache add bash make gcc musl-dev

WORKDIR /var/www/

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o ./cmd main.go

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1

FROM alpine AS runner

COPY --from=builder /var/www/cmd /
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

WORKDIR /var/www/

EXPOSE 8080

CMD ["/cmd"]