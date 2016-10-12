// Reference http://www.grpc.io/docs/tutorials/basic/go.html
package main

import (
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// test
const (
	address     = "127.0.0.1:50051"
	defaultName = "world"
)

// test
func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := NewRemoteRuntimeClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = strings.Join(os.Args[1:], " ")
	}

	stream, err := c.ContainerLogs(
		context.Background(),
		&ContainerLogsRequest{
			Container: name,
			Filter:    "test",
			Type:      0,
		},
	)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// Read Done
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to recv: %v", err)
				return
			}
			log.Printf("Recv: %s", in.Log)
		}
	}()

	stream.CloseSend()
	<-waitc
}
