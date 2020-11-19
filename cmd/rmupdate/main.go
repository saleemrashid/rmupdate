package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: os.Args[0],
	Long: `Utility for fetching software updates from the reMarkable update server.

The reMarkable uses the CoreOS "update_engine" daemon to fetch update payloads from an Omaha server. An update payload contains "installation operations", which describe how to install the update to the device's eMMC.
This utility can fetch and verify these payloads, then reconstruct the installation operations into a root filesystem image.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
