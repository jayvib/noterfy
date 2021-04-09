package main

import (
	"log"
	"noterfy/cli"
	notecli "noterfy/note/cli"
)

func main() {
	cli.RootCmd.AddCommand(notecli.Cmd)
	if err := cli.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
