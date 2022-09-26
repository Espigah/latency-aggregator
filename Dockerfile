##########################
## builder image
##########################
FROM golang:1.17 as builder

ENV APP_HOME /build
WORKDIR $APP_HOME

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ENV GOOS linux
ENV GOARCH ${GOARCH:-amd64}
ENV CGO_ENABLED=0
ENV GIN_MODE=release

RUN go build -v -o latency-aggregator cmd/latency-aggregator/main.go

##########################
## compress binary
##########################
FROM golang:1.17 as upx
ENV APP_HOME /build
WORKDIR $APP_HOME
ARG upx_version=3.96
ARG upx=9
# hadolint ignore=DL3008
RUN apt-get update && apt-get install -y --no-install-recommends xz-utils && \
    curl -Ls https://github.com/upx/upx/releases/download/v${upx_version}/upx-${upx_version}-amd64_linux.tar.xz -o - | tar xvJf - -C /tmp && \
    cp /tmp/upx-${upx_version}-amd64_linux/upx /usr/local/bin/ && \
    chmod +x /usr/local/bin/upx && \
    apt-get remove -y xz-utils && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder $APP_HOME/latency-aggregator latency-aggregator
# Compress the binary and copy it to final image
RUN upx -${upx} latency-aggregator

##########################
## create user
##########################
FROM ubuntu:22.04 as user

ENV APP_HOME                            /build
ENV APP_USER                            app
ENV APP_GROUP                           app

COPY --from=builder $APP_HOME/latency-aggregator  $APP_HOME/latency-aggregator

RUN mkdir -p ${APP_HOME}/ && \
    chmod +x $APP_HOME/latency-aggregator

RUN groupadd -r ${APP_GROUP} && \
    useradd -g ${APP_GROUP} -d ${APP_HOME} -s /sbin/nologin  -c "Unprivileged User" ${APP_USER} && \
    chown -R ${APP_USER}:${APP_GROUP} ${APP_HOME} 

RUN chsh --shell /sbin/nologin root

################################
## create ca-certificates
################################ 
FROM alpine:3.6 as alpine

RUN apk add -U --no-cache ca-certificates


################################
## generate clean, final image
################################ 
FROM scratch

ARG VERSION
ENV APP_VERSION $VERSION
ENV GIN_MODE release
ENV APP_HOME /build
ENV APP_USER app

COPY --from=upx $APP_HOME/latency-aggregator  $APP_HOME/latency-aggregator
COPY --from=user /etc/passwd /etc/passwd
COPY --from=user /etc/group /etc/group
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY  README.md /README.md


USER $APP_USER
WORKDIR $APP_HOME

EXPOSE 7070

CMD ["./latency-aggregator"]
