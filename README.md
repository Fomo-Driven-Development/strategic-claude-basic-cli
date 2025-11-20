# Strategic Claude Basic CLI

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Version](https://img.shields.io/badge/version-0.1.0-blue.svg)](https://github.com/Fomo-Driven-Development/strategic-claude-basic-cli/releases)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A command-line tool that simplifies the integration of the [Strategic Claude Basic framework](https://github.com/Fomo-Driven-Development/strategic-claude-base) into your development projects.

> **For framework usage, slash commands, and workflow documentation**, see the [Strategic Claude Basic Framework](https://github.com/Fomo-Driven-Development/strategic-claude-base) repository.

## Overview

Strategic Claude Basic CLI automates the complex setup process of integrating Strategic Claude Basic framework with Claude Code. Instead of manually cloning repositories, copying directories, and creating symlinks, this tool handles everything with a single command.

### Key Features

- **One-command installation** - Set up Strategic Claude Basic in any project instantly
- **Smart updates** - Update framework core while preserving your custom content
- **Automatic symlinks** - Creates proper `.claude` integration for Claude Code
- **Status validation** - Verify installation health and detect issues
- **Clean removal** - Safely remove framework without affecting your work
- **Shell completions** - Tab completion for bash, zsh, fish, and PowerShell
- **Zero dependencies** - Uses only Git and Go standard library

### What it does

1. **Clones** Strategic Claude Basic template at a fixed, tested commit
2. **Installs** `.strategic-claude-basic/` directory with framework content
3. **Creates** `.claude/` directory structure for Claude Code integration
4. **Links** Claude Code directories to framework components via symlinks
5. **Preserves** your custom content during updates

## Installation

### Prerequisites

- **Git** - Must be installed and available in your PATH
- **Go 1.21+** - Required for building from source

### Install with Go

```bash
go install github.com/Fomo-Driven-Development/strategic-claude-basic-cli/cmd/strategic-claude@latest
```

### Build from source

```bash
git clone https://github.com/Fomo-Driven-Development/strategic-claude-basic-cli.git
cd strategic-claude-basic-cli
make build
# Binary will be available at ./bin/strategic-claude
```

### Install binary to PATH

```bash
make install  # Installs to $GOPATH/bin or $HOME/go/bin
```

## Quick Start

1. **Initialize** Strategic Claude Basic in your project:
   ```bash
   strategic-claude init
   ```

2. **Check status** of your installation:
   ```bash
   strategic-claude status
   ```

3. **Update** framework core files (preserves your custom content):
   ```bash
   strategic-claude init --force-core
   ```

## Usage

### Initialize Framework (`init`)

Install Strategic Claude Basic in a new or existing project:

```bash
# Install in current directory
strategic-claude init

# Install in specific directory
strategic-claude init ./my-project

# Preview what would be installed (dry run)
strategic-claude init --dry-run

# Install with auto-confirmation
strategic-claude init --yes
```

**Update existing installations:**

```bash
# Update only framework files, preserve user content
strategic-claude init --force-core

# Completely overwrite existing installation
strategic-claude init --force
```

### Check Status (`status`)

Verify your installation and diagnose issues:

```bash
# Check current directory
strategic-claude status

# Check specific directory
strategic-claude status --target ./my-project

# Verbose output with detailed diagnostics
strategic-claude status --verbose
```

### Clean Installation (`clean`)

Remove Strategic Claude Basic from your project:

```bash
# Remove with confirmation prompt
strategic-claude clean

# Force removal without confirmation
strategic-claude clean --force

# Clean specific directory
strategic-claude clean ./my-project
```

### Shell Completions (`completions`)

Set up tab completion for your shell:

```bash
# Generate completions for your shell
strategic-claude completions bash
strategic-claude completions zsh
strategic-claude completions fish
strategic-claude completions powershell

# Install bash completions (example)
strategic-claude completions bash > /usr/local/etc/bash_completion.d/strategic-claude
```

## Directory Structure

After installation, your project will have this structure:

```
your-project/
├── .claude/                         # Claude Code integration
│   ├── agents/                      # Custom agent definitions
│   │   └── strategic -> ../../.strategic-claude-basic/core/agents
│   ├── commands/                    # Custom commands
│   │   └── strategic -> ../../.strategic-claude-basic/core/commands
│   ├── hooks/                       # Git hooks
│   │   └── strategic -> ../../.strategic-claude-basic/core/hooks
│   └── settings.json
└── .strategic-claude-basic/         # Framework installation
    ├── archives/                    # Archived documentation (preserved)
    │   ├── decisions/               # Archived decisions
    │   ├── issues/                  # Archived issues
    │   ├── plan/                    # Archived plans
    │   ├── product/                 # Archived product docs
    │   ├── research/                # Archived research
    │   └── summary/                 # Archived summaries
    ├── core/                        # Commands and agent definitions (updated)
    │   ├── agents/                  # Core agent definitions
    │   ├── commands/                # Core commands
    │   └── hooks/                   # Core hooks
    ├── decisions/                   # Decision records (preserved)
    │   └── CLAUDE.md
    ├── issues/                      # Issue tracking (preserved)
    │   └── CLAUDE.md
    ├── plan/                        # Implementation plans (preserved)
    │   └── CLAUDE.md
    ├── product/                     # Product documentation (preserved)
    │   └── CLAUDE.md
    ├── research/                    # Research documentation (preserved)
    │   └── CLAUDE.md
    ├── summary/                     # Work summaries (preserved)
    │   └── CLAUDE.md
    ├── templates/                   # Document templates (updated)
    │   ├── agents/                  # Agent templates
    │   ├── commands/                # Command templates
    │   ├── hooks/                   # Hook templates
    │   ├── ignore/                  # Ignore file templates
    │   └── mcps/                    # MCP templates
    └── tools/                       # Utility tools (preserved)
        └── .gitkeep
```

### Framework vs User Content

**Framework directories** (replaced during `--force-core` updates):
- `core/` - Strategic Claude agents, commands, and hooks
- `templates/` - Framework templates and examples

**User directories** (preserved during `--force-core` updates):
- `archives/` - Your archived work and completed projects
- `decisions/` - Decision records and architectural documentation
- `issues/` - Project-specific issue tracking
- `plan/` - Your planning documents and strategies
- `product/` - Product documentation and roadmaps
- `research/` - Your research notes and findings
- `summary/` - Your summaries and reports
- `tools/` - Custom utility tools and scripts

## Framework Usage

Once installed, the Strategic Claude Basic framework provides structured workflows for AI-assisted development:

### Basic Workflow
```
/research → /plan → /read_execute_plan → /summarize
```

### Context Management
- **Always run `/context` between commands** to monitor context usage
- Keep context under 40% for optimal performance  
- Use `/compact` or `/clear` when approaching limits
- Run `/summarize` before clearing context to capture incomplete work

### Key Commands
- `/research` - Analyze codebase and requirements with parallel sub-agents
- `/plan` - Create detailed implementation plans with phases and checkboxes
- `/read_execute_plan` - Execute plans systematically, tracking progress
- `/summarize` - Document problems and incomplete work for future sessions

For complete documentation on workflows, slash commands, hooks, and advanced usage, see the [Strategic Claude Basic Framework](https://github.com/Fomo-Driven-Development/strategic-claude-base) repository.

## Installation Types

### New Installation
For projects without Strategic Claude Basic:
```bash
strategic-claude init
```
- Installs complete `.strategic-claude-basic/` directory
- Creates `.claude/` symlinks
- Safe - fails if installation already exists

### Core Update (`--force-core`)
For updating framework while preserving your work:
```bash
strategic-claude init --force-core
```
- Updates `core/`, `templates/` directories
- Preserves `archives/`, `decisions/`, `issues/`, `plan/`, `product/`, `research/`, `summary/`, `tools/`
- Maintains your custom content and configurations

### Full Overwrite (`--force`)
For complete reinstallation:
```bash
strategic-claude init --force
```
- Replaces entire `.strategic-claude-basic/` directory
- **Warning**: This will overwrite all your custom user content
- Creates backup unless `--no-backup` is specified

## Commands Reference

| Command | Purpose | Key Flags |
|---------|---------|-----------|
| `init` | Install/update Strategic Claude Basic | `--force-core`, `--force`, `--yes`, `--dry-run` |
| `status` | Check installation health | `--verbose` |
| `clean` | Remove Strategic Claude Basic | `--force` |
| `completions` | Generate shell completions | Shell type argument |
| `version` | Show version information | - |

For detailed help on any command:
```bash
strategic-claude [command] --help
```

## Development

### Building

```bash
# Build binary
make build

# Build and run
make run

# Run tests
make test

# Run tests with coverage
make test-coverage

# Lint code
make lint

# Clean build artifacts
make clean
```

### Testing

```bash
# Run all tests
make test

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
make test-coverage
```

## Version Management

Strategic Claude Basic CLI pins the framework to specific, tested commits to ensure stability:
- All installations use the same framework version
- Framework updates require new CLI releases
- No unexpected changes from upstream repository updates
- Consistent, predictable behavior across installations

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: Check the docs in this repository
- **Issues**: Report bugs or request features on [GitHub Issues](https://github.com/Fomo-Driven-Development/strategic-claude-basic-cli/issues)
- **Discussions**: Join conversations on [GitHub Discussions](https://github.com/Fomo-Driven-Development/strategic-claude-basic-cli/discussions)
