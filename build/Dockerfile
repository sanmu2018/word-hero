
FROM golang:1.25-alpine
WORKDIR /app/word-hero

ENV GOSUMDB=sum.golang.google.cn
ENV GO111MODULE=on

#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#https://mirrors.tuna.tsinghua.edu.cn/alpine#g' /etc/apk/repositories

RUN apk --no-cache add git pkgconfig build-base

ADD go.mod .
ADD go.sum .
RUN go mod download
ADD . .

RUN CGO_ENABLED=0 GO111MODULE=${GO111MODULE} go build -buildmode=pie  -o /go/bin/word-hero cmd/main.go


FROM alpine:3.11
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk --no-cache add ca-certificates libdrm curl zip
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata
WORKDIR /app/
COPY --from=0 /go/bin/word-hero .
COPY --from=0 /app/word-hero/configs ./configs
COPY --from=0 /app/word-hero/web ./web
CMD ["./word-hero"]