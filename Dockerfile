# Build the manager binary
FROM golang:1.21 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download -x

# Copy the go source
COPY cmd/main.go cmd/main.go
COPY api/ api/
COPY internal/controller/ internal/controller/
COPY pkg/ pkg/
COPY version/ version/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager cmd/main.go

# Download lastpass-cli
FROM ubuntu:22.04 as lastpass-cli

ENV LASTPASS_VERSION=1.5.0

# https://github.com/lastpass/lastpass-cli?tab=readme-ov-file#building
RUN apt update && apt --no-install-recommends -yqq install \
    bash-completion \
    build-essential \
    cmake \
    libcurl4  \
    libcurl4-openssl-dev  \
    libssl-dev  \
    libxml2 \
    libxml2-dev  \
    libssl3 \
    pkg-config \
    ca-certificates \
    xclip \
    wget

RUN wget -O lastpass-cli.tar.gz https://github.com/lastpass/lastpass-cli/releases/download/v${LASTPASS_VERSION}/lastpass-cli-${LASTPASS_VERSION}.tar.gz; \
    tar -xf lastpass-cli.tar.gz

WORKDIR /lastpass-cli-${LASTPASS_VERSION}

RUN make && make install

RUN mkdir -p /usr/lib/lastpass-cli \
&& ldd /usr/bin/lpass | grep '=>' | awk '{ print $3 }' | xargs cp -t /usr/lib/lastpass-cli

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/base:nonroot

ENV USER_UID=1001 \
    USER_NAME=lastpass-operator

WORKDIR /
COPY --from=builder /workspace/manager .
COPY --from=lastpass-cli /usr/bin/lpass /usr/bin/which /usr/bin/
COPY --from=lastpass-cli /usr/lib/lastpass-cli /lib
COPY --from=lastpass-cli /bin/sh /bin/echo /bin/

USER nonroot:nonroot

ENTRYPOINT ["/manager"]
