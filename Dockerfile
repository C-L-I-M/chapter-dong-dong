FROM rust

COPY ./target/release/chapter-dong-dong /chapter-dong-dong

CMD ["/chapter-dong-dong"]
