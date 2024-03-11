FROM alpine:3.19.1

RUN mkdir /project
WORKDIR /project

COPY config .
COPY scripts/sql .
COPY build/* .

ENTRYPOINT ["./cmd"]
