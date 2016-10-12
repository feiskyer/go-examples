package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"net"
	"strings"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"

	"github.com/rackspace/gophercloud/openstack/blockstorage/v2/extensions/volumeactions"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v2/volumes"
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
	remove := flag.Bool("r", true, "remove all volumes")
	flag.Parse()

	// Option 1: Pass in the values yourself
	auth_opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://192.168.0.2:5000/v2.0",
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
		Size: 1,
		Name: "test",
	}
	_, err = volumes.Create(cinderClient, cOpts).Extract()
	if err != nil {
		fmt.Println("Create volume error: ", err)
		os.Exit(1)
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
	if err != nil {
		fmt.Println("List volumes error: ", err)
		os.Exit(1)
	}

	if len(volumeList) == 0 {
		fmt.Println("There is no volumes found")
		os.Exit(0)
	}

	res, err := volumes.Get(cinderClient, volumeList[0].ID).Extract()
	if err != nil {
		fmt.Println("Get volume err: ", err)
		os.Exit(1)
	} else {
		fmt.Printf("Got volume: %q\n", res.Name)
	}

	// attach
	attachOpts := volumeactions.AttachOpts{
		MountPoint: "/tmp/adfa/-ada/sdasdfs",
		Mode:       "rw",
		HostName:   hostname,
		//InstanceUUID: 	"UUID",
	}

	attachResult := volumeactions.Attach(cinderClient, volumeList[0].ID, attachOpts)
	if attachResult.Err != nil {
		fmt.Println("Attach volume failed,", attachResult.Err)
	}

	connector := volumeactions.ConnectorOpts{
		//IP:        get_local_ip(),
		Host:      hostname,
		Initiator: get_iscsi_initiator(),
	}
	connectionInfo := volumeactions.InitializeConnection(cinderClient, volumeList[0].ID, &connector)
	if connectionInfo.Err != nil {
		fmt.Println("InitializeConnection volume failed,", connectionInfo.Err)
	} else {
		fmt.Printf("InitializeConnection volume result: %q\n", connectionInfo.Body)
	}

	terminateResult := volumeactions.TerminateConnection(cinderClient, volumeList[0].ID, &connector)
	if terminateResult.Err != nil {
		fmt.Println("terminateResult.Err: ", terminateResult.Err)
	}

	// detach
	detachResult := volumeactions.Detach(cinderClient, volumeList[0].ID)
	if detachResult.Err != nil {
		fmt.Println("Detach volume failed,", detachResult.Err)
	}

	// delete
	if *remove {
		for _, v := range volumeList {
			volumeactions.Detach(cinderClient, v.ID)
			err = volumes.Delete(cinderClient, v.ID).ExtractErr()
			if err != nil {
				fmt.Printf("Delete volume %s error: %v\n", v.ID, err)
			}
		}
	}
}
