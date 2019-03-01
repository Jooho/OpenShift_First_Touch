# PIP
파이썬 패키지 관리자. 

## Installation

```
curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
python get-pip.py
```

## Usage

- install package
  ```
  pip install -user $package_name
  ```
  `-user`는 현재 유저 폴더에 설치를 하는 것임. 
- remove package
  ```
  pip remove -user $package_name
  ```
