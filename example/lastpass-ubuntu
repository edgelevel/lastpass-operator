FROM phusion/baseimage:latest

RUN apt-get update && apt-get upgrade -y

RUN apt-get update && apt-get --no-install-recommends install -y \
  bash-completion \
  build-essential \
  cmake \
  libcurl3  \
  libcurl3-openssl-dev  \
  libssl1.0.0 \
  libssl-dev \
  libxml2 \
  libxml2-dev  \
  pkg-config \
  ca-certificates \
  xclip \
  lastpass-cli && \
  apt-get clean
