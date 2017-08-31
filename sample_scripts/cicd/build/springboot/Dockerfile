FROM registry.access.redhat.com/redhat-openjdk-18/openjdk18-openshift

MAINTAINER Jooho Lee <jlee@redhat.com>

# Build source
git clone https://github.com/${YOUR_GITHUB}/spring-cloud-kubernetes-sample.git

# Copy jar into /deployments


USER 185

WORKDIR /opt/java

CMD [ "/usr/local/s2i/run" ]

