package main

import (
	"flag"
	"fmt"
	"os"
)

const server = "127.0.0.1:22318"

func main() {
	flag.Parse()

	cl, err := NewHyperClient(server)
	if err != nil {
		fmt.Println(err)
		return
	}

	containerList, err := cl.GetContainerList(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(containerList) == 0 {
		return
	}

	err = cl.GetContainerLogs(containerList[0].ContainerID, os.Stdout)
	if err != nil {
		fmt.Println(err)
		return
	}

}
