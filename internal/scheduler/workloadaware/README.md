# WorkloadAware Scheduler Plugin

The WorkloadAware plugin is designed to make scheduling and resource binding decisions based on the workload type of a Pod. It classifies workloads into four categories:

- **MemoryBound:** Workloads with execution time that depends on the available memory bandwidth.
- **CPUBound:** Workloads with execution time that depends on the available CPU resources.
- **IOBound:** Workloads that have threads with high IO wait time.
- **BestEffort:** Workloads that place every thread on the same logical CPU (oversubscription).

## Resource Allocation Based on Workload Type

- **MemoryBound:** Threads are placed on different memory nodes (sockets).
- **CPUBound:** Threads are placed on different, non-utilized cores, preferably on the same socket.
- **IOBound:** Threads are placed on the same physical core, or more cores are utilized if needed.
- **BestEffort:** Every thread is placed on the same logical CPU.

## Policies

- **MaximumUtilization:** The plugin tries to maximize the utilization of the resources.
- **Balanced:** The plugin places the Pods in a balanced way across the cluster. 

## Extra Features

- **PhysicalCores:** Use only physical cores for scheduling.
- **MemoryBoundExclusiveSockets:** Allocate memory nodes (sockets) exclusively for MemoryBound workloads.
- **BestEffortSharedCPUs:** Allow BestEffort workloads to share logical CPUs.

## Installation

Apply the following KubeSchedulerConfiguration to enable the plugin:

```yaml
apiVersion: kubescheduler.config.k8s.io/v1
kind: KubeSchedulerConfiguration
leaderElection:
  leaderElect: false
profiles:
  - schedulerName: maestro
    plugins:
      preFilter:
        enabled:
          - name: WorkloadAware
      filter:
        enabled:
          - name: TaintToleration
          - name: WorkloadAware
      score:
        enabled:
          - name: WorkloadAware
      bind:
        enabled:
          - name: WorkloadAware
        disabled:
          - name: '*'
    pluginConfig:
      - name: WorkloadAware
        args:
          policy: MaximumUtilization
          features:
            - PhysicalCores
            - BestEffortSharedCPUs 
#           - MemoryBoundExclusiveSockets 
```