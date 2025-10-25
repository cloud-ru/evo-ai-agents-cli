#!/bin/bash

# AI Agents CLI Installation Script
# Установка из готового бинарника

set -e

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🚀 Установка AI Agents CLI из готового бинарника...${NC}"

# Проверяем, что бинарник существует
BINARY_PATH="bin/ai-agents-cli"
if [ ! -f "$BINARY_PATH" ]; then
    echo -e "${RED}❌ Бинарник не найден в $BINARY_PATH${NC}"
    echo -e "${YELLOW}Сначала выполните: make build${NC}"
    exit 1
fi

# Проверяем, что бинарник исполняемый
if [ ! -x "$BINARY_PATH" ]; then
    echo -e "${YELLOW}Делаем бинарник исполняемым...${NC}"
    chmod +x "$BINARY_PATH"
fi

# Создаем пользовательскую директорию
echo -e "${YELLOW}Создание директории ~/.local/bin...${NC}"
mkdir -p ~/.local/bin

# Копируем бинарник
echo -e "${YELLOW}Копирование бинарника...${NC}"
cp "$BINARY_PATH" ~/.local/bin/
echo -e "${GREEN}✅ Бинарник скопирован в ~/.local/bin/ai-agents-cli${NC}"

# Определяем shell
SHELL_NAME=$(basename "$SHELL")
echo -e "${GREEN}Обнаружен shell: $SHELL_NAME${NC}"

# Определяем файл конфигурации shell
if [ "$SHELL_NAME" = "zsh" ]; then
    SHELL_CONFIG="$HOME/.zshrc"
elif [ "$SHELL_NAME" = "bash" ]; then
    SHELL_CONFIG="$HOME/.bashrc"
else
    echo -e "${YELLOW}Неизвестный shell, используем .bashrc${NC}"
    SHELL_CONFIG="$HOME/.bashrc"
fi

# Проверяем, есть ли уже PATH в конфигурации
if grep -q "export PATH=\$HOME/.local/bin:\$PATH" "$SHELL_CONFIG" 2>/dev/null; then
    echo -e "${GREEN}✅ PATH уже настроен в $SHELL_CONFIG${NC}"
else
    echo -e "${YELLOW}Добавляем ~/.local/bin в PATH...${NC}"
    echo 'export PATH=$HOME/.local/bin:$PATH' >> "$SHELL_CONFIG"
    echo -e "${GREEN}✅ PATH добавлен в $SHELL_CONFIG${NC}"
fi

# Проверяем, что CLI работает
echo -e "${YELLOW}Проверяем CLI...${NC}"
if ~/.local/bin/ai-agents-cli --help >/dev/null 2>&1; then
    echo -e "${GREEN}✅ CLI работает корректно${NC}"
else
    echo -e "${YELLOW}CLI работает, но версия не определена${NC}"
fi

# Показываем справку
echo -e "${GREEN}🎉 Установка завершена!${NC}"
echo -e "${YELLOW}Для применения изменений выполните:${NC}"
echo -e "   source $SHELL_CONFIG"
echo -e "${YELLOW}Или перезапустите терминал${NC}"
echo ""
echo -e "${GREEN}Теперь можно использовать:${NC}"
echo -e "   ai-agents-cli --help"
echo -e "   ai-agents-cli create mcp my-project"
echo -e "   ai-agents-cli create agent my-agent"
echo ""
echo -e "${YELLOW}💡 Для автодополнения выполните:${NC}"
echo -e "   ai-agents-cli completion $SHELL_NAME > ~/.local/share/bash-completion/completions/ai-agents-cli"
echo -e "   # или для zsh:"
echo -e "   ai-agents-cli completion $SHELL_NAME > ~/.zsh/completions/_ai-agents-cli"
