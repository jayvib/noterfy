package cli

import "github.com/spf13/cobra"

// RootCmd is the root cli command for the noteapp.
var RootCmd = &cobra.Command{
	Use:   "noteapp",
	Short: "noteapp is a CLI application for noteapp",
}
