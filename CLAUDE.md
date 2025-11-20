# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Strategic Claude Basic CLI is a Go-based command-line tool that automates the setup and management of Strategic Claude Basic framework integration. It clones the strategic-claude-base repository, copies contents to `.strategic-claude-basic/`, and creates symlinks in `.claude/` subdirectories for Claude Code integration.

## Common Commands

### Development Commands
```bash
make build          # Build the CLI binary with version injection
make run            # Build and run the CLI
make test           # Run all tests
make test-coverage  # Run tests with coverage report (generates coverage.out and HTML)
make lint           # Run golangci-lint for code quality
make clean          # Clean build artifacts and coverage files
make deps           # Download and tidy Go modules
make install        # Install binary to $GOPATH/bin or $HOME/go/bin
make help           # Show all available targets
```

### CLI Usage
```bash
# Build and test the CLI commands
./bin/strategic-claude --help
./bin/strategic-claude version
./bin/strategic-claude init --help
./bin/strategic-claude status --help
```

### Version Management
The build system automatically injects version information using ldflags:
- `VERSION`: Defaults to "0.1.0", can be overridden
- `COMMIT`: Auto-detected from git HEAD
- `DATE`: Auto-generated build timestamp

## Architecture

### High-Level Structure
The codebase follows a layered architecture:
- **CLI Layer**: Cobra-based commands in `cmd/strategic-claude-basic-cli/`
- **Models Layer**: Data structures, configuration, and error types in `internal/models/`
- **Services Layer**: Business logic (planned for Phase 2) in `internal/services/`
- **Utils Layer**: Validation and user interaction utilities in `internal/utils/`
- **Config Layer**: Constants and configuration in `internal/config/`

### Key Components

#### Commands (cmd/strategic-claude-basic-cli/)
- `main.go`: Entry point, calls Execute()
- `root.go`: Root command with global flags (--verbose, --target)
- `init.go`: Installation command (stubbed, Phase 4 implementation)
- `status.go`: Status checking command (stubbed, Phase 3 implementation)
- `clean.go`: Cleanup command (stubbed, Phase 5 implementation)
- `completions.go`: Shell completion generation
- `version.go`: Version information display

#### Models (internal/models/)
- `config.go`: InstallConfig and CleanConfig with validation methods
- `status.go`: StatusInfo, SymlinkStatus, and InstallationPlan structures
- `errors.go`: Structured error handling with ErrorCode enum and user-friendly messages

#### Configuration (internal/config/constants.go)
Critical constants for framework behavior:
- **Repository**: `Fomo-Driven-Development/strategic-claude-base.git` at fixed commit
- **Directories Replaced**: core/, templates/ (framework content)
- **Directories Preserved**: archives/, decisions/, issues/, plan/, product/, research/, summary/, tools/ (user content)
- **Symlinks Created**:
  - `.claude/agents/strategic` → `../../.strategic-claude-basic/core/agents`
  - `.claude/commands/strategic` → `../../.strategic-claude-basic/core/commands`
  - `.claude/hooks/strategic` → `../../.strategic-claude-basic/core/hooks`

#### Utilities (internal/utils/)
- `interaction.go`: User prompts, confirmation dialogs, formatted output
- `validation.go`: Path validation, directory checks, input validation, symlink validation

### Framework Integration Behavior

#### Installation Types
1. **New Installation**: Full `.strategic-claude-basic/` directory copy
2. **Core Update** (`--force-core`): Replace only core/, templates/ (preserve user content in archives/, decisions/, issues/, plan/, product/, research/, summary/, tools/)
3. **Full Overwrite** (`--force`): Replace entire directory

#### Directory Structure
```
target-project/
├── .strategic-claude-basic/          # Framework installation
│   ├── core/                        # Framework code (replaced)
│   │   ├── agents/                  # Symlink targets
│   │   ├── commands/                # Symlink targets
│   │   └── hooks/                   # Symlink targets
│   ├── templates/                   # Framework templates (replaced)
│   ├── archives/                    # User content (preserved)
│   │   ├── decisions/               # Archived decisions
│   │   ├── issues/                  # Archived issues
│   │   ├── plan/                    # Archived plans
│   │   ├── product/                 # Archived product docs
│   │   ├── research/                # Archived research
│   │   └── summary/                 # Archived summaries
│   ├── decisions/                   # User content (preserved)
│   ├── issues/                      # User content (preserved)
│   ├── plan/                        # User content (preserved)
│   ├── product/                     # User content (preserved)
│   ├── research/                    # User content (preserved)
│   ├── summary/                     # User content (preserved)
│   └── tools/                       # User content (preserved)
└── .claude/                         # Claude Code integration
    ├── agents/
    │   └── strategic -> ../../.strategic-claude-basic/core/agents
    ├── commands/
    │   └── strategic -> ../../.strategic-claude-basic/core/commands
    └── hooks/
        └── strategic -> ../../.strategic-claude-basic/core/hooks
```

## Development Status

### Completed (Phase 1 - Project Foundation)
- ✅ Basic CLI framework with Cobra
- ✅ Version injection and build system
- ✅ Core data models and error handling
- ✅ Validation utilities and user interaction
- ✅ Shell completion support
- ✅ Pre-commit hooks with golangci-lint

### Next Phases
- **Phase 2**: Core Services (Git operations, file system utilities)
- **Phase 3**: Status Command (installation detection, symlink validation)
- **Phase 4**: Init Command (installation logic, symlink creation)
- **Phase 5**: Clean Command (safe cleanup with user content preservation)

## Quality Controls

### Pre-commit Hooks
- `go-fmt`: Code formatting
- `go-imports`: Import organization
- `go-mod-tidy`: Module maintenance
- `go-build`: Compilation verification
- `golangci-lint`: Comprehensive linting (includes vet checks)
- Standard file checks (trailing whitespace, end-of-file, YAML validation)

### Linting Configuration
Uses `.golangci.yml` with additional linters enabled:
- goconst, gocritic, gocyclo, gosec, misspell, unparam, unconvert, whitespace, prealloc, predeclared
- Test files have relaxed linting rules
- File system operations exclude G304 (file path from variable) for installer paths

### Testing Strategy
- Unit tests for utilities and validation
- Integration tests for complete command flows
- Mock file system for testing edge cases
- Target >80% test coverage
