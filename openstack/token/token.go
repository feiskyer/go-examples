package main

import (
	"fmt"
	"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/identity/v2/tenants"
	"github.com/rackspace/gophercloud/openstack/identity/v2/tokens"
	"github.com/rackspace/gophercloud/pagination"
)

func main() {
	// Option 1: Pass in the values yourself
	auth_opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://192.168.0.171:5000/v2.0",
		Username:         "admin",
		Password:         "admin",
		AllowReauth:      true,
	}
	//TenantName:       "admin",

	// Option 2: Use a utility function to retrieve all your environment variables
	//auth_opts, err := openstack.AuthOptionsFromEnv()

	provider, err := openstack.AuthenticatedClient(auth_opts)
	if err != nil {
		fmt.Println("Authentication error: ", err)
		os.Exit(1)
	}

	identity := openstack.NewIdentityV2(provider)
	if identity == nil {
		fmt.Println("Failed to find identity endpoint")
	}

	topts := tenants.ListOpts{}
	tpager := tenants.List(identity, &topts)
	p, _ := tpager.AllPages()
	fmt.Println("Body: ", p.GetBody())

	terr := tpager.EachPage(func(page pagination.Page) (bool, error) {
		tenantList, err := tenants.ExtractTenants(page)
		if err != nil {
			return false, err
		}

		for _, t := range tenantList {

			fmt.Printf("Tenant name: %s, ID: %s\n", t.Name, t.ID)
		}

		return true, nil
	})
	if terr != nil {
		fmt.Println("Error %s", terr)
	}

	ttopts := tokens.AuthOptions{
		IdentityEndpoint: "http://192.168.0.3:5000/v2.0",
		Username:         "admin",
		Password:         "admin",
		TenantID:         "admin",
	}
	token, err := tokens.Create(identity, ttopts).Extract()
	if err != nil {
		fmt.Println("Failed: %v", err)
	} else {
		fmt.Println("Auth success")
	}

}
