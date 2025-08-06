package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type WrongArgsError struct {
	Arg int
	Cmd *cobra.Command
}

func (e *WrongArgsError) Error() string {
	return fmt.Sprintf("Error, exptected %d arguments", e.Arg)
}
