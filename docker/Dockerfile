FROM golang AS builder
WORKDIR /app
COPY . /app
RUN sex -eux; \
    cd /app; \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o drone-kuboard .

FROM alpine:latest
RUN apk update && apk add --no-cache ca-certificates
COPY --from=builder /app/drone-kuboard /bin

ENTRYPOINT ["/bin/drone-kuboard"]
