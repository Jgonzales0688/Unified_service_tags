package mutation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInjectTagsMutate(t *testing.T) {
	want := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test",
			Namespace: "example-namespace",
			Labels: map[string]string{
				"tags.datadoghq.com/env":     "example-namespace",
				"tags.datadoghq.com/service": "example-service",
				"tags.datadoghq.com/version": "1.0.0",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "test",
				Image: "example-service:1.0.0",
			}},
			InitContainers: []corev1.Container{{
				Name: "inittest",
			}},
		},
	}

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test",
			Namespace: "example-namespace",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "test",
				Image: "example-service:1.0.0",
			}},
			InitContainers: []corev1.Container{{
				Name: "inittest",
			}},
		},
	}

	got, err := injectTags{Logger: logger()}.Mutate(pod)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}
