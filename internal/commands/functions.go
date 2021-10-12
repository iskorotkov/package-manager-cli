package commands

import (
	"log"
	"os"

	"github.com/iskorotkov/package-manager-cli/pkg/xlog"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type commandFunc func(cmd *cobra.Command, args []string) error

func wrapCommand(cmd *cobra.Command) *cobra.Command {
	cmd.RunE = wrapCommandFunc(cmd.Use, cmd.RunE)

	return cmd
}

func wrapCommandFunc(name string, f commandFunc) commandFunc {
	return func(cmd *cobra.Command, args []string) error {
		xlog.Push(name)
		defer xlog.Pop()

		log.Printf("invoked command with args: %+v", args)
		defer log.Printf("command completed")

		err := f(cmd, args)
		if err != nil {
			log.Printf("error when executing command: %v", err)
		}

		return err
	}
}

func createTable() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false

	return t
}
