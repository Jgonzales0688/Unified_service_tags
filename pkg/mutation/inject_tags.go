package mutation

import (
	"strings"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// injectTags is a container for the mutation injecting environment vars
type injectTags struct {
	Logger logrus.FieldLogger
}

// injectTags implements the podMutator interface
var _ podMutator = (*injectTags)(nil)

// Name returns the struct name
func (se injectTags) Name() string {
	return "inject_tags"
}

// Mutate returns a new mutated pod according to set label rules
func (se injectTags) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {
	mpod := pod.DeepCopy()

	//Add the desired tags to the labes
	labels := make(map[string]string)

	// Get the environment from the namespace
	env := getEnvFromNamespace(pod.Namespace)
	if env != "" {
		labels["tags.datadoghq.com/env"] = env
	}
	serviceName := getServiceNameFromImage(pod.Spec.Containers[0].Image)
	if serviceName != "" {
		labels["tags.datadoghq.com/service"] = serviceName
	}
	// Get the version from the docker image name
	version := getVersionFromImageName(pod.Spec.Containers[0].Image)
	if version != "" {
		labels["tags.datadoghq.com/version"] = version
	}

	// inject labels into the pod
	injectLabels(mpod, labels)

	return mpod, nil
}

// injectLabels injects labels into the pod
func injectLabels(pod *corev1.Pod, labels map[string]string) {
	if pod.ObjectMeta.Labels == nil {
		pod.ObjectMeta.Labels = make(map[string]string)
	}

	for k, v := range labels {
		pod.ObjectMeta.Labels[k] = v
	}
}

// getEnvFromNamespace exracts the environment form the namespace
func getEnvFromNamespace(namespace string) string {
	parts := strings.Split(namespace, "-")
	return parts[len(parts)-1]
}

// getServiceNameFromImage extracts the service name from the docker image
func getServiceNameFromImage(imageName string) string {
	parts := strings.Split(imageName, ":")
	if len(parts) > 0 {
		serviceNameParts := strings.Split(parts[0], ":")
		return serviceNameParts[len(serviceNameParts)-1]
	}
	return ""
}

// getVersionFromImageName extracts the version from the docker image name
func getVersionFromImageName(imageName string) string {
	parts := strings.Split(imageName, ":")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}
