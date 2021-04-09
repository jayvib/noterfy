package cli

import "github.com/spf13/cobra"

func init() {
	Cmd.AddCommand(UtilsCmd)
}

// Cmd is the root command for the note package.
var Cmd = &cobra.Command{
	Use:   "note",
	Short: "Parent command for any related operation with the note package.",
}
