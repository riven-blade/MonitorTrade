FROM alpine:latest
WORKDIR /app
COPY ./bin/monitor-trade .
ENTRYPOINT ["./monitor-trade"]