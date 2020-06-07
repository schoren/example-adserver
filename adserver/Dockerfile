FROM golang:1.14 as builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/adserver ./adserver/

FROM scratch 
COPY --from=builder /go/bin/adserver /go/bin/adserver
ENTRYPOINT ["/go/bin/adserver"]
