How to install Jenkin plugins?
-----------------------------

[Base Jenkins image](registry.access.redhat.com/openshift3/jenkins-2-rhel7) has minimal [plugins](https://github.com/openshift/jenkins/blob/master/2/contrib/openshift/base-plugins.txt).

Therefore, sometimes you need to install 3rd party plugins on top of the base jenkins image. The good thing is the base Jenkins image has a feature to install plugins easily.

### The concept of plugins installation in base Jenkins image
(Refer [Jenkins Github Repository](https://github.com/openshift/jenkins) README.md)

This use [install-plugin.sh](https://github.com/openshift/jenkins/blob/master/2/contrib/jenkins/install-plugins.sh). This script try to check metadata of the plugins in plugins.txt and extract the dependency list. Then, it will install all dependencies recursively.

*Example* If you install build-monitor-plugin, the token-marco:1.10 will be installed because of the chain of dependency.
```
 build-monitor-plugin -> view-job-filters ->  m2-extra-steps -> maven-plugin -> token-macro:1.10
```

**Note: In above case, the plugin will be conflicted because git-plugin that is installed by default needs token-macro:1.12 at least. Hence, token-macro:1.12 should be installed**


### This doc will show 2 ways:
 - DockerBuild on Local
 - OpenShift Build



### DockerBuild on Local
- Create Dockerfile
```
FROM registry.access.redhat.com/openshift3/jenkins-2-rhel7
COPY plugins.txt /opt/openshift/configuration/plugins.txt
RUN /usr/local/bin/install-plugins.sh /opt/openshift/configuration/plugins.txt
```
- Create plugins.txt
```
build-monitor-plugin:latest
slack:latest
ssh-agent:latest
token-macro:1.12
```

- Build Image && Push it docker registry
```
docker build -t ${DOCKER_REGISTRY_URL}/jenkins-2-rhel7-custom .
docker push ${DOCKER_REGISTRY_URL}/jenkins-2-rhel7-custom
```

### OpenShift Build
This will use above Dockerfile and plugins.txt

- Create BuildConfig
```
$ oc new-build  https://github.com/Jooho/jenkins-test-dockerbuild.git --strategy=docker

or

$ oc clone  https://github.com/Jooho/jenkins-test-dockerbuild.git && cd jenkins-test-dockerbuild
$ oc new-build .
```
- Start Build
```
$ oc start-build jenkins-test-dockerbuild
```

- Create DeploymentConfig to deploy the new image 
```
oc new-app -i jenkins-test-dockerbuild
```

