package triigers

import (
	"os"

	"github.com/cloudru/ai-agents-cli/localizations"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	isVerbose bool
)

// RootCMD represents the base command when called without any subcommands
var RootCMD = &cobra.Command{
	Use:   "ai-agents-cli",
	Short: localizations.Localization.Get("root_short"),
	Long:  localizations.Localization.Get("root_long"),
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.New(os.Stderr)
		logger.SetReportTimestamp(true)
		logger.SetReportCaller(true)
		//logger.SetLevel(log.FatalLevel)
		logger.Info("AAAA")

		if isVerbose {
			logger.Info(logger.GetLevel())
			logger.SetLevel(log.DebugLevel)
			logger.Info("IS Verbose")
		}

		log.SetDefault(logger)
		log.Debug("AAAA")
	},
	Args: cobra.ArbitraryArgs,
}

func init() {
	RootCMD.PersistentFlags().
		BoolVarP(&isVerbose, "verbose", "v", false, "Детализация процесса")
}
