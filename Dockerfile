FROM golang:1.18 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.io

COPY . /app
WORKDIR /app

RUN go get && go build -ldflags="-s -w" -installsuffix cgo

FROM debian:buster-slim

ENV TZ=Asia/Shanghai \
    LANG=C.UTF-8 \
    APP_DIR=/usr/local/go_found

COPY --from=builder /app/gofound ${APP_DIR}/gofound

WORKDIR ${APP_DIR}

RUN ln -snf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && chmod +x gofound

EXPOSE 5678

CMD ["./gofound"]
