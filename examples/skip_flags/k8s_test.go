/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package skip_funcs

import (
	"context"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestSkipFlags(t *testing.T) {
	podFeature := features.New("pod list").WithLabel("env", "prod").
		Assess("pods from kube-system", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			var pods corev1.PodList
			client, err := cfg.NewClient()
			if err != nil {
				t.Fatal(err)
			}
			err = client.Resources("kube-system").List(context.TODO(), &pods)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("found %d pods", len(pods.Items))
			if len(pods.Items) == 0 {
				t.Fatal("no pods in namespace kube-system")
			}
			return ctx
		}).Feature()

	// feature uses pre-generated namespace (see TestMain)
	depFeature := features.New("appsv1/deployment").WithLabel("env", "dev").
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			// insert a deployment
			deployment := newDeployment(cfg.Namespace(), "test-deployment", 1)
			client, err := cfg.NewClient()
			if err != nil {
				t.Fatal(err)
			}

			if err := client.Resources().Create(ctx, deployment); err != nil {
				t.Fatal(err)
			}
			time.Sleep(2 * time.Second)
			return ctx
		}).
		Assess("deployment creation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			var dep appsv1.Deployment
			client, err := cfg.NewClient()
			if err != nil {
				t.Fatal(err)
			}

			if err := client.Resources().Get(ctx, "test-deployment", cfg.Namespace(), &dep); err != nil {
				t.Fatal(err)
			}
			if &dep != nil {
				t.Logf("deployment found: %s", dep.Name)
			}

			return ctx
		}).Feature()

	test.Test(t, podFeature, depFeature)
}

func newDeployment(namespace string, name string, replicas int32) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: map[string]string{"app": "test-skip-label"}},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "test-skip-label"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "test-skip-label"}},
				Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "scratch", Image: "scratch"}}},
			},
		},
	}
}
