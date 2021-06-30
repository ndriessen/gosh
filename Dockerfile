FROM golang:1.16.5-alpine3.14

COPY ./bin/gosh-linux-amd64 /bin/gosh
RUN chmod +x /bin/gosh

WORKDIR /workdir

ENTRYPOINT ["/bin/gosh", "-w", "/workdir"]
