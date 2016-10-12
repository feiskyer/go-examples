package main

import (
	"log"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	runtimeApi "k8s.io/kubernetes/pkg/kubelet/api/v1alpha1/runtime"
)

var apiVersion = "0.1.0"

// dial creates a net.Conn by unix socket addr.
func dial(addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("unix", addr, timeout)
}

func main() {
	conn, err := grpc.Dial("/tmp/frakti.sock", grpc.WithInsecure(), grpc.WithDialer(dial))
	if err != nil {
		log.Fatalf("Connect remote runtime failed: %v\n", err)
	}

	client := runtimeApi.NewRuntimeServiceClient(conn)
	resp, err := client.Version(context.Background(), &runtimeApi.VersionRequest{
		Version: &apiVersion,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Response %#q\n", resp)
}
