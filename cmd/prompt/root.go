package prompt

import (
	"github.com/cloud-ru/evo-ai-agents-cli/localizations"

	"github.com/spf13/cobra"
)

var (
	isVerbose bool
)

// RootCMD represents the base command when called without any subcommands
var RootCMD = &cobra.Command{
	Use:   "prompt",
	Short: localizations.Localization.Get("root_short"),
	Long:  localizations.Localization.Get("root_long"),
	Run: func(cmd *cobra.Command, args []string) {

	},
	Args: cobra.ArbitraryArgs,
}

func init() {

}
