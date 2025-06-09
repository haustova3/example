FROM golang:1.23.4-alpine as builder

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/main.go

FROM scratch
COPY --from=builder server /bin/server
COPY configs/comments-conf.yaml /bin/config/comments-conf.yaml

ENV CONFIG_FILE=/bin/config/comments-conf.yaml

ENTRYPOINT ["/bin/server"]