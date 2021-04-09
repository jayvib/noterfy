package cli

import (
	"github.com/spf13/cobra"
	"noteapp/note/cli/utilscmd"
)

func init() {
	UtilsCmd.AddCommand(utilscmd.ReadProtoFromFile)
}

// UtilsCmd is a cli command where it contains
// any utility tool that is related to note.
var UtilsCmd = &cobra.Command{
	Use:   "utils",
	Short: "A subcommand for any utility operation for note package",
}
