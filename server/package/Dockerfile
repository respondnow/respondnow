# Stage1: BUILD
FROM golang:1.22 AS builder

LABEL maintainer="RespondNow"

ARG TARGETOS=linux
ARG TARGETARCH

ADD . /respond
WORKDIR /respond

ENV GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

RUN go env
RUN GO111MODULE=on GOBIN=/usr/local/bin go install github.com/swaggo/swag/cmd/swag@v1.16.3
RUN CGO_ENABLED=0 go build -o /output/respond -v ./cmd/

# Stage2: PACKAGING
#
# Use RedHat UBI minimal image as base
FROM registry.access.redhat.com/ubi9/ubi-minimal:9.4

LABEL maintainer="RespondNow"

ENV APP_DIR="/respond"

COPY --from=builder /output/respond $APP_DIR/
RUN chown 65534:0 $APP_DIR/respond && chmod 755 $APP_DIR/respond

WORKDIR $APP_DIR
USER 65534

CMD ["./respond"]
EXPOSE 8080
