FROM alpine:latest as builder
RUN apk --update add ca-certificates

FROM scratch
# Import from certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY autobot-plugin-hello-world /

ENTRYPOINT ["/autobot-plugin-hello-world"]
