FROM alpine:latest

ENV OPERATOR=/usr/local/bin/lastpass-operator \
    USER_UID=1001 \
    USER_NAME=lastpass-operator

# https://github.com/lastpass/lastpass-cli
RUN apk add --update --no-cache \
  lastpass-cli && \
  rm -f /var/cache/apk/*

RUN install -d -m 0755 -o ${USER_UID} -g ${USER_UID} .lpass

# install operator binary
COPY build/_output/bin/lastpass-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
