package main

import (
	"fmt"
	"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/portsbinding"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/openstack/networking/v2/ports"
	"github.com/rackspace/gophercloud/pagination"
)

func main() {
	// Option 1: Pass in the values yourself
	auth_opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://192.168.0.3:5000/v2.0",
		Username:         "admin",
		Password:         "admin",
		AllowReauth:      true,
		TenantName:       "admin",
	}

	// Option 2: Use a utility function to retrieve all your environment variables
	//auth_opts, err := openstack.AuthOptionsFromEnv()

	provider, err := openstack.AuthenticatedClient(auth_opts)
	if err != nil {
		fmt.Println("Authentication error: ", err)
		os.Exit(1)
	}

	// Compute service client; which can be created like so:
	client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	opts := networks.ListOpts{Name: "aa"}
	pager := networks.List(client, opts)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		networkList, e := networks.ExtractNetworks(page)
		for _, f := range networkList {
			// "f" will be a flavors.Flavor
			fmt.Println("OpenStack networks: ", f)
		}
		return true, e
	})
	if err != nil {
		fmt.Println("Error ", err)
		os.Exit(1)
	}

	options := portsbinding.CreateOpts{
		Parent: ports.CreateOpts{
			Name:      "private-port",
			NetworkID: "c34e8326-933d-4fed-bc8f-1c5c05e0eac2",
		},
		HostID:   "HOST1",
		VNICType: "normal",
	}
	port, err := portsbinding.Create(client, options).Extract()
	if err != nil {
		fmt.Println("Error ", err)
		os.Exit(1)
	}
	fmt.Printf("Port created: %v\n", port)

}
