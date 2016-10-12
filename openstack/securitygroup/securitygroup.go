package main

import (
	"fmt"
	"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"

	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/security/groups"
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

	client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	if err != nil {
		fmt.Println("NewNetworkV2 Error", err)
		os.Exit(1)
	}

	var securitygroup *groups.SecGroup
	opts := groups.ListOpts{
		Name: "testtt",
	}
	pager := groups.List(client, opts)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		sg, err := groups.ExtractGroups(page)
		if err != nil {
			fmt.Println("Get openstack securitygroups error: %v", err)
			return false, err
		}

		if len(sg) > 0 {
			securitygroup = &sg[0]
		}

		return true, err
	})
	if err != nil {
		fmt.Println("Error %v", err)
		os.Exit(1)
	}

	// If securitygroup doesn't exist, create a new one
	if securitygroup == nil {
		securitygroup, err = groups.Create(client, groups.CreateOpts{
			Name:     "testttt",
			TenantID: "56064caf8d9b4f7aabab297c7d5409dc",
		}).Extract()

		if err != nil {
			fmt.Println("Error %v", err)
			os.Exit(1)
		}
	}

	fmt.Println("SecGroup %s created", securitygroup.ID)
}
