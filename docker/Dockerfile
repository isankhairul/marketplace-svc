## Builder
FROM golang:1.21.1-alpine AS builder

ENV GO111MODULE=on
ENV GOPRIVATE=gitlab.klik.doctor

RUN apk update && apk upgrade && apk add --no-cache git gcc libc-dev make openssl  g++ libc-dev librdkafka-dev pkgconf

WORKDIR /app
COPY . .

RUN cp .netrc /root
RUN chmod 600 /root/.netrc

RUN go get

#RUN git clone https://github.com/go-swagger/go-swagger && cd go-swagger && go install ./cmd/swagger
#RUN swagger generate spec -o swagger.yaml --scan-models
#RUN sed -i 's|localhost:5600|boilerplate-api.medkomtek-stg.com|g' swagger.yaml && sed -i 's/\- https//g' swagger.yaml && sed -i "s|- http|- https|g" swagger.yaml && rm -rf go-swagger

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/klikmedis-svc

## Distribution
FROM alpine:latest

ENV KD_ENV=stg
ENV TZ="Asia/Jakarta"

RUN apk update && apk upgrade && apk add --no-cache tzdata
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app
COPY --from=builder /app /app
COPY --from=builder /go/bin /go/bin

EXPOSE 5600

RUN chmod +x /go/bin/klikmedis-svc

CMD ["/go/bin/klikmedis-svc"]
