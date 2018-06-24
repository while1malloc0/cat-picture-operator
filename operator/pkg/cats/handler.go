package cats

import (
	"context"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/while1malloc0/cat-picture-operator/operator/pkg/apis/cat/v1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct{}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1.CatPicture:
		if err := sdk.Create(newCatDeployment(o)); err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create cat deployment : %v", err)
			return err
		}

		if err := sdk.Create(newCatService(o)); err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create cat service : %v", err)
			return err
		}
	}

	return nil
}

const catContainerImage = "gcr.io/spartan-network-174104/cat-pic-operator-app"

// newCatDeployment creates a new deployment of the cat picture service
func newCatDeployment(cr *v1.CatPicture) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      "cat-deployment",
			Namespace: cr.Namespace,
		},

		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "cat-pictures",
				},
			},
			Replicas: &cr.Spec.Num,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "cat-pictures",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "cat-pictures",
							Image: catContainerImage,
							Ports: []corev1.ContainerPort{
								{
									HostPort:      8080,
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}
}

// newCatService creates a new service for the cat application
func newCatService(cr *v1.CatPicture) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"service": "cats",
			},
			Namespace: cr.Namespace,
			Name:      "cat-pictures",
		},
		Spec: corev1.ServiceSpec{
			Type: "LoadBalancer",
			Selector: map[string]string{
				"app": "cat-pictures",
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.IntOrString{IntVal: 8080},
					Protocol:   "TCP",
				},
			},
		},
	}
}
