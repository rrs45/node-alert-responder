FROM box-registry.jfrog.io/jenkins/box-centos7

LABEL com.box.name="node-alert-responder"

# Required for systemd related things to work
ENV container=docker

ADD ./build/node-alert-responder /node-alert-responder
ADD ./config /config
RUN chown -R container:container /config  && chown container:container /node-alert-responder