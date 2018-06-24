# Cat Picture Operator
We'll be creating a new custom resource type called a CatPicture. This will be acted upon by an operator written with the operator framework.

## Requirements

* A Go development environment
* A working Kubernetes cluster
* THe operator-sdk installed


## Build steps

### Create the project scaffolding

```bash
operator-sdk new operator --api-version=cat.example.com/v1 --kind=CatPicture --skip-git-init
```

### Define the spec for the CatPicture type

```go
// operator/pkg/apis/cat/v1/types.go
type CatPictureSpec struct {
	// The number of containers to have running in the ReplicaSet
	Num int
}
```

### Update the business logic for the Handler

Rename the `stub` package to `cats` by renaming the folder and replacing the reference in `cats/handler.go`.

Remove the `newBusyBox` method. We'll want to create two things: a Service and a Deployment. Create the deployment first

```go
// operator/pkg/cats/handler.go

// catContainerImage is the name of the container image
const catContainerImage = "fill-in-your-location"

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
		},
		Spec: corev1.ServiceSpec{
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