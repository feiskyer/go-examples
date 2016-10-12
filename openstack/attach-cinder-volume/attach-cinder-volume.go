package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	// "github.com/rackspace/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"net"
	"strings"

	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/pagination"
)

func get_iscsi_initiator() string {
	contents, err := ioutil.ReadFile("/etc/iscsi/initiatorname.iscsi")
	if err != nil {
		return ""
	}

	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "InitiatorName=") {
			return strings.Split(line, "=")[1]
		}
	}

	return ""
}

func get_local_ip() string {
	interfaces, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range interfaces {
		return addr.String()
	}

	return ""
}

func get_connector_properties(host, ip string) map[string]string {
	props := make(map[string]string)
	props["ip"] = ip
	props["host"] = host
	props["initiator"] = get_iscsi_initiator()
	props["multipath"] = "False"
	props["platform"] = "x86_64"
	props["os_type"] = "linux2"

	return props
}
func main() {
	// Option 1: Pass in the values yourself
	auth_opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://192.168.0.3:5000/v2.0",
		Username:         "admin",
		Password:         "admin",
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

	opts := networks.ListOpts{}
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

	// nova client
	novaClient, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil || novaClient == nil {
		fmt.Println("Unable to initialize nova client")
		os.Exit(1)
	}

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Failed to get Minion's hostname: ", err)
		os.Exit(1)
	}

	fmt.Println("Hostname is ", hostname)

	serverOpts := servers.ListOpts{}
	serverPager := servers.List(novaClient, serverOpts)
	err = serverPager.EachPage(func(page pagination.Page) (bool, error) {
		sList, err := servers.ExtractServers(page)
		if err != nil {
			return false, err
		}
		for _, server := range sList {
			fmt.Println("Server: ", server)
		}
		return true, nil
	})
	if err != nil {
		fmt.Println("List servers error: ", err)
		os.Exit(1)
	}

	cinderClient, err := openstack.NewBlockStorageV1(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		fmt.Println("Create cinder client error: ", err)
	}

	cinderOpts := volumes.ListOpts{}
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
	fmt.Println("Volume list: ", volumeList)

	// attach
	attachOpts := volumes.AttachOpts{
		MountPoint: "xvda",
		Mode:       "rw",
		HostName:   "localhost",
		//InstanceUUID: 	"UUID",
	}

	attachResult := volumes.Attach(cinderClient, volumeList[0].ID, attachOpts)
	if attachResult.Err != nil {
		fmt.Println("Attach volume failed,", attachResult.Err)
	} else {
		fmt.Println("Attach volume result: ", attachResult.Body)
	}

	connector := volumes.ConnectorOpts{
		IP:        "127.0.0.1",
		Host:      "localhost",
		Initiator: "iqn.1994-05.com.redhat:2c3a5fa29de",
	}
	connectionInfo := volumes.InitializeConnection(cinderClient, volumeList[0].ID, connector)
	if connectionInfo.Err != nil {
		fmt.Println("InitializeConnection volume failed,", connectionInfo.Err)
	} else {
		fmt.Println("InitializeConnection volume result: ", connectionInfo.Body)
	}

	// detach
	detachResult := volumes.Detach(cinderClient, volumeList[0].ID)
	if detachResult.Err != nil {
		fmt.Println("Detach volume failed,", detachResult.Err)
	} else {
		fmt.Println("Detach volume result: ", detachResult.Body)
	}
}
