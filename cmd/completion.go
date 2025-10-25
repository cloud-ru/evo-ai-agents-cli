package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Генерация скриптов автодополнения для оболочки",
	Long: `Генерация скриптов автодополнения для различных оболочек.

Поддерживаемые оболочки:
  - bash: для Bash и совместимых оболочек
  - zsh:  для Zsh
  - fish: для Fish shell
  - powershell: для PowerShell

Примеры использования:
  # Генерация для bash
  ai-agents-cli completion bash > ~/.bash_completion.d/ai-agents-cli
  
  # Генерация для zsh
  ai-agents-cli completion zsh > ~/.zsh/completions/_ai-agents-cli
  
  # Генерация для fish
  ai-agents-cli completion fish > ~/.config/fish/completions/ai-agents-cli.fish
  
  # Генерация для PowerShell
  ai-agents-cli completion powershell > ~/.config/powershell/ai-agents-cli.ps1
  
  # Автоматическая установка для текущей оболочки
  source <(ai-agents-cli completion bash)
`,
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Args:      cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shell := args[0]

		var err error
		switch shell {
		case "bash":
			err = RootCMD.GenBashCompletion(os.Stdout)
		case "zsh":
			err = RootCMD.GenZshCompletion(os.Stdout)
		case "fish":
			err = RootCMD.GenFishCompletion(os.Stdout, true)
		case "powershell":
			err = RootCMD.GenPowerShellCompletion(os.Stdout)
		default:
			fmt.Fprintf(os.Stderr, "Неподдерживаемая оболочка: %s\n", shell)
			os.Exit(1)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка генерации автодополнения: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCMD.AddCommand(completionCmd)
}
