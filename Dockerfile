FROM golang:1.24-alpine AS builder
ARG PKG=github.com/kashalls/minecraft-router-sidehook
ARG VERSION=dev
ARG REVISION=dev
ARG CMD_TARGET=discord
WORKDIR /build
COPY . .
RUN go build -o /build/webhook -ldflags "-s -w -X main.Version=${VERSION} -X main.Gitsha=${REVISION}" ./cmd/${CMD_TARGET}

FROM gcr.io/distroless/static-debian12:nonroot
USER 8675:8675
COPY --from=builder --chmod=555 /build/webhook /webhook
EXPOSE 8888/tcp
ENTRYPOINT ["/webhook"]
