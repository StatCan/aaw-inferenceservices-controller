package main

import (
	"github.com/StatCan/inferenceservices-controller/cmd"

	// Import the Azure auth method for Kubernetes
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
)

func main() {
	cmd.Execute()
}
