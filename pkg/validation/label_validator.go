package validation

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// labelValidator is a container for validating the name of pods
type labelValidator struct {
	Logger logrus.FieldLogger
}

// labelValidator implements the podValidator interface
var _ podValidator = (*labelValidator)(nil)

// Name returns the name of labelValidator
func (lv labelValidator) Name() string {
	return "name_validator"
}

// Validate inspects the labels of a given pod and returns validation.
// The returned validation is only valid if the pod label contains the required labels
func (lv labelValidator) Validate(pod *corev1.Pod) (validation, error) {
	requiredLabels := []string{
		"tags.datadoghq.com/env",
		"tags.datadoghq.com/service",
		"tags.datadoghq.com/version",
	}

	missingLabels := []string{}
	for _, label := range requiredLabels {
		if _, exists := pod.Labels[label]; !exists {
			missingLabels = append(missingLabels, label)
		}
	}

	if len(missingLabels) > 0 {
		v := validation{
			Valid:  false,
			Reason: fmt.Sprintf("missing labels: %s", strings.Join(missingLabels, ", ")),
		}
		return v, nil
	}
	return validation{Valid: true, Reason: "valid labels"}, nil

}
