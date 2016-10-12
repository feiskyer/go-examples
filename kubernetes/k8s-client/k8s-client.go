package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

const (
	kubernetesHost string = "192.168.0.3:8080"
)

var logger *log.Logger

func init() {
	// Ignore timestamps.
	logger = log.New(os.Stdout, "", 0)
}

func usage(msg string) {
	logger.Printf("%s\n\n", msg)
	fmt.Fprintf(os.Stderr, "usage: %s version | deploy", path.Base(os.Args[0]))
	os.Exit(1)
}

func main() {
	// Skip 0-th argument containing the binary's name.
	args := os.Args[1:]
	if len(args) < 1 {
		usage("insufficient number of parameters")
	}

	opName := args[0]
	var op operation

	switch opName {
	case "version":
		op = &versionOperation{}
	case "deploy":
		op = &deployOperation{
			image: "nginx:latest",
			name:  "nginx",
			port:  80,
		}
	default:
		usage(fmt.Sprintf("unknown operation: %s", opName))
	}

	config := &restclient.Config{Host: kubernetesHost}
	c, err := client.New(config)
	if err != nil {
		logger.Fatalf("could not connect to Kubernetes API: %s", err)
	}

	op.Do(c)
}
