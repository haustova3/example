FROM golang:1.23.4-alpine as builder

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/main.go

FROM scratch
COPY --from=builder server /bin/server

ENTRYPOINT ["/bin/server"]