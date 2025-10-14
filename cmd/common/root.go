package common

import (
	"github.com/spf13/cobra"
)

var (
	isVerbose bool
)

// RootCMD represents the base command when called without any subcommands
var RootCMD = &cobra.Command{
	Use:   "common",
	Short: "Общие функции",
	Long:  "Общие функции и утилиты CLI",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	Args: cobra.ArbitraryArgs,
}

func init() {

}
