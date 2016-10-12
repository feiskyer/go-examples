/*
	Query a regular user's tenant list
*/
package main

import (
	"fmt"
	"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/identity/v2/tenants"
	"github.com/rackspace/gophercloud/pagination"
)

func main() {
	auth_opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://192.168.0.2:5000/v2.0",
		Username:         "demouser",
		Password:         "demo",
		AllowReauth:      true,
	}

	provider, err := openstack.AuthenticatedClient(auth_opts)
	if err != nil {
		fmt.Println("Authentication error: ", err)
		os.Exit(1)
	}

	identity := openstack.NewIdentityV2(provider)

	opts := tenants.ListOpts{Limit: 5}
	pager := tenants.List(identity, &opts)
	userTenants := make([]tenants.Tenant, 0, 1)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		tenantList, err := tenants.ExtractTenants(page)
		if err != nil {
			return false, nil
		}
		for _, t := range tenantList {
			userTenants = append(userTenants, t)
		}

		return true, nil
	})
	if err != nil {
		fmt.Println("Tenat list error: ", err)
		os.Exit(1)
	}

	fmt.Println("User Tenants: ", userTenants)
}
