# Ansible Development Tutorial using molecule

## Content

- Ansible Overview
  - Automation 
  - Ansible AD-HOC
    - Example
  - Ansible Playbook
    - Example
  - Ansible Galaxy
    - Example

- [Molecule Overview](./molecule/molecule_overview.md)
  - Testing Process
  - Useful Tools
    - virtualenv
    - virtualenvwrapper
  - Installation
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
    

- Practical Example
  - Create a new role
  - Test the role with a new image built by yourself
  - Test the role with a pre_built image
  - How to build an test image
  - Develop Test Python
  - Test a role on multiple platforms
  - 

- Best Practice 
  - Key Points
    - Solid Test Process 
    - Continous Integration (Travis)
  - Example