# 🚀 Автодополнение для ai-agents-cli

## 📋 Поддерживаемые оболочки

- **Bash** - для Bash и совместимых оболочек
- **Zsh** - для Zsh
- **Fish** - для Fish shell
- **PowerShell** - для PowerShell

## 🚀 Автоматическая установка

### Быстрая установка
```bash
# Автоматическая установка для текущей оболочки
./scripts/install-completion.sh
```

### Ручная установка

#### Bash
```bash
# Генерация completion файла
ai-agents-cli completion bash > ~/.bash_completion.d/ai-agents-cli

# Добавление в .bashrc (если еще не добавлено)
echo 'source ~/.bash_completion.d/ai-agents-cli' >> ~/.bashrc

# Применение изменений
source ~/.bashrc
```

#### Zsh
```bash
# Генерация completion файла
ai-agents-cli completion zsh > ~/.zsh/completions/_ai-agents-cli

# Добавление в .zshrc (если еще не добавлено)
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
echo 'autoload -U compinit && compinit' >> ~/.zshrc

# Применение изменений
source ~/.zshrc
```

#### Fish
```bash
# Автоматическая установка для Fish
fish scripts/install-fish-completion.fish

# Или ручная установка
ai-agents-cli completion fish > ~/.config/fish/completions/ai-agents-cli.fish

# Применение изменений
source ~/.config/fish/config.fish
```

#### PowerShell
```powershell
# Генерация completion файла
ai-agents-cli completion powershell > ~/.config/powershell/ai-agents-cli.ps1

# Добавление в профиль PowerShell
echo 'ai-agents-cli completion powershell | Out-String | Invoke-Expression' >> $PROFILE
```

## 🎯 Использование

После установки автодополнения вы можете использовать **Tab** для автодополнения команд:

### Fish Shell особенности
Fish shell имеет встроенную поддержку автодополнения и не требует дополнительной настройки после установки файла completion. Fish автоматически:
- Загружает completion файлы из `~/.config/fish/completions/`
- Предоставляет красивый интерфейс автодополнения
- Показывает описания команд при наведении
- Поддерживает частичное совпадение команд

```bash
# Автодополнение основных команд
ai-agents-cli <TAB>
# Покажет: agent, ci, create, deploy, mcp-server, prompt, system, trigger

# Автодополнение подкоманд
ai-agents-cli create <TAB>
# Покажет: agent, mcp

# Автодополнение флагов
ai-agents-cli create mcp --<TAB>
# Покажет: --author, --help, --path, --verbose

# Автодополнение значений
ai-agents-cli create mcp my-project --path <TAB>
# Покажет доступные директории
```

## 🔧 Команды автодополнения

### Генерация completion
```bash
# Для bash
ai-agents-cli completion bash

# Для zsh
ai-agents-cli completion zsh

# Для fish
ai-agents-cli completion fish

# Для PowerShell
ai-agents-cli completion powershell
```

### Проверка установки
```bash
# Проверка наличия команды completion
ai-agents-cli completion --help

# Проверка автодополнения
ai-agents-cli <TAB>
```

## 🐛 Устранение неполадок

### Проблема: Автодополнение не работает
```bash
# Проверьте, что completion файл существует
ls ~/.bash_completion.d/ai-agents-cli
ls ~/.zsh/completions/_ai-agents-cli
ls ~/.config/fish/completions/ai-agents-cli.fish

# Проверьте, что файл загружается
source ~/.bashrc
source ~/.zshrc
source ~/.config/fish/config.fish
```

### Проблема: Completion файл не генерируется
```bash
# Проверьте, что ai-agents-cli установлен
which ai-agents-cli

# Проверьте версию
ai-agents-cli --version

# Проверьте права доступа
chmod +x $(which ai-agents-cli)
```

### Проблема: Неправильная оболочка
```bash
# Проверьте текущую оболочку
echo $SHELL

# Проверьте, что completion установлен для правильной оболочки
ls ~/.bash_completion.d/ai-agents-cli
ls ~/.zsh/completions/_ai-agents-cli
ls ~/.config/fish/completions/ai-agents-cli.fish
```

### Проблема: Fish completion не работает
```fish
# Проверьте, что файл существует
ls ~/.config/fish/completions/ai-agents-cli.fish

# Проверьте, что Fish загружает completion
fish -c "complete -c ai-agents-cli"

# Перезапустите Fish shell
exec fish
```

## 📁 Структура файлов

```
~/.bash_completion.d/ai-agents-cli          # Bash completion
~/.zsh/completions/_ai-agents-cli           # Zsh completion
~/.config/fish/completions/ai-agents-cli.fish # Fish completion
~/.config/powershell/ai-agents-cli.ps1     # PowerShell completion
```

## 🔄 Обновление автодополнения

При обновлении ai-agents-cli рекомендуется обновить и автодополнение:

```bash
# Автоматическое обновление
./scripts/install-completion.sh

# Или ручное обновление
ai-agents-cli completion bash > ~/.bash_completion.d/ai-agents-cli
ai-agents-cli completion zsh > ~/.zsh/completions/_ai-agents-cli
ai-agents-cli completion fish > ~/.config/fish/completions/ai-agents-cli.fish
```

## 🎨 Кастомизация

### Добавление собственных completion
Вы можете расширить автодополнение, добавив собственные функции в completion файлы.

### Изменение поведения
Вы можете изменить поведение автодополнения, отредактировав сгенерированные файлы.

## 📚 Дополнительные ресурсы

- [Cobra Completion Documentation](https://github.com/spf13/cobra/blob/master/shell_completions.md)
- [Bash Completion](https://www.gnu.org/software/bash/manual/html_node/Programmable-Completion.html)
- [Zsh Completion](http://zsh.sourceforge.net/Doc/Release/Completion-System.html)
- [Fish Completion](https://fishshell.com/docs/current/completions.html)
- [PowerShell Completion](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_profiles)

## 🆘 Поддержка

Если у вас возникли проблемы с автодополнением:

1. Проверьте, что ai-agents-cli установлен и работает
2. Проверьте, что completion файл существует и загружается
3. Проверьте, что вы используете правильную оболочку
4. Попробуйте переустановить автодополнение

Для получения помощи создайте issue в репозитории проекта.
