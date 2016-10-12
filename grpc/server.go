// Reference http://www.grpc.io/docs/tutorials/basic/go.html
package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct{}

func randString(strSize int, randType string) string {
	var dictionary string
	if randType == "alphanum" {
		dictionary = "0123456789abcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

func (s *server) ContainerLogs(req *ContainerLogsRequest, stream RemoteRuntime_ContainerLogsServer) error {
	container := req.Container
	reqType := req.Type
	log.Printf("Recv request from %s", container)
	for i := 0; i < 10; i++ {
		stream.Send(&ContanerLogsResponse{
			Log: randString(4096, "alphanum"),
		})
	}

	if reqType != int32(0) {
		stream.Send(&ContanerLogsResponse{
			Error: "error response",
		})
		return fmt.Errorf("Error response")
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	RegisterRemoteRuntimeServer(s, &server{})
	s.Serve(lis)
}
