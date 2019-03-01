# Ansible Development Tutorial using molecule

## Content

### Ansible
- Ansible Overview
- Automation 
- Ansible AD-HOC
  - Example
- Ansible Playbook
  - Example
- Ansible Galaxy
  - Example

### Molecule
- [Molecule Overview](./molecule/molecule_overview.md)
- [Why Molecule?](./molecule/why_molecule.md)
- [Useful Tools](./tools/useful_tools.md)
  - [virtualenv](./tools/virtualenv.md)
  - [virtualenvwrapper](./tools/virtualenvwrapper.md)
- [Installation](./molecule/install.md)
- [Setup Test Environment](./molecule/setup_test_env.md)
- Variable Substitution
- Dependency
- Driver(Platforms)
  - Docker
  - Vagrant(virtualBox/Libvirt)
- Provisioner
- Lint
  - yamllint
  - ansible-lint
- Scenario
- Verifier
  - Lint(flake8)
  - Test(TestInfra)
    
## Molecule Practical Example

  - Create a new role
  - Test the role with a new image built by yourself
  - Test the role with a pre_built image
  - How to build an test image
  - Develop Test Python
  - Test a role on multiple platforms
  - 

## Best Practice 
  - Key Points
    - Solid Test Process 
    - Continous Integration (Travis)
  - Example



https://drive.google.com/drive/u/1/my-drive
https://molecule.readthedocs.io/en/latest/index.html
https://medium.com/@dan_kim/%ED%8C%8C%EC%9D%B4%EC%8D%AC-%EC%B4%88%EC%8B%AC%EC%9E%90%EB%A5%BC-%EC%9C%84%ED%95%9C-pip-%EA%B7%B8%EB%A6%AC%EA%B3%A0-virtualenv-%EC%86%8C%EA%B0%9C-a53512fab3c2

http://flake8.pycqa.org/en/latest/manpage.html

https://tox.readthedocs.io/en/latest/example/pytest.html#basic-example

https://www.jeffgeerling.com/blog/2018/testing-your-ansible-roles-molecule

https://docs.travis-ci.com/user/customizing-the-build/#build-matrix

https://testinfra.readthedocs.io/en/latest/

https://testinfra.readthedocs.io/en/latest/modules.html

https://www.jeffgeerling.com/blog/2018/how-i-test-ansible-configuration-on-7-different-oses-docker

https://opensource.com/article/18/12/testing-ansible-roles-molecule

https://blog.codecentric.de/en/2018/12/test-driven-infrastructure-ansible-molecule/

https://developer.rackspace.com/blog/molecule-for-existing-roles/

https://medium.com/@elgallego/ansible-role-testing-molecule-7a64f43d95cb

https://medium.com/@elgallego/molecule-2-x-tutorial-5dc1ee6b29e3

https://blog.codecentric.de/en/2018/12/continuous-infrastructure-ansible-molecule-travisci/