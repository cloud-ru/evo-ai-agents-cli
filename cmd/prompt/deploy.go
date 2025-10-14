package prompt

import (
	"fmt"

	"github.com/cloud-ru/evo-ai-agents-cli/localizations"

	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: localizations.Localization.Get("deploy_short"),
	Long:  localizations.Localization.Get("deploy_long"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deploy called")
	},
}

func init() {
	RootCMD.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
