FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/

COPY autobot /

ENTRYPOINT ["/autobot"]
