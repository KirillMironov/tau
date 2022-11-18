package protoconv

import (
	"github.com/KirillMironov/tau/api"
	"github.com/KirillMironov/tau/resources"
)

func PodFromProto(pod *api.Pod) resources.Pod {
	return resources.Pod{
		Name:       pod.Name,
		Containers: ContainersFromProto(pod.Containers),
	}
}

func PodToProto(pod resources.Pod) *api.Pod {
	return &api.Pod{
		Name:       pod.Name,
		Containers: ContainersToProto(pod.Containers),
	}
}
