## Builder
FROM harbor-pulomas.klik.doctor/base-os/golang:1.21.1-alpine AS builder

ENV GO111MODULE=on
ENV GOPRIVATE=gitlab.klik.doctor

RUN apk update && apk upgrade && apk add --no-cache git gcc libc-dev make openssl  g++ libc-dev librdkafka-dev pkgconf

WORKDIR /app
COPY . .

RUN cp .netrc /root
RUN chmod 600 /root/.netrc

RUN go get

RUN CGO_CFLAGS="-D_LARGEFILE64_SOURCE" CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags musl -a -installsuffix cgo -o /go/bin/marketplace-svc

## Distribution
FROM harbor-pulomas.klik.doctor/base-os/alpine:latest

ENV KD_ENV=stg
ENV TZ="Asia/Jakarta"

RUN apk update && apk upgrade && apk add --no-cache tzdata
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app
COPY --from=builder /app /app
COPY --from=builder /go/bin /go/bin

EXPOSE 5600

RUN chmod +x /go/bin/marketplace-svc

CMD ["/go/bin/marketplace-svc"]
