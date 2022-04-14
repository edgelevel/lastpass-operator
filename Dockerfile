# Build the manager binary
FROM golang:1.17 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download -x

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY pkg/ pkg/
COPY version/ version/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go

# Download lastpass-cli
FROM golang:1.17 as lastpass-cli

ENV LASTPASS_VERSION=1.3.3-4

# https://github.com/lastpass/lastpass-cli
RUN apt update && apt install -y \
  lastpass-cli=${LASTPASS_VERSION} \
  && rm -rf /var/lib/apt/lists/*

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
