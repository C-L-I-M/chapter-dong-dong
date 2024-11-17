FROM golang:1-alpine AS builder

ADD . /build

WORKDIR /build

RUN go build -o /chapter-dong-dong

FROM alpine:3 as runner

COPY --from=builder ./chapter-dong-dong /chapter-dong-dong

CMD ["/chapter-dong-dong"]
