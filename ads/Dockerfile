FROM golang:1.14 as builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/ads ./ads/

FROM scratch 
COPY --from=builder /go/bin/ads /go/bin/ads
ENTRYPOINT ["/go/bin/ads"]
