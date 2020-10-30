# DB
```
yum install mariadb
mysql -uroot1 -ppassword2 -h127.0.0.13 databasename
mysql -uroot -ppassword -h127.0.0.1 databasename < SQL-Dump.sql
```

## wildfly deploy
```
mvn package
mvn wildfly:deploy
mvn wildfly:undeploy
mvn -DskipTests wildfly:deploy
```

# Docker
```
docker run -p 10080:80 webserver
docker run --privileged -v /sbin:/mnt/sbin -v /bin:/mnt/bin -it rhel7 /bin/bash
docker-registry-cli infrastructure:5000 search ose ssl
```


# openshift
```
sudo yum -y install atomic-openshift-clients
```

# OpenShift Environment from source (nodejs)
```
process.env.OPENSHIFT_NODEJS_PORT
```


# NFS
```
mkdir /var/export/remembermysql
chown nfsnobody:nfsnobody /var/export/remembermysql
chmod 700 /var/export/remembermysql
/var/export/remembermysql   *(rw,sync,all_squash)
exportfs -a
```

# Mysql
```
select count(*) from Promotion;
```

## put db data
```
oc port-forward bookstore-mysql-1-dw85k 13306:3306
mysql -h127.0.0.1 -P13306 -ubookstoreapp -psecret \
    bookstoredb < DO290/labs/clustering-bookstore/bookstore.sql

```

#EAP
## cert
```
oc process -f eap70-mysql-persistent-s2i.json \
    -p APPLICATION_NAME=bookstore \
    -p HOSTNAME_HTTP=bookstore-template.cloudapps.lab.example.com \
    -p HOSTNAME_HTTPS=bookstore-template-secure.cloudapps.lab.example.com \
    -p SOURCE_REPOSITORY_URL=http://infrastructure.lab.example.com/bookstore \
    -p SOURCE_REPOSITORY_REF="" \
    -p CONTEXT_DIR="" \
    -p DB_JNDI=java:/jboss/datasources/mysql \
    -p DB_DATABASE=bookstore \
    -p DB_USERNAME=openshif \
    -p DB_PASSWORD=password \
    -p HTTPS_SECRET=eap-app-secret-member \
    -p HTTPS_KEYSTORE=keystore.jks \
    -p HTTPS_NAME=jboss \
    -p HTTPS_PASSWORD=mykeystorepass \
    -p JGROUPS_ENCRYPT_SECRET=eap-app-secret-member \
    -p JGROUPS_ENCRYPT_KEYSTORE=keystore.jks \
```



## Clustering
```
oc adm policy add-role-to-user view --serviceaccount=default

or
oc describe dc bookstore | grep "Service Account"
oc adm policy add-role-to-user view \
    --serviceaccount=eap7-service-account

# Web-inf
<distributable/>

#hiberneate
   			<property name="hibernate.dialect" value="org.hibernate.dialect.MySQLDialect" />
			<property name="hibernate.hbm2ddl.auto" value="none" />
			<property name="hibernate.show_sql" value="false" />
			<property name="hibernate.format_sql" value="true" />
```


## EAP cli connect
```
 oc exec -it <pod_name> -- /opt/eap/bin/jboss-cli.sh --connect
```

## EAP Debug
```
The Java Debug Wire Protocol (JDWP) is the protocol used for communication between a debugger and the Java virtual machine (JVM) that allows for a debugger tool (such as that found in JBoss Developer Studio) to connect to it remotely and step through the source code as it is being executed in real time. One common way to enable debug mode in EAP is to edit the /opt/eap/bin/standalone.conf file by uncommenting the following line:

JAVA_OPTS="$JAVA_OPTS -agentlib:jdwp=transport=dt_socket,address=8787,server=y,suspend=n"
```

### For xPass
```
{
    "kind": "Template",
    "apiVersion": "v1",
    "metadata": {
        "name": "eap70-basic-s2i",
...
        {
            "kind": "DeploymentConfig",
...
                    "spec": {
                        "containers": [
                            {
...
                                "env": [
                                    {
                                        "name": "DEBUG",
                                        "value": "true"
                                    },
...   
```
```
oc env dc/deployment DEBUG=true
```

## jboss-cli.sh
```
/opt/eap/bin/jboss-cli.sh --connect
  --command="/subsystem=logging/logger=org.infinispan/:add(category=org.infinispan,level=TRACE,use-parent-handlers=true)"

```
### connect 
```
oc port-forward pod_name [-c container_name] localPort:8787
```

### if there are multiple pods, make it debuggerable
```
$ oc label pod myapp-1-xyzwx debugging=rightnow
$ oc expose pod myapp-1-xyzwx --name=debug-session --port=8080
$ oc expose svc debug-session --hostname=debug-session-myapp.example.com
```


SQL data preload
```
# src/main/resources/runtime.properties
preloaddata=true

# src/main/resources/META-INF/persistence.xml
hibernate.hbm2ddl.auto=create-drop
```

# Jenkins
```
node('maven') { 1
  stage 'stage-one' 2
    openshiftBuild(buildConfig: 'mybuildconfig', showBuildLogs: 'true') 3
  stage 'stage-two' 4
    openshiftDeploy(deploymentConfig: 'mydeploymentconfig') 5
}
```

## Jenkins pipeline
```
{
    "apiVersion": "v1",
    "kind": "BuildConfig",
    "metadata": {
        "name": "mypipeline",
        ...
    }
    "spec": {
        "strategy": {
            "jenkinsPipelineStrategy": {
                "jenkinsfile": "node('maven') {\n  stage 'stage-one\n    openshiftBuild(buildConfig: 'mybuildconfig', showBuildLogs: 'true')\n  stage 'stage-two'\n    openshiftDeploy(deploymentConfig: 'mydeploymentconfig')\n}" 1
            },
            "type": "JenkinsPipeline"
        }
        ...
    }
}
```

## Diable trigger
```
"imageChangeParams": {
                            "automatic": true,
```


# S2I

## New S2I image
```
sti create mybuilder myproject
```

## S2I script path
```
LABEL io.openshift.s2i.scripts-url=image:///usr/local/sti
```

## s2i build test
```
s2i build http://mygitserver/testapprepo mybuilder testapp

s2i build \
    file:///home/student/httpd-s2i/test/test-app \
    httpd-s2i:latest httpd-s2i-test
```

### S2I notice
```

The docker command can only be run as root from a RHEL, CentOS, or Fedora system. Developer users are advised to use the sudo command to get the necessary root access. The s2i utility makes Docker API calls, so it also needs to be run as root.
```

## IS
```
The image stream resource definition has to include an spec.tags.annotations object attribute with supports and tags attributes. 
```

## 


Tip
```
 dd if=/dev/zero of=/dev/null
 ```

