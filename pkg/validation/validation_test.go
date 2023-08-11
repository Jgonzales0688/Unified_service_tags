package validation

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// Import from the mutation package
	"gitlab.com/heb-engineering/teams/platform-engineering/gke-hybrid-cloud/kon/apps/Mutating-Webhook/pkg/mutation"
)

func Test_LabelValidator(t *testing.T) {
	want := validation{Valid: true, Reason: "valid pod"}

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
			Labels: map[string]string{
				"tags.datadoghq.com/env":     "kondev",
				"tags.datadoghq.com/service": "example-service:1.0.0",
				"tags.datadoghq.com/version": "1.0.0",
			},
		},
	}

	logger := logrus.New().WithField("test", "validation")
	mutator := mutation.NewMutator(logger)
	validator := NewValidator(logger, mutator)

	got, err := validator.ValidatePod(pod)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}
