# ETCD Recovery without SSL

## Download `etcdctl` binary ( Ceheck new release:https://github.com/etcd-io/etcd/releases)
```
wget -c https://github.com/etcd-io/etcd/releases/download/v3.3.12/etcd-v3.3.12-linux-amd64.tar.gz -O -|tar xvz

or 
curl -L https://github.com/etcd-io/etcd/releases/download/v3.3.12/etcd-v3.3.12-linux-amd64.tar.gz|tar xvz
```

## Create folders for ETCD data
```
mkdir /tmp/etcd{1,2,3}
```

## Start ETCD cluster
```
docker-compose up -d
```

## Break etcd-1 member cluster
```
 docker rm -f $(docker ps --filter="NAME=etcd-1" -q)

sudo rm -rf /tmp/etcd1/*
```

## Remove etcd-1 from cluster
```
ETCDCTL_API=3 ./etcdctl --endpoints=http://172.33.21.42:2379,http://172.33.21.43:2379 member remove $(ETCDCTL_API=3 ./etcdctl --endpoints=http://172.33.21.41:2379,http://172.33.21.42:2379,http://172.33.21.43:2379 member list|grep etcd-1|awk '{print $1}'|tr -d ',')
```
## Update etcd-1 to etcd-4 in docker-compose.yml
```
sed 's/etcd-1/etcd-4/g' -i ./docker-compose.yml 
```

## Add etcd-4 to cluster
```
ETCDCTL_API=3 ./etcdctl --endpoints=http://172.33.21.42:2379,http://172.33.21.43:2379 member add etcd-4  --peer-urls="http://etcd-4:2380"
```

## Start etcd-4
```
docker-compose up -d
```

## Check member list
```
ETCDCTL_API=3 ./etcdctl --endpoints=http://172.33.21.41:2379,http://172.33.21.42:2379,http://172.33.21.43:2379 member list
ETCDCTL_API=3 ./etcdctl --endpoints=http://172.33.21.41:2379,http://172.33.21.42:2379,http://172.33.21.43:2379 endpoint health
ETCDCTL_API=3 ./etcdctl --endpoints=http://172.33.21.41:2379,http://172.33.21.42:2379,http://172.33.21.43:2379 endpoint status -w table

```

## Clean up
```
docker-compose down --volume

sudo rm -rf /tmp/etcd{1,2,3}
sed 's/etcd-4/etcd-1/g'  -i ./docker-compose.yml 
```


## Tip
- Docker network 
```
docker network ls
docker network rm ${id}
```
