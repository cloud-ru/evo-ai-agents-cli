class AiAgentsCli < Formula
  desc "Command-line tool for managing AI agents and MCP servers"
  homepage "https://github.com/cloud-ru/evo-ai-agents-cli"
  url "https://github.com/cloud-ru/evo-ai-agents-cli/releases/download/v1.0.0/ai-agents-cli-darwin-amd64.tar.gz"
  sha256 "PLACEHOLDER_SHA256"
  license "MIT"
  version "1.0.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/cloud-ru/evo-ai-agents-cli/releases/download/v1.0.0/ai-agents-cli-darwin-arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_ARM64"
    end
  end

  on_linux do
    url "https://github.com/cloud-ru/evo-ai-agents-cli/releases/download/v1.0.0/ai-agents-cli-linux-amd64.tar.gz"
    sha256 "PLACEHOLDER_SHA256_LINUX"
  end

  def install
    bin.install "ai-agents-cli"
    
    # Install bash completion
    bash_completion.install "completions/bash/ai-agents-cli.bash" if File.exist?("completions/bash/ai-agents-cli.bash")
    
    # Install zsh completion
    zsh_completion.install "completions/zsh/_ai-agents-cli" if File.exist?("completions/zsh/_ai-agents-cli")
    
    # Install fish completion
    fish_completion.install "completions/fish/ai-agents-cli.fish" if File.exist?("completions/fish/ai-agents-cli.fish")
    
    # Install man pages
    man1.install "man/ai-agents-cli.1" if File.exist?("man/ai-agents-cli.1")
  end

  test do
    # Test basic functionality
    assert_match "AI Agents CLI", shell_output("#{bin}/ai-agents-cli --help", 1)
    
    # Test version
    assert_match "1.0.0", shell_output("#{bin}/ai-agents-cli version", 1)
  end

  def caveats
    <<~EOS
      AI Agents CLI has been installed!
      
      To get started:
        1. Set your IAM credentials: export IAM_KEY_ID="your-iam-key-id" IAM_SECRET="your-iam-secret"
        2. Set your project ID: export PROJECT_ID="your-project-id"
        3. Run: ai-agents-cli --help
      
      For more information, visit: https://github.com/cloud-ru/evo-ai-agents-cli
    EOS
  end
end
