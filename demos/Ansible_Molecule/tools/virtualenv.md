# Virtualenv

Virtualenv help avoid confliction of python projects on a host.

## Installation

If you don't have [pip](./pip.md), please install it.

```
pip install virtualenv 
```


## Usage

- Create virtual environment(folder)
  ```
  virtualenv ~/ansible 
  ```
  `~/ansible` is the folder that you will store python packages so you can change

- Activate a virtual environment
  ```
  source ~/ansible/bin/activate
  ```

- Deactivate a virtual environment
  ```
  deactivate
  ```


## Reference

- [Dan Kim Medium](https://medium.com/@dan_kim/%ED%8C%8C%EC%9D%B4%EC%8D%AC-%EC%B4%88%EC%8B%AC%EC%9E%90%EB%A5%BC-%EC%9C%84%ED%95%9C-pip-%EA%B7%B8%EB%A6%AC%EA%B3%A0-virtualenv-%EC%86%8C%EA%B0%9C-a53512fab3c2)