FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/

COPY autobot-plugin-slack /

ENTRYPOINT ["/autobot-plugin-slack"]
