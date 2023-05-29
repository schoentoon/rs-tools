FROM golang as builder

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download && go install golang.org/x/tools/cmd/stringer

COPY . /app

RUN go generate ./... && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rs-cli ./cmd/rs-cli/...

FROM alpine:latest

WORKDIR /bin/

COPY --from=builder /app/rs-cli .

WORKDIR /app/

ENTRYPOINT [ "/bin/rs-cli" ]
