# Build Stage
FROM golang:1.14 AS build-stage

LABEL APP="build-example-api"
LABEL REPO="https://github.com/flamefatex/example-api"

ADD . /go/src/github.com/flamefatex/example-api
WORKDIR /go/src/github.com/flamefatex/example-api

RUN make build-alpine

# Final Stage
FROM alpine:3.12

ARG GIT_COMMIT
ARG VERSION
ARG APP_NAME

LABEL REPO="https://github.com/flamefatex/example-api"
LABEL GIT_COMMIT=${GIT_COMMIT}
LABEL VERSION=${VERSION}
LABEL APP_NAME=${APP_NAME}

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache tcpdump lsof net-tools tzdata curl dumb-init libc6-compat

ENV TZ Asia/Shanghai
ENV PATH=$PATH:/opt/example-api/bin

WORKDIR /opt/example-api/bin

COPY --from=build-stage /go/src/github.com/flamefatex/example-api/bin/example-api .
RUN chmod +x /opt/example-api/bin/example-api

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/example-api/bin/example-api"]