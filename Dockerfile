FROM alpine:3.19.1

WORKDIR /project

COPY config config
COPY scripts/sql scripts/sql
COPY build/alpine/* .
COPY stubs/clothes images

ENTRYPOINT ["./cmd"]
