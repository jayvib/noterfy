package main

import (
	"log"
	"noteapp/cli"
	notecli "noteapp/note/cli"
)

func main() {
	cli.RootCmd.AddCommand(notecli.Cmd)
	if err := cli.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
