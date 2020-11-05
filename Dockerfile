FROM golang as builder

WORKDIR /app

COPY go.mod /app/

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/app ./cmd/rsge-json-server/...

FROM gcr.io/distroless/base

COPY --from=builder /bin/app /bin/app

CMD [ "/bin/app" ]

EXPOSE 8000/tcp
