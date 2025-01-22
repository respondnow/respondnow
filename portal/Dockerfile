FROM registry.access.redhat.com/ubi8/ubi-minimal:8.8
RUN microdnf update --refresh --best --noplugins --setopt=install_weak_deps=0
RUN microdnf module enable nginx:1.20
RUN microdnf install nginx

COPY dist /opt/respond-now
COPY nginx/nginx.conf /etc/nginx/
COPY ./entrypoint.sh /opt

WORKDIR /opt/respond-now

RUN chown 65534:65534 -R /opt/respond-now
RUN chown 65534:65534 -R /var/log/nginx
RUN chown 65534:65534 -R /etc/nginx

USER 65534

EXPOSE 8191

ENTRYPOINT ["sh", "/opt/entrypoint.sh"]
