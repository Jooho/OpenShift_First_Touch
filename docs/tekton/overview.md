Tekton Overview
---------------

![Image of Tekton](./images/tekton-horizontal-color.png)


Tekton Piplelines project provides Kubernetes-style resources for declaring CI/CD-style pipelines.

**Tekton** is :
- not only for knative
- Knative is first class citizen
- very early days


**Tekton Roadmap** :
- Conditional Execution
- Build Results and Logs
- Pluggable tasks
- Triggering
- community library

## History
```
2018      - Kanative build
2018 half - knative build-pipelines
2019      - Tekton Pipelines
```

*Knative Pipelines* become promoted from Knative project and it is *TEKTON Pipelines.*

## How Tekton works

**Architecture**

![Image of Tekton](./images/tekton-architecture.svg)

**Tekton based on CRDs**
- CRDs(Custom Resource Definition)
- Controllers + CRDs = Declarative API
- Controllers ==> CRD changes state: *Create/Update/Delete*


**Task CRD**
- New CRD
- *Sequence of steps*
- Run in sequential order
- Run on the same k8s node
- In the same node

**Pipeline CRD**
- Express Tasks order
  - Sequentially
  - Concurrently
  - (Graph)
- Can execute Tasks on different nodes
- Link inputs and outputs

**Runitme CRDs**
- Instances of Pipeline/Task:
  - PipelineRun
  - TaksRun
- PipelineResource
  - *Runtime info* like image registry,git repo, etc


## Kubernetes Challenges
**Ordering Steps Issue**
- Task steps run in order
- Task steps should be able to share disk
- How to schedule containers to both:
  - Run one at a time
  - Run on the same node

**Ordering Steps Solution**
- All steps in the same pod
- "Entrypoint" binary is injected
- All containers start
- Each step pauses until it should run

**Sidecars Issue**
- Don't start steps until sidecars are ready
- Sidecars can keep your pods alive forever

**Sidecars Solution**
- Check for sidecar readiness via Downward API
- Stop pods when steps finish:
  - Change sidecar image to no-op image

**Outputs --> Inputs Issue**
- Outputs of a Task can be porvided as Input to dependent Tasks
- Outputs = 
  - Binary data
  - Strings (TBD)

**Outputs --> Inputs Solution**
1. Attach PVCs
2. upload to blob store
3. (Theoretical) Schedule onto same node
   



## Reference
- [Modern CI/CD with Tekton and Prow Automated via Jenkins X - James Rawlings, Cloudbees](https://www.youtube.com/watch?v=4EyTGYB7GvA)
