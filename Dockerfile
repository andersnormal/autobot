FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
# Import from certs
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY autobot /

ENTRYPOINT ["/autobot"]
