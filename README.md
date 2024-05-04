# actimanager

Fine-grained resource management for Kubernetes Pods.

## Description

This repository is a framework for fine-grained orchestration of Kubernetes Pods. It utilizes the Operator Pattern, CRD's, Scheduling Framework plugins and other extension points of Kubernetes to provide a solution for managing the CPU resources of Pods at a fine-grained level.

## Components

- **Custom Resource Definitions (CRD's)**
    - `NodeCPUTopology`: The CPU and NUMA topology of a node.
    - `PodCPUBinding`: The allocated CPU resources of a Pod.
      ```yaml
      apiVersion: cslab.ece.ntua.gr/v1alpha1
      kind: PodCPUBinding
      metadata:
        name: benchpod-1-binding
        namespace: benchmarks
      spec:
        exclusivenessLevel: Core  # None, CPU, Core, Socket, NUMA
        podName: benchpod-3
        cpuSet:
          - cpuID: 10
          - cpuID: 12
      ```
- **Controller - Manager**
    - Watches for changes in the CRD's and reconciles them.
- **Custom Scheduler**
  - A custom scheduler that schedules Pods based on the CRD's.
  - Implements the `WorkloadAware` plugin, which schedules and binds Pods based on the workload family they belong in (MemoryBound, CPUBound, IOBound, BestEffort).
- **Daemon**
    - A gRPC server that runs on each node as a DaemonSet.
    - Exposes `Topology` and `CPUPinning` services to interact with each node.
    - Reconciles the resources of the Pods inside each host.
- **Utility Libraries**:
    - Code-generated clientset, informers, listers for CRD's
    - Protobuf definitions for the gRPC services.
    - Utility functions for interacting with the Kubernetes API.
    - Linux cgroup utilities for interacting with the host.


## Getting Started

### Prerequisites
- go version v1.20.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster

1. Clone the repository and navigate to the root directory.

    ```sh
    git clone
    cd actimanager
    ```
2. Edit the DaemonSet manifest under `config/daemon/daemon.yaml` to match your nodes' configuration.
    
    ```yaml
     args:
        - '--node-name=$(NODE_NAME)'
        - '--container-runtime=containerd'  # containerd, docker, kind
        - '--cgroups-path=/cgroup'
        - '--cgroups-driver=systemd'        # systemd, cgroupfs
        - '--reconcile-period=15s'
        - '--verbosity=3'
    ```
3. Install the components on the cluster.

    ```sh
    kubectl apply -k config/default
    ```

### To Uninstall

```sh
kubectl delete -k config/default
```

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

