# Environment Configuration

## Quick Setup

1. **Copy the example file:**
   ```bash
   cp env.example .env
   ```

2. **Edit `.env` with your values:**
   ```bash
   nano .env
   ```

3. **Required variables:**
   - `IAM_KEY_ID` - Your IAM Key ID
   - `IAM_SECRET` - Your IAM Secret
   - `PROJECT_ID` - Your AI Agents Project ID

## Environment Variables

### IAM Authentication
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `IAM_KEY_ID` | IAM Key ID for authentication | - | ✅ |
| `IAM_SECRET` | IAM Secret for authentication | - | ✅ |
| `IAM_ENDPOINT` | IAM API endpoint | `https://iam.api.cloud.ru` | ❌ |

### Project Configuration
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PROJECT_ID` | AI Agents Project ID | - | ✅ |
| `PUBLIC_API_ENDPOINT` | AI Agents API endpoint | `ai-agents.api.cloud.ru` | ❌ |

### Service Configuration
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SERVICE_APP_ENVIRONMENT` | Environment (dev/stage/prod) | `dev` | ❌ |
| `SERVICE_NAME` | Service name | `ai-agents-cli` | ❌ |
| `SERVICE_LOG_LEVEL` | Log level (debug/info/warn/error) | `debug` | ❌ |
| `SERVICE_VERSION` | Service version | `1.0.0` | ❌ |

### IAM Service Account
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `IAM_BFF_HOST` | IAM BFF host | `<HOST>` | ❌ |
| `IAM_CLIENT_ID` | IAM Client ID | `<CLIENT_ID>` | ❌ |
| `IAM_CLIENT_SECRET` | IAM Client Secret | `<CLIENT_SECRET>` | ❌ |

### Performance
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `BULK_OPERATIONS_CONCURRENCY` | Bulk operations concurrency factor | `20` | ❌ |

## Usage

The CLI will automatically load environment variables from `.env` file if present.

### Manual Environment Setup

If you prefer to set environment variables manually:

```bash
export IAM_KEY_ID="your-key-id"
export IAM_SECRET="your-secret"
export PROJECT_ID="your-project-id"
./bin/ai-agents-cli validate examples/
```

### Docker Environment

For Docker usage, you can pass environment variables:

```bash
docker run -e IAM_KEY_ID="your-key-id" \
           -e IAM_SECRET="your-secret" \
           -e PROJECT_ID="your-project-id" \
           ai-agents-cli validate examples/
```

## Security Notes

- Never commit `.env` file to version control
- Use strong, unique secrets for production
- Rotate IAM credentials regularly
- Consider using secret management systems for production
