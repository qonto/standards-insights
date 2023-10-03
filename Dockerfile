FROM alpine:3.18

ARG HOME=/app

RUN apk add --update --no-cache ca-certificates

RUN addgroup -g 10001 -S standards \
    && adduser --home ${HOME} -u 10001 -S standards -G standards \
    && mkdir -p /app \
    && chown standards:standards -R /app

WORKDIR $HOME

USER 10001
WORKDIR ${HOME}
COPY standards-insights /app/

EXPOSE 3000

ENTRYPOINT ["/app/standards-insights"]

