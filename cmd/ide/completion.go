package ide

import (
	"github.com/cloud-ru/evo-ai-agents-cli/localizations"
	"github.com/cloud-ru/evo-ai-agents-cli/localizations/i18n_labels"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   i18n_labels.CompletionCommandLabelName,
	Short: localizations.Localization.Get(i18n_labels.CompletionShortDescLabelName),
	Long:  localizations.Localization.Get(i18n_labels.CompletionLongDescCommandName),
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	RootCMD.AddCommand(completionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
