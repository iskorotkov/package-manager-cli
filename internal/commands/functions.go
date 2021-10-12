package commands

import (
	"log"

	"github.com/iskorotkov/package-manager-cli/pkg/xlog"
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
