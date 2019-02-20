# Molecule Overview


> **Molecule** is designed to aid in the development and testing of Ansible roles. Molecule provides support for testing with multiple instances, operating systems and distributions, virtualization providers, test frameworks and testing scenarios. Molecule is opinionated in order to encourage an approach that results in consistently developed roles that are well-written, easily understood and maintained.

from [molecule homepage](https://molecule.readthedocs.io/en/latest/index.html)



## Documentation
- https://molecule.readthedocs.io/


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