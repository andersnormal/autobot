FROM scratch

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY autobot-plugin-slack /

ENTRYPOINT ["/autobot-plugin-slack"]
