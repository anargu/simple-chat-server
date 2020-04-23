FROM golang:1.12-alpine AS build

RUN apk add --update --no-cache ca-certificates git

WORKDIR /simple-chat-server
COPY . ./

RUN GO111MODULE=on go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/simple-chat-server server/cmd/main.go

# This results in a single layer image
FROM scratch

# adding ca certificates
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# copying binary
COPY --from=build /bin/simple-chat-server /bin/simple-chat-server

ENTRYPOINT ["/bin/simple-chat-server"]