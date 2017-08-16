Retrieve Dockerfile from Image
-----------------------------

**Usage**
```
docker run -v /var/run/docker.sock:/var/run/docker.sock dockerfile-from-image <IMAGE_TAG_OR_ID>

docker run --rm -v /var/run/docker.sock:/var/run/docker.sock  ljhiyh/docker-from-file-image  docker-from-file-image
```

**Build Image**

- Clone git project
```
git clone https://github.com/CenturyLinkLabs/dockerfile-from-image.git
cp ./dockerfile-from-image/dockerfile-from-image.rb .
```

- Create Dockerfile
~~~
FROM alpine:3.2
MAINTAINER CenturyLink Labs <clt-labs-futuretech@centurylink.com>

RUN apk --update add ruby-dev ca-certificates && \
    gem install --no-rdoc --no-ri docker-api && \
    apk del ruby-dev ca-certificates && \
    apk add ruby ruby-json && \
    rm /var/cache/apk/*

ADD dockerfile-from-image.rb /usr/src/app/dockerfile-from-image.rb
RUN chmod +x /usr/src/app/dockerfile-from-image.rb

ENTRYPOINT ["/usr/src/app/dockerfile-from-image.rb"]
CMD ["--help"]
~~~

- Build Image
~~~
docker build -t docker-from-file-image .
~~~


