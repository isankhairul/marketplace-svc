## Builder
FROM harbor-pulomas.klik.doctor/base-image/marketplace-svc:20231220-0824 AS builder

WORKDIR /app
COPY . .
RUN cp .netrc /root
RUN chmod 600 /root/.netrc

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/klikmedis-svc

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

RUN chmod +x /go/bin/klikmedis-svc

CMD ["/go/bin/klikmedis-svc"]
