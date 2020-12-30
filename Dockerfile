FROM gcr.io/distroless/static
ARG VERSION=none
COPY bin/gokmp-linux-${VERSION} /gokmp

ENTRYPOINT ["/gokmp"]
