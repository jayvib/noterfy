package utilscmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"noteapp/note"
	"noteapp/note/proto/protoutil"
	"os"
)

var (
	fileName string
)

func init() {
	ReadProtoFromFile.Flags().StringVarP(&fileName, "filename", "f", "note.pb", "The filepath to the file.")
}

// ReadProtoFromFile is a cli cmd that reads a protobuf binary file
// and print the note content in the terminal.
var ReadProtoFromFile = &cobra.Command{
	Use:   "read-proto-file",
	Short: "Use to read the protocol buffers from file",
	Long: `Use to read the protocol buffers from file.

This will read all the notes that are stored in the file then
print to the terminal.
`,
	Example: "noteapp_cli note utils read-proto-file --filename ./note.pb",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(fileName)
		if err != nil {
			logrus.Fatal(err)
		}
		defer func() { _ = file.Close() }()

		notes, err := protoutil.ReadAllProtoMessages(file)
		if err != nil {
			logrus.Fatal(err)
		}

		note.Notes(notes).ForEach(func(n *note.Note) bool {
			fmt.Println(n)
			return false
		})
	},
}
