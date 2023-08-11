package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLabelValidator(t *testing.T) {
	want := validation{Valid: true, Reason: "valid labels"}

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

	lv := labelValidator{}
	got, err := lv.Validate(pod)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}
