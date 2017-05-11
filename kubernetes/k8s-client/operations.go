package main

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/resource"
	vapi "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

const namespace string = "default"

// operation represents a Kubernetes operation.
type operation interface {
	Do(c *client.Client)
}

type versionOperation struct{}

func (op *versionOperation) Do(c *client.Client) {
	info, err := c.Discovery().ServerVersion()
	if err != nil {
		logger.Fatalf("failed to retrieve server API version: %s\n", err)
	}

	logger.Printf("server API version information: %s\n", info)
}

type deployOperation struct {
	image string
	name  string
	port  int
}

func (op *deployOperation) Do(c *client.Client) {
	appName := op.name
	port := op.port
	deploy := c.Extensions().Deployments(namespace)

	// Define Deployments spec.
	d := &extensions.Deployment{
		TypeMeta: vapi.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: api.ObjectMeta{
			Name: appName,
		},
		Spec: extensions.DeploymentSpec{
			Replicas: 1,
			Template: api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Name:   appName,
					Labels: map[string]string{"app": appName},
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						{
							Name:  appName,
							Image: op.image,
							Ports: []api.ContainerPort{
								{ContainerPort: port, Protocol: api.ProtocolTCP},
							},
							Resources: api.ResourceRequirements{
								Limits: api.ResourceList{
									api.ResourceCPU:    resource.MustParse("100m"),
									api.ResourceMemory: resource.MustParse("256Mi"),
								},
							},
							ImagePullPolicy: api.PullIfNotPresent,
						},
					},
					RestartPolicy: api.RestartPolicyAlways,
					DNSPolicy:     api.DNSClusterFirst,
				},
			},
		},
	}

	// Implement update-or-create logic (i.e., `kubectl apply`).
	depl, err := deploy.Update(d)
	switch {
	case err == nil:
		op.logSuccess(false, depl)
	case !errors.IsNotFound(err):
		logger.Fatalf("failed to update deployment controller: %s", err)
	default:
		depl, err = deploy.Create(d)
		if err != nil {
			logger.Fatalf("failed to create deployment controller: %s", err)
		}
		op.logSuccess(true, depl)
	}
}

func (op *deployOperation) logSuccess(isCreate bool, depl *extensions.Deployment) {
	mode := "updated"
	if isCreate {
		mode = "created"
	}

	logger.Printf("deployment controller %s (received definition: %+v)\n", mode, depl)
}
