FROM alpine:latest as builder

RUN apk --update add ca-certificates

FROM scratch
# Import the certificates into the image
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY autobot /

ENTRYPOINT ["/autobot"]
