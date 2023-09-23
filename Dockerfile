FROM golang:1.19-alpine as builder
WORKDIR /root
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -o bin/app -ldflags "-s -w" main.go

FROM alpine:3.18 as app
WORKDIR /root
EXPOSE 8000
COPY --from=builder /root/bin/app app
ENTRYPOINT ["/root/app"]
