# syntax=docker/dockerfile:1.6

########################
# BUILD STAGE
########################
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
ARG RELEASE=dev
ARG COMMIT=
ARG DATE=

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w -X github.com/m87/ctx/server.Release=${RELEASE} -X github.com/m87/ctx/server.Commit=${COMMIT} -X github.com/m87/ctx/server.Date=${DATE}" -o /out/ctx ./


########################
# RUNTIME STAGE
########################
FROM --platform=$TARGETPLATFORM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
  && adduser -D -H -u 10001 appuser

COPY --from=build /out/ctx /usr/local/bin/ctx

RUN mkdir /data
RUN mkdir /blobs

EXPOSE 8080
ENV TZ=Europe/Warsaw
ENV DATABASE_PATH=/data/ctx.db

CMD ["ctx", "serve", "--addr", ":8080"]
