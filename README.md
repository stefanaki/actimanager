# actimanager

Fine-grained resource management for Kubernetes Pods.

## Description

This repository is a framework for fine-grained orchestration of Kubernetes Pods. It utilizes the Operator Pattern, CRD's, Scheduling Framework plugins and other extension points of Kubernetes to provide a solution for fine-grained resource allocation. It consists of the following components:

- **Custom Resource Definitions (CRD's)**
    - `NodeCpuTopology`: The CPU and NUMA topology of a node.
    - `PodCpuBinding`: The allocated CPU resources of a Pod.
      ```yaml
      apiVersion: cslab.ece.ntua.gr/v1alpha1
      kind: PodCpuBinding
      metadata:
        name: benchpod-1-binding
        namespace: benchmarks
      spec:
        exclusivenessLevel: Core
        podName: benchpod-3
        cpuSet:
          - cpuId: 10
          - cpuId: 12
      ```
- **Controller - Manager**
    - Watches for changes in the CRD's and takes action accordingly.
- **Custom Scheduler**
  - A custom scheduler that schedules Pods based on the CRD's.
  - Implements multiple scheduling policies for different use cases.
- **Daemon**
    - A gRPC server that runs on each node and tries to bind the Pods to the CPU's according to the CRD's.
- **Utility Libraries**:
    - Code-generated clientset, informers, listers for CRD's
    - Utility functions for interacting with the Kubernetes API.
    - Linux `cgroup` utilities for interacting with the host.


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
       - '--container-runtime=docker' # containerd, kind
       - '--cgroups-path=/cgroup'
       - '--cgroups-driver=systemd' # cgroupfs
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

