# syntax=docker/dockerfile:1.2
FROM golang:1.21.0-alpine3.17

RUN go version
ENV GOPATH=/

COPY ./ /binance-converter-backend

WORKDIR /binance-converter-backend

# build go app
RUN go mod download
RUN go mod tidy
RUN go build -o binance-converter-backend ./cmd/backend-server/main.go

FROM alpine:3.17
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=0 /binance-converter-backend/binance-converter-backend .
COPY config/.env .

#COPY entrypoint.sh .
#
#RUN chmod 777 /root/entrypoint.sh
#
#ENTRYPOINT ["/root/entrypoint.sh"]
