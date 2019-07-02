# UPI

- `Ignition config` files for the bootstrap, master, and worker machines
- Every `control plane machine` in an OpenShift Container Platform 4.1 cluster must use `RHCOS`
  - critical first-boot provisioning tool called `Ignition`
    - updates are delivered as an `Atomic OSTree repository`
  - The installation program uses Ignition config files to set the exact state of each machine, and the Machine Config Operator completes more changes to the machines, such as the application of new certificates or keys, after installation.



Necessary Files
- install-config.yaml
- Kubernetes manifests
- ignition config
  
**Necessary components:**
- Control Planes
- Compute machines
- Load Balancers
- Cluster networking(including the DNS)
- Storage for cluster infrastructure

**Required machines**
- One bootstrap machine
- Three control plane, or master, machines
- At least two compute, or worker, machines

**Minimum resource requirements**

| Machine   | OS    | vCPU | RAM  | Storage |
| --------- | ----- | ---- | ---- | ------- |
| Bootstrap | RHCOS | 4    | 16GB | 120GB   |
| Control   | RHCOS | 4    | 16GB | 120GB   |
| APP       | RHCOS or RHEL 7.6 | 4    | 16GB | 120GB   |


**Required components**
- Provision the required load balancers.
- Configure the ports for your machines.
- Configure DNS.
- Ensure network connectivity.



**What bootstrap do:**
- The bootstrap machine boots and starts hosting the remote resources - required for the master machines to boot. (Requires manual intervention - if you provision the infrastructure)
 
- The master machines fetch the remote resources from the bootstrap machine - and finish booting. (Requires manual intervention if you provision the - infrastructure)
 
- The master machines use the bootstrap machine to form an etcd cluster.
 
- The bootstrap machine starts a temporary Kubernetes control plane using the new etcd cluster.
 
- The temporary control plane schedules the production control plane to the - master machines.
 
- The temporary control plane shuts down and passes control to the - production control plane.
 
- The bootstrap machine injects OpenShift Container Platform components - into the production control plane.
 
- The installation program shuts down the bootstrap machine. (Requires - manual intervention if you provision the infrastructure)
  
- The control plane sets up the worker nodes.
  
- The control plane installs additional services in the form of a set of Operators.
 
 
 **Your machines must have direct internet access to install the cluster.**

- Access the OpenShift Infrastructure Providers page to download the installation program

- Access quay.io to obtain the packages that are required to install your cluster

- Obtain the packages that are required to perform cluster updates

- Access Red Hat’s software as a service page to perform subscription management

