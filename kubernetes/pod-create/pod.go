package main

import (
	"fmt"
	"time"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/client/restclient"
	kubeclient "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
)

func NewKubeClient() (*kubeclient.Client, error) {
	client, err := kubeclient.New(&restclient.Config{
		Host: "192.168.0.3:8080",
	})

	if err != nil {
		return nil, err
	}
	return client, nil
}

func CreatePod() error {
	var (
		resp *api.Pod = &api.Pod{}
		w    watch.Interface
	)

	req := &api.Pod{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: api.ObjectMeta{
			Name: "test-pod",
		},
		Spec: api.PodSpec{
			Containers: []api.Container{
				{
					Name:  "nginx",
					Image: "nginx",
				},
			},
		},
	}
	cli, err := NewKubeClient()
	if err != nil {
		return err
	}

	if resp, err = cli.Pods("default").Create(req); err != nil {
		return err
	}

	fmt.Printf("Pod created: %s", resp)

	status := resp.Status
	if w, err = cli.Pods("default").Watch(api.ListOptions{
		Watch:           true,
		ResourceVersion: resp.ResourceVersion,
		FieldSelector:   fields.Set{"metadata.name": "test-pod"}.AsSelector(),
		LabelSelector:   labels.Everything(),
	}); err != nil {
		return err
	}

	func() {
		for {
			select {
			case events, ok := <-w.ResultChan():
				if !ok {
					return
				}
				resp = events.Object.(*api.Pod)
				fmt.Println("Pod status:", resp.Status.Phase)
				status = resp.Status
				if resp.Status.Phase != api.PodPending {
					w.Stop()
				}
			case <-time.After(10 * time.Second):
				fmt.Println("timeout to wait for pod active")
				w.Stop()
			}
		}
	}()
	if status.Phase != api.PodRunning {
		return fmt.Errorf("Pod is unavailable: %v", status.Phase)
	}
	return nil
}

func main() {
	err := CreatePod()
	if err != nil {
		fmt.Println("CreatePod error: ", err)
	}
}
