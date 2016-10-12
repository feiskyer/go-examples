package main

import (
	"fmt"
	"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

func main() {
	// Option 1: Pass in the values yourself
	auth_opts := gophercloud.AuthOptions{
		IdentityEndpoint: "https://192.168.0.3:5000/v2.0",
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

	client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		fmt.Println("Err: ", err)
		os.Exit(1)
	}

	opts := floatingips.CreateOpts{
		FloatingNetworkID: "06f456ea-5c72-458b-9875-d59f6a81f9b5",
		FloatingIP:        "1.2.3.4",
	}
	fip, err := floatingips.Create(client, opts).Extract()
	if err != nil {
		fmt.Println("Create openstack flaotingip failed: ", err)
	} else {
		fmt.Println("FloatingIP created: ", fip)
	}
}
