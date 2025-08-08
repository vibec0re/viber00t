# viber00t

Containerized development environments with Podman 💜

## Quick Start

```bash
# Build the base image
./viber00t build

# Initialize config in your project
cd my-project
../viber00t init

# Run container
../viber00t run
```

## Configuration

Edit `viber00t.toml`:

```toml
[project]
name = "my-app"
agent = "claude-code"
privileged = true  # Enable Docker-in-Podman

[[install]]
packages = ["docker-compose", "postgresql-client"]

[[volumes]]
source = "~/data"
target = "/c0de/data"

[[ports]]
host = 3000
container = 3000
```

## Features

- Auto-mounts project to `/c0de/project`
- Claude credentials from `~/.claude/.credentials.json`
- Docker-in-Podman support with privileged mode
- Dynamic package installation
- Port forwarding and volume mounts
- Multiple AI agents support

## Usage

```bash
viber00t         # Run container (default)
viber00t init    # Create viber00t.toml
viber00t build   # Build base image
viber00t run     # Run container
```