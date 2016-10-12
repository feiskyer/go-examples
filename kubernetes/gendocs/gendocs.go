package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/kubernetes/cmd/genutils"
	"k8s.io/kubernetes/cmd/kube-apiserver/app"
)

func main() {
	// use os.Args instead of "flags" because "flags" will mess up the man pages!
	path := "docs/"
	if len(os.Args) == 2 {
		path = os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [output directory]\n", os.Args[0])
		os.Exit(1)
	}

	outDir, err := genutils.OutDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get output directory: %v\n", err)
		os.Exit(1)
	}

	// Set environment variables used so the output is consistent,
	// regardless of where we run.
	os.Setenv("HOME", "/home/username")
	s := app.NewAPIServer()
	s.AddFlags(pflag.CommandLine)
	kubectl := &cobra.Command{
		Use:  "kube-apiserver",
		Long: "The Kubernetes API server validates and configures data for the api objects which include pods, services, replicationcontrollers, and others. The API Server services REST operations and provides the frontend to the cluster's shared state through which all other components interact.",
		Run: func(cmd *cobra.Command, args []string) {
			s.Run(pflag.CommandLine.Args())
		},
	}
	cobra.GenMarkdownTree(kubectl, outDir)
}
