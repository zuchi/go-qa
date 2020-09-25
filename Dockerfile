## Compile Image
FROM golang:1.15.2-alpine3.12 as BUILDER
RUN apk update && apk add --no-cache curl ca-certificates git
RUN mkdir /go/qa
COPY . /go/qa/
WORKDIR /go/qa
RUN env GOOS=linux GOARCH=amd64 go build -o dist/qa ./cmd/api.go


### Runner image
FROM alpine:3.12 as RUNNER
RUN apk update && apk add --update bash ca-certificates curl
RUN apk add -U tzdata && cp /usr/share/zoneinfo/America/Sao_Paulo /etc/localtime && echo "America/Sao_Paulo" /etc/timezone
COPY --from=builder /go/qa/dist/qa /qa
RUN ["chmod", "+x", "qa"]
ENTRYPOINT ["./qa"]