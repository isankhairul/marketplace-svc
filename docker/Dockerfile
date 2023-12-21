## Builder
# FROM harbor-pulomas.klik.doctor/base-image/marketplace-svc:20231220-0824 AS builder
## Builder
# FROM harbor-pulomas.klik.doctor/base-os/golang:1.21.1-alpine3.18 AS builder
FROM golang:1.21.1-alpine3.18 AS builder
ENV GO111MODULE=on
ENV GOPRIVATE=gitlab.klik.doctor

RUN apk update && apk upgrade && apk add --no-cache git gcc libc-dev make openssl  g++ libc-dev librdkafka-dev pkgconf

# WORKDIR /app
# RUN git clone https://github.com/edenhill/librdkafka.git

# WORKDIR /app/librdkafka
# RUN ls
# RUN ./configure --prefix /usr
# RUN make
# RUN make install

WORKDIR /app
COPY . .

RUN cp .netrc /root
RUN chmod 600 /root/.netrc

RUN go get

# RUN rm -rf go-swagger

RUN CGO_CFLAGS="-D_LARGEFILE64_SOURCE" CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags musl -a -installsuffix cgo -o /go/bin/marketplace-svc

## Distribution
FROM harbor-pulomas.klik.doctor/base-os/ubuntu:22.04

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
