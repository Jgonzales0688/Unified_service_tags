package validation

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/heb-engineering/teams/platform-engineering/gke-hybrid-cloud/kon/apps/Mutating-Webhook/pkg/mutation"
	corev1 "k8s.io/api/core/v1"
)

// Validator is a container for mutation
type Validator struct {
	Logger  *logrus.Entry
	Mutator *mutation.Mutator // Mutator field for accessing mutations
}

// NewValidator returns an initialised instance of Validator
func NewValidator(logger *logrus.Entry, mutator *mutation.Mutator) *Validator {
	return &Validator{
		Logger:  logger,
		Mutator: mutator,
	}
}

// podValidators is an interface used to group functions mutating pods
type podValidator interface {
	Validate(*corev1.Pod) (validation, error)
	Name() string
}

type validation struct {
	Valid  bool
	Reason string
}

// ValidatePod validates a pod and applies mutations if necessary
func (v *Validator) ValidatePod(pod *corev1.Pod) (validation, error) {
	var podName string
	if pod.Name != "" {
		podName = pod.Name
	} else {
		if pod.ObjectMeta.GenerateName != "" {
			podName = pod.ObjectMeta.GenerateName
		}
	}
	log := v.Logger.WithField("pod_name", podName)
	log.Print("delete me")

	// List of all validations to be applied to the pod
	validations := []podValidator{
		labelValidator{Logger: v.Logger},
	}

	// Apply all validations
	for _, val := range validations {
		vp, err := val.Validate(pod)
		if err != nil {
			return validation{Valid: false, Reason: err.Error()}, err
		}
		if !vp.Valid {
			log.Println("Pod is not valid. Applying mutations...")
			// Apply mutations using the Mutator
			patchedPod, err := v.Mutator.MutatePodPatch(pod)
			if err != nil {
				return validation{Valid: false, Reason: err.Error()}, err
			}
			// You can apply the patched pod in your logic as needed
			log.Println("Patched Pod:", string(patchedPod))
			return validation{Valid: false, Reason: vp.Reason}, nil
		}
	}

	return validation{Valid: true, Reason: "valid pod"}, nil
}
