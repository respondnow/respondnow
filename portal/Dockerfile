# Use AlmaLinux minimal base image with ARM64 & AMD64 support
FROM debian:stable-slim AS base

RUN dnf update -y && \
    dnf install -y nginx && \
    dnf clean all

# Copy application files
COPY dist /opt/respond-now
COPY nginx/nginx.conf /etc/nginx/
COPY ./entrypoint.sh /opt

WORKDIR /opt/respond-now

# Ensure proper permissions
RUN chown 65534:65534 -R /opt/respond-now
RUN chown 65534:65534 -R /var/log/nginx
RUN chown 65534:65534 -R /etc/nginx

USER 65534

EXPOSE 8191

ENTRYPOINT ["sh", "/opt/entrypoint.sh"]
