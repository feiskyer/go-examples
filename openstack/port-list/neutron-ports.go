package main

import (
	"fmt"
	"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/openstack/networking/v2/ports"
	"github.com/rackspace/gophercloud/pagination"
)

func main() {
	auth_opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://192.168.0.3:5000/v2.0",
		Username:         "admin",
		TenantName:       "admin",
		Password:         "admin",
		AllowReauth:      true,
	}

	provider, err := openstack.AuthenticatedClient(auth_opts)
	if err != nil {
		fmt.Println("Authentication error: ", err)
		os.Exit(1)
	}

	client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: "us-west-1.cell-01",
	})
	if err != nil {
		fmt.Println("Get neutron client error:", err)
		os.Exit(1)
	}

	network_id := os.Getenv("NETWORK")
	opts := networks.ListOpts{ID: network_id}
	pager := networks.List(client, opts)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		networkList, e := networks.ExtractNetworks(page)
		if len(networkList) != 1 {
			return false, fmt.Errorf("Can not find network")
		}
		return true, e
	})
	if err != nil {
		fmt.Println("Error ", err)
		os.Exit(1)
	}

	var results []ports.Port
	listOpts := ports.ListOpts{
		NetworkID:   network_id,
		DeviceOwner: "network:dhcp",
	}

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
		fmt.Println(p.ID, p.FixedIPs[0].IPAddress)
	}
}
