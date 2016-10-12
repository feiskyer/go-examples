package main

import (
	"fmt"
	"os"

	"github.com/fsouza/go-dockerclient"
)

func main() {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", img.Created)
		fmt.Println("Size: ", img.Size)
		fmt.Println("VirtualSize: ", img.VirtualSize)
		fmt.Println("ParentId: ", img.ParentID)
	}

	fmt.Println()

	opts := docker.LogsOptions{
		Container:    "a739618663b59178223f4e975a93f035c29ebe48d1e2dc0bc5c4063648a02c15",
		Stdout:       true,
		Stderr:       true,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Timestamps:   false,
		RawTerminal:  false,
		Follow:       true,
	}

	err := client.Logs(opts)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return
}
