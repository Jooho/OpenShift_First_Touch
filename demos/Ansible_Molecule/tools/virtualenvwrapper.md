# Virtualenvwrapper
`virtualenvwrapper` is a set of extensions to Ian Bicking’s virtualenv tool. The extensions include wrappers for creating and deleting virtual environments and otherwise managing your development workflow, making it easier to work on more than one project at a time without introducing conflicts in their dependencies.



## Installation
```
pip install --user virtualenvwrapper
```

## Usage

- Create a directory for `virtualenvwrapper`
  ```
  export WORKON_HOME=~/VirtualEnvs
  mkdir -p $WORKON_HOME
  ```

- Initialize `virtualenvwrapper`
  ```
  source ~/.local/bin/virtualenvwrapper.sh 
  ```

- Create a virtual env
  ```
  mkvirtualenv ansible
  ```

- Remove a virtual env
  ```
  rmvirtualenv ansible
  ```

- Activate a virtual env
  ```
  workon ansible
  ```

- Deactivate a virtual env
  ```
  deactivate
  ```

- Install a package on all virtual environments
  ```
  allvirtualenv pip install -U ansible
  ```

- List site-packages
  ```
  lssitepackages
  ```

- List Virtual Environment
  ```
  ls $WORKON_HOME
  ```

- Post packages installation (postmkvirtualenv)
  ```
  echo 'pip install sphinx' >> $WORKON_HOME/postmkvirtualenv
  ```
  

## Update ~/.bashrc
```
export WORKON_HOME=~/VirtualEnvs
source ~/.local/bin/virtualenvwrapper.sh 
```


## Reference
- [Official doc](https://virtualenvwrapper.readthedocs.io/en/latest/)