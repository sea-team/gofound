FROM golang:1.18 as builder

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io

COPY . /app
WORKDIR /app

RUN go get && go build

FROM debian:buster-slim

ENV TZ=Asia/Shanghai
ENV LANG=C.UTF-8
ENV APP_DIR=/usr/local/go_found

COPY --from=builder /app /usr/local/go_found
WORKDIR ${APP_DIR}

RUN ln -snf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && chmod +x gofound

EXPOSE 5678

CMD ["./gofound"]
