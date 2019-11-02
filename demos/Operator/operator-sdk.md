# Operator SDK

## Build
```
mkdir -p $GOPATH/src/github.com/operator-framework
cd $GOPATH/src/github.com/operator-framework
git clone https://github.com/operator-framework/operator-sdk
cd operator-sdk
git checkout master
make dep
make install
```

## Download
```
wget https://github.com/operator-framework/operator-sdk/releases/download/v0.12.0/operator-sdk-v0.12.0-x86_64-linux-gnu -O operator-sdk

chmod 77 ./operator-sdk
cp ./operator-sdk /usr/bin/.
```
