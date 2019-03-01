# Setup Test Environment

## [Setup virtualenvwrapper](../tools/virtualenvwrapper.md)
```
pip install --user virtualenvwrapper
echo "export WORKON_HOME=~/VirtualEnvs" >> ~/.bashrc
echo "source ~/.local/bin/virtualenvwrapper.sh >> ~/.bashrc"

source ~/.bashrc

mkvirtualenv ansible

```

## [Setup molecule](./install.md)
```
sudo yum install -y gcc python-pip python-devel openssl-devel libselinux-python

pip install ansible molecule
```

## Creaet a new test role
```
workon ansible
molecule init role -r test_role -d docker
```

## Role Directory
```
cd ./test_role
tree

.
├── defaults
│   └── main.yml
├── handlers
│   └── main.yml
├── meta
│   └── main.yml
├── molecule
│   └── default
│       ├── Dockerfile.j2
│       ├── INSTALL.rst
│       ├── molecule.yml
│       ├── playbook.yml
│       └── tests
│           ├── __pycache__
│           │   ├── test_default.cpython-36.opt-1.pyc
│           │   └── test_default.cpython-36.pyc
│           └── test_default.py
├── README.md
├── tasks
│   └── main.yml
└── vars
    └── main.yml

9 directories, 13 files

```


## Docker
```
docker run --rm -it --privileged \
    -v "$(pwd)":/tmp/$(basename "${PWD}"):ro \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -w /tmp/$(basename "${PWD}") \
    quay.io/ansible/molecule:latest \
    sudo molecule test
```
*Note* Destroy show some error messages because the behavior of molecule using docker image is like docker in docker.
## 