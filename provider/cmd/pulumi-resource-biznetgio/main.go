package main

import (
	"context"
	"fmt"
	"os"

	"github.com/biznetgio/pulumi-biznetgio/provider"
)

func main() {
	err := provider.Provider().Run(context.Background(), provider.Name, provider.Version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}
