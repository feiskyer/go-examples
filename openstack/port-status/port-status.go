package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/networking/v2/ports"
	"github.com/rackspace/gophercloud/pagination"
)

const (
	portID = "888c106c-dbd2-46de-89c8-24c9e5dc2729"
)

func main() {
	auth_opts, err := openstack.AuthOptionsFromEnv()
	provider, err := openstack.AuthenticatedClient(auth_opts)
	if err != nil {
		fmt.Println("Authentication error: ", err)
		os.Exit(1)
	}

	client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		fmt.Println("Get neutron client error:", err)
		os.Exit(1)
	}

	for i := 0; i < 100; i++ {
		p, err := ports.Get(client, portID).Extract()
		if err != nil {
			fmt.Println("Get neutron port error:", err)
			os.Exit(1)
		}
		fmt.Printf("Port %q status %q\n", p.ID, p.Status)
		time.Sleep(100 * time.Millisecond)
	}

	var results []ports.Port
	listOpts := ports.ListOpts{ID: portID}

	portPager := ports.List(client, listOpts)
	err = portPager.EachPage(func(page pagination.Page) (bool, error) {
		portList, err := ports.ExtractPorts(page)
		if err != nil {
			fmt.Println("Get openstack ports error: %v", err)
			return false, err
		}

		for _, port := range portList {
			results = append(results, port)
		}

		return true, err
	})

	if err != nil {
		fmt.Printf("Get ports error: %s", err)
		os.Exit(1)
	}

	for _, p := range results {
		fmt.Println(p.ID, p.Status)
	}
}
