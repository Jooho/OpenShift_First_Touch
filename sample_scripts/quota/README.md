LimitRange Demo
---

**Scenario 1. Set Limit on Pod**

*CPU*
1. Deploy a pod-cpu-load that is BestEffort QoS.

2. Check CPU usage

3. Kill

4. Depoy a limit-pod-cpu-load that is Burstable QoS

5. Check CPU usage.

*Memory*
1. Deploy a pod-memory-load that is BestEffort QoS.

2. Check Memory usage

3. Kill

4. Depoy a limit-pod-cpu-load that is Burstable QoS

5. Check CPU usage.



**Scenario 2. Set Limit by defaultProjectTemplate**
1. Create defaultProjectTemplate
2. Apply it to master-config.yml && Restart Master
3. Deploy dc-cpu-mem-load that is BestEffor QoS.


Quota Demo
---

Job

Build

BestEffort

NotBestEffort

compute-resources
- pods
- request.memory
- limits.cpu
- limits.memory


Require Explicit Quota to Consume a Resource
