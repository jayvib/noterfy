package cli

import "github.com/spf13/cobra"

// RootCmd is the root cli command for the noterfy.
var RootCmd = &cobra.Command{
	Use:   "noterfy",
	Short: "noterfy is a CLI application for noterfy",
}
