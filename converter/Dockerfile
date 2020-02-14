FROM golang:1.13-alpine AS builder

COPY .  /go/src/github.com/eedama/quiz-for-kids/converter
WORKDIR /go/src/github.com/eedama/quiz-for-kids/converter

RUN cd /go/src/github.com/eedama/quiz-for-kids/converter \
 	&& go install

FROM alpine:3.11 AS runtime

# Build argument
ARG VERSION
ARG BUILD
ARG NOW

# Install tini to /usr/local/sbin
ADD https://github.com/krallin/tini/releases/download/v0.18.0/tini-muslc-amd64 /usr/local/sbin/tini

# Install runtime dependencies & create runtime user
RUN apk --no-cache --no-progress add ca-certificates \
 && chmod +x /usr/local/sbin/tini && mkdir -p /opt \
 && adduser -D x0rzkov -h /opt/quiz-for-kids -s /bin/sh \
 && su x0rzkov -c 'cd /opt/quiz-for-kids; mkdir -p bin config data'

# Switch to user context
USER x0rzkov
WORKDIR /opt/quiz-for-kids/data

# Copy twint-docker binary to /opt/twint-docker/bin
COPY --from=builder /go/bin/converter /opt/quiz-for-kids/bin/converter
ENV PATH $PATH:/opt/quiz-for-kids/bin

# Container metadata
LABEL name="twint" \
      version="$VERSION" \
      build="$BUILD" \
      architecture="x86_64" \
      build_date="$NOW" \
      vendor="eedama" \
      maintainer="x0rzkov <x0rzkov@protonmail.com>" \
      url="https://github.com/eedama/quiz-for-kids" \
      summary="Dockerized converter from kahoot to opentdb cli tool" \
      description="Dockerized converter from kahoot to opentdb cli tool" \
      vcs-type="git" \
      vcs-url="https://github.com/eedama/quiz-for-kids" \
      vcs-ref="$VERSION" \
      distribution-scope="public"

# Container configuration
VOLUME ["/opt/quiz-for-kids/data"]
ENTRYPOINT ["tini", "-g", "--"]
CMD ["/opt/quiz-for-kids/bin/converter"]
