package cmd

import (
	"github.com/cloudru/ai-agents-cli/cmd/agent"
	"github.com/cloudru/ai-agents-cli/cmd/common"
	"github.com/cloudru/ai-agents-cli/cmd/mcp_server"
	"github.com/cloudru/ai-agents-cli/cmd/prompt"
	"github.com/cloudru/ai-agents-cli/cmd/system"
	"github.com/cloudru/ai-agents-cli/cmd/trigger"
)

func init() {
	RootCMD.AddCommand(
		agent.RootCMD,
		common.RootCMD,
		mcp_server.RootCMD,
		prompt.RootCMD,
		system.RootCMD,
		trigger.RootCMD,
	)
}
