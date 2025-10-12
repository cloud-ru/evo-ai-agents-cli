# Environment Configuration

## Quick Start

1. **Copy the example file:**
   ```bash
   cp env.example .env
   ```

2. **Edit `.env` with your values:**
   ```bash
   nano .env
   ```

3. **Run the CLI:**
   ```bash
   ./bin/ai-agents-cli validate examples/
   ```

## Required Variables

| Variable | Description |
|----------|-------------|
| `IAM_KEY_ID` | Your IAM Key ID |
| `IAM_SECRET` | Your IAM Secret |
| `PROJECT_ID` | Your AI Agents Project ID |

## Example Configuration

```bash
# IAM Authentication
IAM_KEY_ID=your-iam-key-id
IAM_SECRET=your-iam-secret
IAM_ENDPOINT=https://iam.api.cloud.ru

# Project Configuration
PROJECT_ID=your-project-id
PUBLIC_API_ENDPOINT=ai-agents.api.cloud.ru

# Service Configuration
SERVICE_APP_ENVIRONMENT=dev
SERVICE_NAME=ai-agents-cli
SERVICE_LOG_LEVEL=debug
SERVICE_VERSION=1.0.0
```

## Files

- `env.example` - Template with all available variables
- `env.local` - Local test configuration (not committed)
- `.env` - Your local configuration (not committed)
- `ENV_SETUP.md` - Detailed setup instructions

## Security

- Never commit `.env` file to version control
- Use strong, unique secrets for production
- Rotate IAM credentials regularly
