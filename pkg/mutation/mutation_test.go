package mutation

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMutatePodPatch(t *testing.T) {
	m := NewMutator(logger())
	got, err := m.MutatePodPatch(pod())
	if err != nil {
		t.Fatal(err)
	}

	p := patch()
	g := string(got)
	assert.JSONEq(t, p, g)
}

func pod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "lifespan",
			Labels: map[string]string{
				"acme.com/lifespan-requested": "7",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "lifespan",
				Image: "busybox",
			}},
		},
	}
}

func patch() string {
	patch := `[
		{"op":"add","path":"/metadata/labels/tags.datadoghq.com~1env","value":""},
		{"op":"add","path":"/metadata/labels/tags.datadoghq.com~1service","value":"busybox"},
		{"op":"add","path":"/metadata/labels/tags.datadoghq.com~1version","value":""}
]`

	patch = strings.ReplaceAll(patch, "\n", "")
	patch = strings.ReplaceAll(patch, "\t", "")
	patch = strings.ReplaceAll(patch, " ", "")

	return patch
}

func logger() *logrus.Entry {
	mute := logrus.StandardLogger()
	mute.Out = ioutil.Discard
	return mute.WithField("logger", "test")
}
