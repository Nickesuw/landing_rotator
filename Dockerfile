FROM golang:1.17 as build

COPY ./.netrc /root/.netrc
RUN chmod 600 /root/.netrc

WORKDIR /go/src/app

RUN go env -w GOPRIVATE=gitlab.tubecorporate.com

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o  ./bin/nft-presale ./cmd/nft-presale

FROM alpine:3.10 as app
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /usr/bin
COPY --from=build /go/src/app /go
COPY ./configs configs
COPY ./assets assets
COPY --from=build /go/src/app /go
EXPOSE 8001
ENTRYPOINT /go/bin/nft-presale
