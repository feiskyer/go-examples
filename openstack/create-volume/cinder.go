package main

import (
	"fmt"
	"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v2/volumes"
	"github.com/rackspace/gophercloud/pagination"
)

func main() {
	// Option 1: Pass in the values yourself
	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://192.168.0.3:5000/v2.0",
		Username:         "admin",
		Password:         "admin",
		TenantName:       "admin",
	}

	// Option 2: Use a utility function to retrieve all your environment variables
	//authOpts, err := openstack.AuthOptionsFromEnv()

	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		fmt.Println("Authentication error: ", err)
		os.Exit(1)
	}

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Failed to get Minion's hostname: ", err)
		os.Exit(1)
	}

	fmt.Println("Hostname is ", hostname)

	cinderClient, err := openstack.NewBlockStorageV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		fmt.Println("Create cinder client error: ", err)
	}

	cOpts := volumes.CreateOpts{
		Size:     1,
		TenantID: "f0df9b076efb4ca198f328fc6e65097b",
		Name:     "test",
	}
	testVol, err := volumes.Create(cinderClient, cOpts).Extract()
	if err != nil {
		fmt.Println("Create volume error: ", err)
		os.Exit(1)
	}

	cinderOpts := volumes.ListOpts{AllTenants: true}
	vlistPager := volumes.List(cinderClient, cinderOpts)
	var volumeList []volumes.Volume
	err = vlistPager.EachPage(func(page pagination.Page) (bool, error) {
		vList, err := volumes.ExtractVolumes(page)

		if err != nil {
			return false, err
		}

		for _, v := range vList {
			volumeList = append(volumeList, v)
		}
		return true, nil
	})
	if err != nil {
		fmt.Println("List volumes error: ", err)
		os.Exit(1)
	}

	if len(volumeList) == 0 {
		fmt.Println("There is no volumes found")
		os.Exit(0)
	}

	for _, vol := range volumeList {
		res, err := volumes.Get(cinderClient, vol.ID).Extract()
		if err != nil {
			fmt.Println("Get volume err: ", err)
			os.Exit(1)
		} else {
			fmt.Printf("Got volume: %q\tTenant: %v\n", res.Name, res.TenantID)
		}
	}

	volumes.Delete(cinderClient, testVol.ID)
}
