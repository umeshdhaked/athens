FROM golang:1.22.0

WORKDIR /go/src/github.com/fastbiztech/hastinapura

COPY / /go/src/github.com/fastbiztech/hastinapura/

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/Web/main.go

EXPOSE 8080

CMD ["/go/src/github.com/fastbiztech/hastinapura/main"]
