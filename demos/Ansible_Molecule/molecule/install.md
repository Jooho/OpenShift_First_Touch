# Molecule Installation

## Pre-requisites
Python 2
```
sudo yum install -y gcc python-pip python-devel openssl-devel libselinux-python
```

Python 3
```
sudo yum install -y gcc python3-pip python3-devel openssl-devel libselinux-python3
```


## Molecule

```
pip install --user molecule
```

## With virtualenvwrapper

```
mkvirtualenv molecule

or 

workon molecule

pip install --user molecule
```



## Using Docker
```
docker run --rm -it --privileged \
    -v "$(pwd)":/tmp/$(basename "${PWD}"):ro \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -w /tmp/$(basename "${PWD}") \
    quay.io/ansible/molecule:latest \
    sudo molecule test
```