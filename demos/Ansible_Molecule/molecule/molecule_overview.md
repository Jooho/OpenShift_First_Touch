# Molecule Overview


> **Molecule** is designed to aid in the development and testing of Ansible roles. Molecule provides support for testing with multiple instances, operating systems and distributions, virtualization providers, test frameworks and testing scenarios. Molecule is opinionated in order to encourage an approach that results in consistently developed roles that are well-written, easily understood and maintained.

from [molecule homepage](https://molecule.readthedocs.io/en/latest/index.html)



## Documentation
- https://molecule.readthedocs.io/
- https://media.readthedocs.org/pdf/molecule/latest/molecule.pdf

## Main Functions

- Handle project linting by yamllint/ansible-lint.
- Execute roles in a specific platform
- Test role result by invoking configurable verifier(TestInfra)
  

## Scenario
```
scenario:
  name: default
  create_sequence:
    - create
    - prepare
  check_sequence:
    - destroy
    - dependency
    - create
    - prepare
    - converge
    - check
    - destroy
  converge_sequence:
    - dependency
    - create
    - prepare
    - converge
  destroy_sequence:
    - destroy
  test_sequence:
    - lint
    - destroy
    - dependency
    - syntax
    - create
    - prepare
    - converge
    - idempotence
    - side_effect
    - verify
    - destroy
```

### Scenario detail
- Lint 
  -  Run yamllint to format yaml style.
- Destory
  - Stop test VM/container
- Dependency
  - Download dependent Galaxy roles
- Syntax
  - Run ansible-lint to check Ansible syntax
- Create
  - Run VM/Container for test
- Prepare
  - **(Optional)** pre-requisite playbook
- Converge
  - Run your playbook
- Idempotence
  - Run the playbook again and check if there is any changed.
- Verify
  - Run TestInfra script
- Side_effect
  - **(Optional)** Intended to test HA failover senarios or the like
  - Experimential function