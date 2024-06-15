FROM rust:latest AS builder

COPY . .

RUN cargo build --release

FROM alpine:latest AS alpine

RUN apk add gcompat libc6-compat libc++

COPY --from=builder ./target/release/chapter-dong-dong /app/chapter-dong-dong

CMD ["/app/chapter-dong-dong"]
