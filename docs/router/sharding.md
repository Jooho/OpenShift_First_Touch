Router Sharding
---------------

[Router Sharding](https://docs.openshift.com/container-platform/3.5/architecture/core_concepts/routes.html#router-sharding) is running multiple router deployments each serving different sets of routes.  It's described in detail in the Routes section of the documentation.  Also see, Using the Default HAProxy Router

Two new nodes were provisioned and labeled region=dmzinfra and region=dmzprimary.

The new router (representing a DMZ) would run on nodes with region=dmzinfra, while the old router would continue to run on all nodes labeled region=infra.  

Router sharding is accomplished by giving different routers different selection criteria to limit the routes provided by a router.  For this experiment, namespaces/projects were used as a filter, although labels on routes can also be used to filter.

Routes in projects labeled zone=dmz would be serviced by the DMZ router.  All other routes would be serviced by the regular router.

First to set up the DMZ router:

```
[user@host ~ ]$ oc project default 
[user@host ~ ]$ oadm router router-shard-dmz --replicas=0 --selector='region=dmzinfra'            
[user@host ~ ]$ oc set env dc router-shard-dmz NAMESPACE_LABEL='zone=dmz'  
[user@host ~ ]$ oc scale dc router-shard-dmz --replicas=1 
```
Second modify the old router to not accept:
```
[user@host ~ ]$ oc project default 
[user@host ~ ]$ oc set env dc router NAMESPACE_LABEL='zone notin (dmz)'  

# This should trigger a redeploy, but if not manually redeploy it:
[user@host ~ ]$ oc rollout latest router
```

Third, label a project:
```
[user@host ~ ]$ oc label namespace dmztest zone=dmz
```
