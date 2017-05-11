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

func NewClient() (*kubeclient.Client, error) {
	client, err := kubeclient.New(&restclient.Config{
		Host: "192.168.0.3:8080",
	})

	if err != nil {
		return nil, err
	}
	return client, nil
}

func CreateNetwork() error {
	var (
		netCreateResp *api.Network = &api.Network{}
		w             watch.Interface
	)

	netCreateReq := &api.Network{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Network",
			APIVersion: "v1",
		},
		ObjectMeta: api.ObjectMeta{
			Name: "test-network",
		},
		Spec: api.NetworkSpec{
			TenantID: "test-tenant",
			Subnets: map[string]api.Subnet{
				"subnet1": {
					CIDR:    "192.176.0.0/24",
					Gateway: "192.176.0.1",
				},
			},
		},
	}
	cli, err := NewClient()
	if err != nil {
		return err
	}

	if netCreateResp, err = cli.Networks().Create(netCreateReq); err != nil {
		return err
	}

	fmt.Printf("Network created: %s", netCreateResp)

	status := netCreateResp.Status
	if w, err = cli.Networks().Watch(api.ListOptions{
		Watch:         true,
		FieldSelector: fields.Set{"metadata.name": "test-network"}.AsSelector(),
		LabelSelector: labels.Everything(),
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
				netCreateResp = events.Object.(*api.Network)
				status = netCreateResp.Status
				if status.Phase != api.NetworkInitializing && status.Phase != api.NetworkPending {
					w.Stop()
				}
			case <-time.After(5 * time.Second):
				fmt.Println("timeout to wait for network active")
				w.Stop()
			}
		}
	}()

	if status.Phase != api.NetworkActive {
		return fmt.Errorf("network is unavailable: %v", status.Phase)
	}
	return nil
}

func main() {
	err := CreateNetwork()
	fmt.Println(err)
}
