# Use Debian as the base image (supports both ARM64 and AMD64)
FROM --platform=$TARGETPLATFORM debian:stable-slim AS base

# Update package lists and install Nginx
RUN apt-get update && \
    apt-get install -y nginx && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy application files
COPY dist /opt/respond-now
COPY nginx/nginx.conf /etc/nginx/
COPY ./entrypoint.sh /opt

WORKDIR /opt/respond-now

# Set permissions
RUN chown 65534:65534 -R /opt/respond-now
RUN chown 65534:65534 -R /var/log/nginx
RUN chown 65534:65534 -R /etc/nginx

USER 65534

EXPOSE 8191

ENTRYPOINT ["sh", "/opt/entrypoint.sh"]
