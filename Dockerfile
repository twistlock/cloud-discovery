FROM golang:1.11.1
WORKDIR /go/src/github.com/twistlock/cloud-discovery/
COPY . .
RUN go fmt ./...
RUN go vet ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/server/main.go

FROM alpine:3.8
RUN apk --no-cache add ca-certificates nmap
WORKDIR /root/
COPY --from=0 /go/src/github.com/twistlock/cloud-discovery/app .
CMD ["./app"]