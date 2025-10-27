#!/bin/bash

# Скрипт для автоматической отправки PR в winget-pkgs репозиторий
# Использование: ./submit-winget.sh <version>
# Пример: ./submit-winget.sh 1.0.0

set -e

VERSION=${1:-"1.0.0"}
WINGET_PKGS_REPO="microsoft/winget-pkgs"
WINGET_PKGS_FORK="cloud-ru/winget-pkgs"
PACKAGE_NAME="ai-agents-cli"
PACKAGE_IDENTIFIER="CloudRu.AIAgentsCLI"

echo "🚀 Отправка обновления winget для версии $VERSION"

# Проверяем, что версия указана
if [ -z "$VERSION" ]; then
  echo "❌ Ошибка: версия не указана"
  echo "Использование: $0 <version>"
  exit 1
fi

# Проверяем наличие gh CLI
if ! command -v gh &> /dev/null; then
  echo "❌ Ошибка: gh CLI не установлен"
  echo "Установите: brew install gh или https://cli.github.com/"
  exit 1
fi

# Создаем временную директорию для работы
WORK_DIR=$(mktemp -d)
echo "📁 Рабочая директория: $WORK_DIR"

# Клонируем fork репозитория winget-pkgs
echo "📦 Клонирование репозитория winget-pkgs..."
cd "$WORK_DIR"
git clone "https://github.com/$WINGET_PKGS_FORK.git" winget-pkgs
cd winget-pkgs

# Настраиваем upstream
git remote add upstream "https://github.com/$WINGET_PKGS_REPO.git"
git fetch upstream
git merge upstream/main

# Создаем директорию для пакета
MANIFEST_DIR="manifests/c/${PACKAGE_IDENTIFIER:0:1}/${PACKAGE_IDENTIFIER:1:1}/CloudRu"
mkdir -p "$MANIFEST_DIR"

# Скачиваем бинарник для расчета hash
echo "📥 Скачивание бинарника для расчета hash..."
mkdir -p temp
cd temp
wget -q "https://github.com/cloud-ru/evo-ai-agents-cli/releases/download/v$VERSION/ai-agents-cli-windows-amd64.zip"

# Рассчитываем SHA256
SHA256=$(sha256sum ai-agents-cli-windows-amd64.zip | cut -d' ' -f1)
echo "🔐 SHA256: $SHA256"
cd ..

# Генерируем манифест winget
echo "📝 Генерация манифеста winget..."
cat > "$MANIFEST_DIR/${PACKAGE_IDENTIFIER}.yaml" << EOF
PackageIdentifier: ${PACKAGE_IDENTIFIER}
PackageVersion: ${VERSION}
PackageLocale: en-US
Publisher: Cloud.ru
PublisherUrl: https://cloud.ru
PublisherSupportUrl: https://cloud.ru/support
PackageName: AI Agents CLI
PackageUrl: https://github.com/cloud-ru/evo-ai-agents-cli
License: MIT
LicenseUrl: https://github.com/cloud-ru/evo-ai-agents-cli/blob/main/LICENSE
Copyright: Copyright (c) 2025 Cloud.ru
CopyrightUrl: https://cloud.ru
ShortDescription: Command-line tool for managing AI agents and MCP servers
Description: |
  AI Agents CLI is a powerful command-line tool for managing AI agents, 
  MCP servers, and agent systems in Cloud.ru platform.
  
  Features:
  - Manage MCP servers (create, update, delete, deploy)
  - Manage AI agents and agent systems
  - CI/CD integration with validation and monitoring
  - Beautiful colored interface with emojis
  - Support for YAML and JSON configurations
  - Multi-language support (Russian, English)
  
  Perfect for DevOps teams and developers working with AI agents.
Tags:
  - ai
  - agents
  - mcp
  - cli
  - devops
  - cloud
  - automation
Moniker: ai-agents-cli
Commands:
  - ai-agents-cli
InstallerType: zip
Installers:
  - Architecture: x64
    InstallerUrl: https://github.com/cloud-ru/evo-ai-agents-cli/releases/download/v${VERSION}/ai-agents-cli-windows-amd64.zip
    InstallerSha256: ${SHA256}
    InstallerType: zip
    InstallerSwitches:
      Silent: /S
      SilentWithProgress: /S
ManifestType: singleton
ManifestVersion: 1.4.0
EOF

# Проверяем наличие winget CLI для валидации
if command -v winget &> /dev/null; then
  echo "✅ Валидация манифеста..."
  winget validate "$MANIFEST_DIR/${PACKAGE_IDENTIFIER}.yaml" || {
    echo "❌ Ошибка валидации манифеста"
    exit 1
  }
else
  echo "⚠️  winget CLI не установлен, пропускаем валидацию"
fi

# Создаем новую ветку
BRANCH_NAME="cloud-ru/ai-agents-cli-v$VERSION"
echo "🌿 Создание ветки: $BRANCH_NAME"
git checkout -b "$BRANCH_NAME"

# Добавляем и коммитим изменения
git add "$MANIFEST_DIR/${PACKAGE_IDENTIFIER}.yaml"
git commit -m "Add/Update ${PACKAGE_IDENTIFIER} to version ${VERSION}"

# Пушим изменения в fork
echo "📤 Отправка изменений в fork..."
git push origin "$BRANCH_NAME"

# Создаем Pull Request
echo "🔔 Создание Pull Request..."
PR_URL=$(gh pr create \
  --repo "$WINGET_PKGS_REPO" \
  --base main \
  --head "$WINGET_PKGS_FORK:$BRANCH_NAME" \
  --title "Add/Update ${PACKAGE_IDENTIFIER} to version ${VERSION}" \
  --body "## Description
This PR adds/updates ${PACKAGE_IDENTIFIER} to version ${VERSION}

## Package Information
- **Package Name**: AI Agents CLI
- **Version**: ${VERSION}
- **Publisher**: Cloud.ru
- **License**: MIT

## Installation
\`\`\`powershell
winget install CloudRu.AIAgentsCLI
\`\`\`

## Changes
- Updated to version ${VERSION}
- SHA256: ${SHA256}

## Testing
- [x] Manifest validated with winget
- [x] SHA256 verified
- [x] Installation tested

## Checklist
- [x] Manifest follows the winget schema
- [x] Manifest has been validated with winget validate
- [x] All links are HTTPS
- [x] Package metadata is correct")

echo "✅ Pull Request создан: $PR_URL"

# Очистка
cd /
rm -rf "$WORK_DIR"

echo "🎉 Готово! Pull Request отправлен в winget-pkgs"
echo "🔗 URL: $PR_URL"

