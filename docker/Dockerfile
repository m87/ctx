FROM alpine:3.20

RUN addgroup -S ctx && \
    adduser -S ctx -G ctx -h /home/ctx

RUN apk add --no-cache tzdata
ENV TZ=Europe/Warsaw

WORKDIR /home/ctx

COPY ctx /usr/local/bin/ctx
COPY ctx.yaml /home/ctx/.ctx

RUN mkdir /data
RUN chown -R ctx:ctx /home/ctx
RUN chown -R ctx:ctx /data

EXPOSE 8080

USER ctx

CMD ["ctx", "serve"]
