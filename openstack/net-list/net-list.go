package main

import (
	"fmt"
	"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/identity/v2/tenants"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/pagination"
)

func main() {
	auth_opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://192.168.0.3:5000/v2.0",
		Username:         "admin",
		Password:         "admin",
		TenantName:       "admin",
		AllowReauth:      true,
	}

	provider, err := openstack.AuthenticatedClient(auth_opts)
	if err != nil {
		fmt.Println("Authentication error: ", err)
		os.Exit(1)
	}

	identity, err := openstack.NewIdentityAdminV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		fmt.Println("Failed to find identity endpoint: %s", err)
		os.Exit(1)
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

	// Neutron
	client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})

	opts := networks.ListOpts{}
	pager := networks.List(client, opts)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		networkList, e := networks.ExtractNetworks(page)
		for _, f := range networkList {
			// "f" will be a flavors.Flavor
			fmt.Println("OpenStack networks: ", f.Name, f.Subnets)
		}
		return true, e
	})
	if err != nil {
		fmt.Println("Error ", err)
		os.Exit(1)
	}
}
