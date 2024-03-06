package cpupinning

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// ParsePodInfo extracts relevant information from a Pod to create a cpupinning.Pod object
func ParsePodInfo(pod *v1.Pod) *Pod {
	p := &Pod{
		Id:         string(pod.ObjectMeta.UID),
		Name:       pod.Name,
		Namespace:  pod.Namespace,
		Containers: nil,
	}

	containers := make([]*Container, 0)
	for _, containerStatus := range pod.Status.ContainerStatuses {
		containers = append(containers, &Container{
			Id:        containerStatus.ContainerID,
			Name:      containerStatus.Name,
			Resources: ParseContainerResources(containerStatus.Name, pod),
		})
	}

	p.Containers = containers

	return p
}

// ParseContainerResources extracts resource information from a container
func ParseContainerResources(containerName string, pod *v1.Pod) *ResourceInfo {
	resources := &ResourceInfo{}
	for _, container := range pod.Spec.Containers {
		if container.Name == containerName {
			requestCpus := container.Resources.Requests.Cpu()
			limitCpus := container.Resources.Limits.Cpu()
			requestMemory := container.Resources.Requests.Memory()
			limitMemory := container.Resources.Limits.Memory()
			if requestCpus == nil {
				requestCpus = resource.NewMilliQuantity(0, resource.DecimalSI)
			}
			if limitCpus == nil {
				limitCpus = resource.NewMilliQuantity(0, resource.DecimalSI)
			}
			resources = &ResourceInfo{
				RequestedCpus:   int32(requestCpus.MilliValue()),
				LimitCpus:       int32(limitCpus.MilliValue()),
				RequestedMemory: requestMemory.String(),
				LimitMemory:     limitMemory.String(),
			}
			return resources
		}
	}
	return resources
}
