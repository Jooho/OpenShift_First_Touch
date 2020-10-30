# Concept (Control Plane)

## Envoy Proxy
Features
- Dynamic service discovery
- Load balancing
- TLS termination
- /2 and gRPC proxies
- Circuit breakers
- Health checks
- Staged rollouts with percent-based traffic split
- Fault injection
- Rich metrics

Envoy deployed as sidecar to relevant service in same Kubernetes pod
Sidecar proxy model lets you add Istio capabilities to existing deployment
No need to rearchitect or rewrite code

## Mixer
- Enforces access control and usage policies across service mesh
- Collects telemetry data from Envoy proxy and other services
- Proxy extracts request-level attributes, sends to Mixer for evaluation

## Pilot

- Provides:
  - Service discovery for Envoy sidecars
  - Traffic management capabilities for intelligent routing—A/B tests, canary deployments, etc.
  - Resiliency—timeouts, retries, circuit breakers, etc.

- Converts high-level routing rules that control traffic behavior into Envoy-specific configurations
  - Propagates them to sidecars at runtime

- RulesAPI
  - Manages and configures Envoy proxy (sidecar) instances in service mesh
  - Allows you to specify routing rules to use between services in service mesh
  - Enables service discovery, dynamic updates for load balancing, routing tables

## Auth/Citadel

- Provides service-to-service and end user authentication
- Built-in identity and credential management
- Use to upgrade unencrypted traffic in service mesh
- Enables operators to enforce policies based on service identity, not network controls


## Galley

- Validates user-authored Istio API configuration on behalf of other control plane components
- Intended for future releases:
  - Take over responsibility as Istio's top-level configuration ingestion, processing distribution component
  - Insulate other components from details of obtaining user configuration from underlying platform (e.g., Kubernetes)


# Installation

## Using Operator
- `istio-operator` project
- `istio-system` project
- https://github.com/Maistra/istio-operator

```
oc new-project istio-operator
oc new-project istio-system
oc apply -n istio-operator -f ./deploy/maistra-operator.yaml

oc apply -n istio-system -f ./deploy/examples/maistra_v1_servicemeshcontrolplane_cr_basic.yaml
```

## Using Ansible
- https://github.com/Maistra/openshift-ansible/tree/maistra-0.11/istio
- https://github.com/Maistra/openshift-ansible/blob/maistra-0.11/istio/Installation.md    <Deprecated>
```

```

## Using Helm Chart
- https://github.com/istio/installer

```
helm template --namespace $NAMESPACE -n $COMPONENT $CONFIGDIR -f global.yaml | \
   kubectl apply -n $NAMESPACE --prune -l release=$COMPONENT -f -
```


## Note
- For user to access istio-system namespace, you should give edit permission to `istio-system`
  ```
  oc adm policy add-role-to-user edit user1 -n istio-system
  ```
- For user to run istio commands you should give system:admin cluster role
  ```
  oc adm policy add-cluster-role-to-user sudoer user1 --as=system:admin
  ```


# Kiali

Sample Repo: https://github.com/gpe-mw-training/ocp-service-mesh-foundations

## URL
```
export KIALI_URL=https://$(oc get route kiali -n istio-system -o template --template='{{.spec.host}}')
```