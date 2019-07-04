FROM scratch

COPY server/autobot /

ENTRYPOINT ["/autobot"]
