FROM golang:1.22-alpine AS builder

WORKDIR /

RUN wget https://github.com/ksmartdata/mysqld_exporter/archive/refs/heads/customizations.zip -O mysqld_exporter.zip && \
    unzip mysqld_exporter.zip && \
    rm -rf mysqld_exporter.zip

RUN cd /mysqld_exporter-customizations && \
    go mod tidy && \
    CGO_ENABLED=0 go build -o /mysqld_exporter

FROM quay.io/prometheus/busybox:latest

COPY --from=builder /mysqld_exporter /bin/

RUN chown nobody:nobody /bin/mysqld_exporter && \
    chmod 755 /bin/mysqld_exporter

EXPOSE 9104

USER nobody

ENTRYPOINT ["/bin/mysqld_exporter"]
