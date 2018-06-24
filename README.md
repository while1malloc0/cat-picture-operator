# Cat Picture Operator
A tutorial for creating a new custom resource type called a CatPicture. This will be acted upon by an operator written with the operator framework.

## Requirements

* A Go development environment
* A working Kubernetes cluster
* THe operator-sdk installed

## Build Steps

### Create the project scaffolding

```bash
operator-sdk new operator --api-version=cat.example.com/v1 --kind=CatPicture --skip-git-init
```

### Define the spec for the CatPicture type

```go
// operator/pkg/apis/cat/v1/types.go
type CatPictureSpec struct {
	// The number of containers to have running in the ReplicaSet
	Num int32 `json:"num,omitempty"`
}
```

### Update the business logic for the Handler

Rename the `stub` package to `cats` by renaming the folder and replacing the reference in `cats/handler.go`. Make sure to rename the references in `operator/cmd/operator/main.go`

Remove the `newBusyBox` method. We'll want to create two things: a Service and a Deployment. Create the deployment first. We'll need import Deployments from the Kubernetes client-go library.

```go
import (
    // ... a bunch of stuff
    appsv1 "k8s.io/api/apps/v1"
)
```

```go
// operator/pkg/cats/handler.go

// catContainerImage is the name of the container image
const catContainerImage = "fill-in-your-location"

// newCatDeployment creates a new deployment of the cat picture service
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
```

And then the service

```go
// operator/pkg/cats/handler.go

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
```

Plug the two creation functions into the handler

```go
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
```

### Build the operator

```bash
# Replace IMAGE_REPO with what's being pushed to.
# For google container images, use gcr.io/[YOUR_PROJECT_ID]

# From within the operator folder
operator-sdk build [IMAGE_REPO]/cat-picture-operator

docker push [IMAGE_REPO]/cat-picture-operator:latest
```

### Deploy the operator

```bash
kubectl apply -f ./operator/deploy/crd.yaml
# Optional, but recommended
kubectl apply -f ./operator/deploy/rbac.yaml
# This file is generated by the build step above
kubectl apply -f ./operator/deploy/operator.yaml
```

### Deploy the example resource

```bash
kubectl apply -f ./operator/deploy/cr.yaml
```

### Get some cat picture

Get the cat-pictures service and copy the ExternalIP

```bash
kubectl get svc
```

Copy the IP into your browser and enjoy the cats
