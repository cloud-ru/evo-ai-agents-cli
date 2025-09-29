package cmd

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/cloudru/ai-agents-cli/localizations"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: localizations.Localization.Get("deploy_short"),
	Long:  localizations.Localization.Get("deploy_long"),
	Run: func(cmd *cobra.Command, args []string) {
		currentDir, err := filepath.Abs(".")
		if err != nil {
			log.Error("failed to determine current working directory", err)
			os.Exit(1)
		}
		fmt.Println(currentDir)
	},
}

func init() {
	RootCMD.AddCommand(deployCmd)

	deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
