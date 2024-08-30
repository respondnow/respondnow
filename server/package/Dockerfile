# Stage1: BUILD
FROM golang:1.22 AS builder

LABEL maintainer="RespondNow"

ARG TARGETOS=linux
ARG TARGETARCH

ADD . /respondnow
WORKDIR /respondnow

ENV GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

RUN go env
RUN GO111MODULE=on GOBIN=/usr/local/bin go install github.com/swaggo/swag/cmd/swag@v1.16.3
RUN CGO_ENABLED=0 go build -o /output/respondnow -v ./cmd/

# Stage2: PACKAGING
#
# Use RedHat UBI minimal image as base
FROM registry.access.redhat.com/ubi9/ubi-minimal:9.4

LABEL maintainer="RespondNow"

ENV APP_DIR="/respondnow"

COPY --from=builder /output/respondnow $APP_DIR/
RUN chown 65534:0 $APP_DIR/respondnow && chmod 755 $APP_DIR/respondnow

WORKDIR $APP_DIR
USER 65534

CMD ["./respondnow"]
EXPOSE 8080
