package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumeactions"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/volumes"
)

func main() {
	volumeID := flag.String("volume", "", "volume to detach")
	region := flag.String("region", "RegionOne", "region for operations")
	flag.Parse()

	if *volumeID == "" {
		fmt.Println("Must set volumeID")
		os.Exit(1)
	}

	auth_opts, err := openstack.AuthOptionsFromEnv()
	provider, err := openstack.AuthenticatedClient(auth_opts)
	if err != nil {
		fmt.Println("Authentication error: ", err)
		os.Exit(1)
	}

	cinderClient, err := openstack.NewBlockStorageV2(provider, gophercloud.EndpointOpts{
		Region: *region,
	})
	if err != nil {
		fmt.Println("Create cinder client error: ", err)
	}

	res, err := volumes.Get(cinderClient, *volumeID).Extract()
	if err != nil {
		fmt.Println("Get volume err: ", err)
		os.Exit(1)
	} else {
		fmt.Printf("Got volume: %#v\n", res)
	}

	for _, attach := range res.Attachments {
		// detach
		detachResult := volumeactions.Detach(cinderClient, *volumeID, volumeactions.DetachOpts{
			AttachmentID: attach.AttachmentID,
		})
		if detachResult.Err != nil {
			fmt.Println("Detach volume failed,", detachResult.Err)
			os.Exit(1)
		}
	}

	fmt.Println("success")
}
