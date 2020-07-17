ssh dell-r60


ssh stack@10.37.197.2 BRQ  

DNS server 10.37.192.157
vi overcloudrc

/home/stack/heat_templates

USer
Project
key

commands
```
openstack loadbalancer amphora list
openstack endpoint list
openstack  server list --all-projects
```

Editable Variables
- master_vms: 3
- infra_vms: 3
- app_vms: 3
- cluster_name: ${user_name}_ocp_3_11
- container_engine: docker
  - container_engine_storage_type:  overlay2
  - container_engine_storage_size:  20G
  - container_engine_storage_disk:  vdb
  - 


Optional Variables
- master_node_prefix: master
- infra_node_prefix: infra
- app_node_prefix: app
- lb_node_prefix


Role Based Test.

Roles
- Config OpenStack
  - Create project
  - Change project quota
  - Create user
  - Associate roles
  - Create keypair
  - Add keypair
  - Add security group rule
    - TCP
    - UDP
    - ICMP
    - SSH
  - Create flavor
    - master
    - infra
    - app
  - Create image
  - Add network
  - Add subnet to network
  - Create router for private and external routing

- Create VMs
  - master 
  - infra
  - app
  - lb
  - 
- Create Load Balancer
  - [Test] domain name

- Create DNS records
  - [Test] DDNS
    - https://mojo.redhat.com/people/jhutar/blog/2018/12/18/how-to-get-your-own-usersysredhatcom-hostname-aka-ddns
    - https://mojo.redhat.com/people/miabbott/blog/2018/05/31/conatinerizing-redhat-internal-ddns-client
  - subdomain
  - 
-----------------------------------------------------------

- Generate ansible inventory file
- Create 




[Test List]

- Is it possible to access to LB domain name by Red Hat VPN network?
- Subdomain using ddns can map to LB IP?


version: '3'
services:
  postgres:
    image: "postgres:9.6"
    environment:
      POSTGRES_USER: awx
      POSTGRES_PASSWORD: awxpass
      POSTGRES_DB: awx

  rabbitmq:
    image: "rabbitmq:3"
    environment:
      RABBITMQ_DEFAULT_VHOST: awx

  memcached:
    image: "memcached:alpine"

  awx_web:
    # image: "geerlingguy/awx_web:latest"
    image: "ansible/awx_web:latest"
    links:
      - rabbitmq
      - memcached
      - postgres
    ports:
      - "80:8052"
    hostname: awxweb
    user: root
    environment:
      SECRET_KEY: aabbcc
      DATABASE_USER: awx
      DATABASE_PASSWORD: awxpass
      DATABASE_NAME: awx
      DATABASE_PORT: 5432
      DATABASE_HOST: postgres
      RABBITMQ_USER: guest
      RABBITMQ_PASSWORD: guest
      RABBITMQ_HOST: rabbitmq
      RABBITMQ_PORT: 5672
      RABBITMQ_VHOST: awx
      MEMCACHED_HOST: memcached
      MEMCACHED_PORT: 11211

  awx_task:
    # image: "geerlingguy/awx_task:latest"
    image: "ansible/awx_task:latest"
    links:
      - rabbitmq
      - memcached
      - awx_web:awxweb
      - postgres
    hostname: awx
    user: root
    environment:
      SECRET_KEY: aabbcc
      DATABASE_USER: awx
      DATABASE_PASSWORD: awxpass
      DATABASE_NAME: awx
      DATABASE_PORT: 5432
      DATABASE_HOST: postgres
      RABBITMQ_USER: guest
      RABBITMQ_PASSWORD: guest
      RABBITMQ_HOST: rabbitmq
      RABBITMQ_PORT: 5672
      RABBITMQ_VHOST: awx
      MEMCACHED_HOST: memcached
      MEMCACHED_PORT: 11211


