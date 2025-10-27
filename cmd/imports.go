package cmd

import (
	"github.com/cloud-ru/evo-ai-agents-cli/cmd/agent"
	"github.com/cloud-ru/evo-ai-agents-cli/cmd/ci"
	"github.com/cloud-ru/evo-ai-agents-cli/cmd/common"
	"github.com/cloud-ru/evo-ai-agents-cli/cmd/mcp_server"
	"github.com/cloud-ru/evo-ai-agents-cli/cmd/prompt"
	registryCmd "github.com/cloud-ru/evo-ai-agents-cli/cmd/registry"
	"github.com/cloud-ru/evo-ai-agents-cli/cmd/system"
	"github.com/cloud-ru/evo-ai-agents-cli/cmd/trigger"
)

func init() {
	RootCMD.AddCommand(
		agent.RootCMD,
		ci.RootCMD,
		common.RootCMD,
		mcp_server.RootCMD,
		prompt.RootCMD,
		registryCmd.RootCMD,
		system.RootCMD,
		trigger.RootCMD,
	)
}
