# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build
```bash
go build -o viber00t
```

### Run
```bash
./viber00t         # Run container (default)
./viber00t init    # Create Viber00t.toml  
./viber00t clean   # Clean cached images
```

### Dependencies
```bash
go mod tidy        # Update module dependencies
```

## Architecture

viber00t is a containerized development environment tool built in Go that uses Podman to manage development containers. 

### Core Components

1. **Configuration System** (`main.go:17-103`)
   - Project config via `Viber00t.toml` 
   - Global config at `~/.config/viber00t/config.toml`
   - Environment templates for language-specific packages
   - Config hashing for image caching

2. **Container Management** (`main.go:454-590`)
   - Dynamic Dockerfile generation based on config
   - Project-specific image building with caching
   - Automatic volume mounting (project dir, credentials, SSH keys)
   - Port forwarding and privileged mode support

3. **Image Building** (`main.go:379-452`)
   - Config-based image generation
   - State tracking in XDG directories
   - Automatic rebuild on config changes
   - Cached image reuse when possible

### Key Features

- Auto-mounts project to `/c0de/project` in container
- Claude Code integration with credentials mounting
- Docker-in-Podman support via privileged mode
- Language environment templates (Python, Rust, Node, Go, etc.)
- XDG Base Directory compliance for config/cache/state

### Configuration Flow

1. Global config provides defaults and base packages
2. Project config (`Viber00t.toml`) specifies project-specific needs
3. Config hash determines if rebuild is needed
4. Dockerfile generated dynamically from merged configs
5. Images cached and reused based on config hash