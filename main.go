package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Project struct {
		Name       string
		Agent      string
		Privileged bool
	}
	Install []struct {
		Packages []string
		Envs     []string
	}
	Volumes []struct {
		Source string
		Target string
	}
	Ports []struct {
		Host      int
		Container int
	}
}

type GlobalConfig struct {
	DefaultAgent      string
	DefaultPrivileged bool
	DefaultImage      string
	ClaudeFlags       []string
	DefaultEnvs       []string
	DefaultPackages   []string
	BasePackages      []string // Core packages for all containers
}

var envTemplates = map[string][]string{
	"python": {"python3", "python3-dev", "python3-pip", "python3-venv", "pipx", "poetry", "pyenv", "python3-setuptools"},
	"rust":   {"rustc", "cargo", "rustfmt", "rust-src", "pkg-config", "libssl-dev"},
	"node":   {"nodejs", "npm", "yarn", "n"},
	"go":     {"golang", "gopls"},
	"ruby":   {"ruby-full", "ruby-dev", "bundler", "rbenv"},
	"java":   {"openjdk-17-jdk", "maven", "gradle"},
	"cpp":    {"clang", "clang-tools", "clang-format", "cmake", "ninja-build", "ccache", "gdb", "valgrind"},
	"php":    {"php", "php-cli", "php-mbstring", "php-xml", "composer"},
	"dotnet": {"dotnet-sdk-8.0", "nuget"},
}

const defaultConfig = `[project]
name = "my-project"
agent = "claude"
privileged = false

[[install]]
packages = []
envs = []  # Available: python, rust, node, go, ruby, java, cpp, php, dotnet

[[volumes]]
# source = "~/extra"
# target = "/c0de/extra"

[[ports]]
# host = 3000
# container = 3000
`

const defaultGlobalConfig = `# viber00t global configuration
# ~/.config/viber00t/config.toml

default_agent = "claude"
default_privileged = false
default_image = "viber00t/base:latest"

# Flags passed to claude
claude_flags = ["--dangerously-skip-permissions"]

# Base packages installed in every container
base_packages = [
  "git", "git-lfs", "build-essential", "make",
  "vim", "nano", "htop", "tmux", "tree", "ncdu",
  "jq", "ripgrep", "fd-find", "fzf", "bat",
  "httpie", "netcat-openbsd", "iputils-ping",
  "zip", "unzip", "tar", "xz-utils",
  "docker.io", "docker-compose",
  "postgresql-client", "redis-tools", "sqlite3"
]

# Default environments for all projects
default_envs = []

# Default packages for all projects
default_packages = []
`

func getXDGConfigHome() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return xdg
	}
	return filepath.Join(os.Getenv("HOME"), ".config")
}

func getXDGCacheHome() string {
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return xdg
	}
	return filepath.Join(os.Getenv("HOME"), ".cache")
}

func getXDGStateHome() string {
	if xdg := os.Getenv("XDG_STATE_HOME"); xdg != "" {
		return xdg
	}
	return filepath.Join(os.Getenv("HOME"), ".local", "state")
}

func main() {
	if len(os.Args) < 2 {
		runContainer()
		return
	}

	switch os.Args[1] {
	case "init":
		initConfig()
	case "clean":
		cleanImages()
	case "shell":
		runShell()
	case "help", "-h", "--help":
		showHelp()
	case "version", "-v", "--version":
		showVersion()
	default:
		runContainer()
	}
}

func showHelp() {
	banner := `
╦  ╦╦╔╗ ╔═╗╦═╗╔═╗╔═╗╔╦╗
╚╗╔╝║╠╩╗║╣ ╠╦╝║ ║║ ║ ║ 
 ╚╝ ╩╚═╝╚═╝╩╚═╚═╝╚═╝ ╩ 
      [ FULL SPECTRUM CYBER ]
`
	fmt.Print("\033[35m" + banner + "\033[0m")
	fmt.Println("\n\033[36mContainerized Development Environments\033[0m")
	fmt.Println("\033[90m═══════════════════════════════════════\033[0m")
	fmt.Println()
	fmt.Println("\033[33mUSAGE:\033[0m")
	fmt.Println("  viber00t              \033[90m# Run container (default)\033[0m")
	fmt.Println("  viber00t init         \033[90m# Create viber00t.toml\033[0m")
	fmt.Println("  viber00t shell        \033[90m# Interactive bash shell\033[0m")
	fmt.Println("  viber00t clean        \033[90m# Clean cached images\033[0m")
	fmt.Println()
	fmt.Println("\033[33mENVIRONMENTS:\033[0m")
	fmt.Println("  python, rust, node, go, ruby, java, cpp, php, dotnet")
	fmt.Println()
	fmt.Println("\033[35m» vibec0re.github.io\033[0m")
}

func showVersion() {
	fmt.Println("\033[35mviber00t v1.0.0\033[0m - Full Spectrum Cyber")
	fmt.Println("\033[90mvibec0re.github.io\033[0m")
}

func initConfig() {
	// Initialize global config first
	initGlobalConfig()

	if _, err := os.Stat("viber00t.toml"); err == nil {
		fmt.Println("\033[33m⚠\033[0m  viber00t.toml already exists")
		return
	}

	err := ioutil.WriteFile("viber00t.toml", []byte(defaultConfig), 0644)
	if err != nil {
		log.Fatal("\033[31m✗\033[0m Failed to create config:", err)
	}
	fmt.Println("\033[32m✓\033[0m Created viber00t.toml")
}

func initGlobalConfig() {
	configDir := filepath.Join(getXDGConfigHome(), "viber00t")
	configPath := filepath.Join(configDir, "config.toml")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatal("\033[31m✗\033[0m Failed to create config directory:", err)
	}

	// Check if global config already exists
	if _, err := os.Stat(configPath); err == nil {
		return
	}

	// Write default global config
	if err := ioutil.WriteFile(configPath, []byte(defaultGlobalConfig), 0644); err != nil {
		log.Fatal("\033[31m✗\033[0m Failed to create global config:", err)
	}

	fmt.Println("\033[32m✓\033[0m Created global config at ~/.config/viber00t/config.toml")
}

func loadGlobalConfig() (*GlobalConfig, error) {
	var config GlobalConfig
	configPath := filepath.Join(getXDGConfigHome(), "viber00t", "config.toml")

	// Initialize if not exists
	initGlobalConfig()

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(data), &config); err != nil {
		return nil, err
	}

	// Set defaults if not specified
	if config.DefaultAgent == "" {
		config.DefaultAgent = "claude"
	}
	if config.DefaultImage == "" {
		config.DefaultImage = "viber00t/base:latest"
	}
	if len(config.ClaudeFlags) == 0 {
		config.ClaudeFlags = []string{"--dangerously-skip-permissions"}
	}

	return &config, nil
}

func loadConfig() (*Config, error) {
	var config Config

	// Load global config for defaults
	globalConfig, _ := loadGlobalConfig()

	data, err := ioutil.ReadFile("viber00t.toml")
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(data), &config); err != nil {
		return nil, err
	}

	// Apply global defaults if not specified in project config
	if config.Project.Agent == "" {
		config.Project.Agent = globalConfig.DefaultAgent
	}

	// Add global default packages and envs
	if len(globalConfig.DefaultPackages) > 0 && len(config.Install) > 0 {
		config.Install[0].Packages = append(globalConfig.DefaultPackages, config.Install[0].Packages...)
	}
	if len(globalConfig.DefaultEnvs) > 0 && len(config.Install) > 0 {
		config.Install[0].Envs = append(globalConfig.DefaultEnvs, config.Install[0].Envs...)
	}

	return &config, nil
}

func getConfigHash(config *Config) string {
	// Create hash of entire config that affects the build
	h := sha256.New()
	h.Write([]byte(config.Project.Name))
	h.Write([]byte(config.Project.Agent))
	h.Write([]byte(fmt.Sprintf("%v", config.Project.Privileged)))

	// Hash install packages and envs
	if len(config.Install) > 0 {
		for _, pkg := range config.Install[0].Packages {
			h.Write([]byte(pkg))
		}
		for _, env := range config.Install[0].Envs {
			h.Write([]byte(env))
		}
	}

	// Also hash the config file modification time
	if info, err := os.Stat("viber00t.toml"); err == nil {
		h.Write([]byte(fmt.Sprintf("%d", info.ModTime().Unix())))
	}

	return hex.EncodeToString(h.Sum(nil))[:12]
}

func getProjectImageName(config *Config) string {
	hash := getConfigHash(config)
	return fmt.Sprintf("viber00t/%s:%s", config.Project.Name, hash)
}

func generateDockerfile(config *Config, globalConfig *GlobalConfig) string {
	// Separate base packages from project-specific ones
	var basePackages []string
	var projectPackages []string

	// Always add minimal required packages to base
	basePackages = append(basePackages, "curl", "wget", "sudo", "ca-certificates", "gnupg", "lsb-release", "git", "vim", "nano", "htop", "less", "man-db")

	// Add base packages from global config
	if len(globalConfig.BasePackages) > 0 {
		basePackages = append(basePackages, globalConfig.BasePackages...)
	}

	// Add global default packages to project layer
	if len(globalConfig.DefaultPackages) > 0 {
		projectPackages = append(projectPackages, globalConfig.DefaultPackages...)
	}

	// Add project-specific packages
	if len(config.Install) > 0 {
		if len(config.Install[0].Packages) > 0 {
			projectPackages = append(projectPackages, config.Install[0].Packages...)
		}

		// Expand environment templates
		for _, env := range config.Install[0].Envs {
			if packages, ok := envTemplates[env]; ok {
				projectPackages = append(projectPackages, packages...)
			}
		}

		// Add global default envs
		for _, env := range globalConfig.DefaultEnvs {
			if packages, ok := envTemplates[env]; ok {
				projectPackages = append(projectPackages, packages...)
			}
		}
	}

	// Build multi-stage Dockerfile
	dockerfile := `FROM ubuntu:latest AS base

# Set non-interactive frontend
ENV DEBIAN_FRONTEND=noninteractive

# Install base packages (cached layer)
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ` + strings.Join(basePackages, " \\\n    ") + ` && \
    rm -rf /var/lib/apt/lists/*

# Install Claude Code (cached layer)
RUN curl -fsSL https://claude.ai/install.sh | bash

FROM base AS final
`

	// Only add project packages if there are any
	if len(projectPackages) > 0 {
		dockerfile += `
# Install project-specific packages
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ` + strings.Join(projectPackages, " \\\n    ") + ` && \
    rm -rf /var/lib/apt/lists/*
`
	}

	// Add final configuration
	dockerfile += `
# Setup environment
ENV PATH="/root/.local/bin:${PATH}"
WORKDIR /c0de/project

CMD ["claude"]
`

	return dockerfile
}

func buildProjectImage(config *Config) error {
	imageName := getProjectImageName(config)
	currentHash := getConfigHash(config)
	globalConfig, _ := loadGlobalConfig()

	// Check state file for previous build
	stateDir := filepath.Join(getXDGStateHome(), "viber00t", "images")
	stateFile := filepath.Join(stateDir, config.Project.Name+".state")

	needsBuild := true

	// Check if we have a previous build state
	if data, err := ioutil.ReadFile(stateFile); err == nil {
		parts := strings.Split(string(data), ":")
		if len(parts) == 2 {
			oldHash := parts[1]
			if oldHash == currentHash {
				// Check if image still exists
				checkCmd := exec.Command("podman", "images", "-q", imageName)
				output, _ := checkCmd.Output()
				if len(output) > 0 {
					fmt.Printf("\033[35m◉\033[0m Using cached image: %s\n", imageName)
					needsBuild = false
				}
			} else {
				// Config changed, remove old image
				oldImage := parts[0] + ":" + oldHash
				fmt.Printf("\033[33m⟳\033[0m Config changed, removing old image: %s\n", oldImage)
				exec.Command("podman", "rmi", oldImage).Run()
			}
		}
	}

	if !needsBuild {
		return nil
	}

	// Clean up any existing containers using this image
	fmt.Printf("\033[33m⟳\033[0m Cleaning up existing containers...\n")
	exec.Command("podman", "ps", "-a", "--filter", fmt.Sprintf("ancestor=%s", imageName), "--format", "{{.Names}}", "|", "xargs", "-r", "podman", "rm", "-f").Run()

	// Remove the old image before building new one
	exec.Command("podman", "rmi", "-f", imageName).Run()

	fmt.Printf("\033[35m◉\033[0m Building project image: %s\n", imageName)

	// Generate Dockerfile from configs
	dockerfile := generateDockerfile(config, globalConfig)

	// Create temp build directory
	buildDir := filepath.Join(getXDGCacheHome(), "viber00t", "builds", config.Project.Name)
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("failed to create build directory: %w", err)
	}

	// Write Dockerfile
	dockerfilePath := filepath.Join(buildDir, "Dockerfile")
	if err := ioutil.WriteFile(dockerfilePath, []byte(dockerfile), 0644); err != nil {
		return fmt.Errorf("failed to write Dockerfile: %w", err)
	}

	// Build image
	cmd := exec.Command("podman", "build", "-t", imageName, buildDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build image %s: %w", imageName, err)
	}

	// Store build state with full image name
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	// Store as "imagename:hash" format
	stateData := fmt.Sprintf("viber00t/%s:%s", config.Project.Name, currentHash)
	ioutil.WriteFile(stateFile, []byte(stateData), 0644)

	return nil
}

func runContainer() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("\033[31m✗\033[0m No viber00t.toml found. Run 'viber00t init' first.")
		os.Exit(1)
	}

	// Build project-specific image
	if err := buildProjectImage(config); err != nil {
		log.Fatal("\033[31m✗\033[0m Failed to build image:", err)
	}

	cwd, _ := os.Getwd()
	containerName := fmt.Sprintf("viber00t-%s", filepath.Base(cwd))

	// Check if container already exists
	checkCmd := exec.Command("podman", "ps", "-a", "--format", "{{.Names}}")
	output, _ := checkCmd.Output()
	if strings.Contains(string(output), containerName) {
		fmt.Printf("\033[33m⟳\033[0m Removing existing container %s\n", containerName)
		exec.Command("podman", "rm", "-f", containerName).Run()
	}

	args := []string{
		"run", "-it",
		"--name", containerName,
		"--hostname", "viber00t",
		"--userns=keep-id:uid=0,gid=0",
		"-v", fmt.Sprintf("%s:/c0de/project", cwd),
	}

	// Mount Claude config directory if it exists
	claudeDir := filepath.Join(os.Getenv("HOME"), ".claude")
	if _, err := os.Stat(claudeDir); err == nil {
		args = append(args, "-v", fmt.Sprintf("%s:/root/.claude:rw", claudeDir))
	}

	// Mount claude.json config file if it exists
	claudeJSON := filepath.Join(os.Getenv("HOME"), ".claude.json")
	if _, err := os.Stat(claudeJSON); err == nil {
		args = append(args, "-v", fmt.Sprintf("%s:/root/.claude.json:rw", claudeJSON))
	}

	// Mount git config
	gitConfig := filepath.Join(os.Getenv("HOME"), ".gitconfig")
	if _, err := os.Stat(gitConfig); err == nil {
		args = append(args, "-v", fmt.Sprintf("%s:/root/.gitconfig:ro", gitConfig))
	}

	// Mount git credentials
	gitCreds := filepath.Join(os.Getenv("HOME"), ".git-credentials")
	if _, err := os.Stat(gitCreds); err == nil {
		args = append(args, "-v", fmt.Sprintf("%s:/root/.git-credentials:ro", gitCreds))
	}

	// Mount SSH keys for git
	sshDir := filepath.Join(os.Getenv("HOME"), ".ssh")
	if _, err := os.Stat(sshDir); err == nil {
		args = append(args, "-v", fmt.Sprintf("%s:/root/.ssh:ro", sshDir))
	}

	// Add privileged mode if requested
	if config.Project.Privileged {
		args = append(args, "--privileged", "--security-opt", "label=disable")
		// Mount docker socket if it exists
		if _, err := os.Stat("/var/run/docker.sock"); err == nil {
			args = append(args, "-v", "/var/run/docker.sock:/var/run/docker.sock")
		}
	}

	// Add volumes
	for _, vol := range config.Volumes {
		if vol.Source != "" && vol.Target != "" {
			source := expandPath(vol.Source)
			args = append(args, "-v", fmt.Sprintf("%s:%s:Z", source, vol.Target))
		}
	}

	// Add ports
	for _, port := range config.Ports {
		if port.Host != 0 && port.Container != 0 {
			args = append(args, "-p", fmt.Sprintf("%d:%d", port.Host, port.Container))
		}
	}

	// Environment variables
	args = append(args, "-e", "TERM=xterm-256color")
	args = append(args, "-e", "VIBER00T_PROJECT="+config.Project.Name)
	args = append(args, "-e", "IS_SANDBOX=true")

	// Create package install script if needed
	if len(config.Install) > 0 {
		var allPackages []string

		// Add explicit packages
		if len(config.Install[0].Packages) > 0 {
			allPackages = append(allPackages, config.Install[0].Packages...)
		}

		// Expand environment templates
		if len(config.Install[0].Envs) > 0 {
			for _, env := range config.Install[0].Envs {
				if packages, ok := envTemplates[env]; ok {
					allPackages = append(allPackages, packages...)
				}
			}
		}

		if len(allPackages) > 0 {
			packages := strings.Join(allPackages, " ")
			args = append(args, "-e", "VIBER00T_INSTALL="+packages)
		}
	}

	// Load global config for flags
	globalConfig, _ := loadGlobalConfig()

	// Use project-specific image
	imageName := getProjectImageName(config)
	args = append(args, imageName)

	// Run with specified agent and flags
	if config.Project.Agent != "" {
		agentCmd := []string{config.Project.Agent}

		// Add claude specific flags
		if config.Project.Agent == "claude" && len(globalConfig.ClaudeFlags) > 0 {
			agentCmd = append(agentCmd, globalConfig.ClaudeFlags...)
		}

		args = append(args, agentCmd...)
	}

	fmt.Printf("\033[35m◉\033[0m Starting viber00t for \033[36m%s\033[0m...\n", config.Project.Name)
	fmt.Println("\033[90m───────────────────────────────────\033[0m")

	cmd := exec.Command("podman", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal("\033[31m✗\033[0m Container failed:", err)
	}
}

func runShell() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("\033[31m✗\033[0m No viber00t.toml found. Run 'viber00t init' first.")
		os.Exit(1)
	}

	// Build project-specific image
	if err := buildProjectImage(config); err != nil {
		log.Fatal("\033[31m✗\033[0m Failed to build image:", err)
	}

	cwd, _ := os.Getwd()
	containerName := fmt.Sprintf("viber00t-shell-%s", filepath.Base(cwd))

	// Check if container already exists
	checkCmd := exec.Command("podman", "ps", "-a", "--format", "{{.Names}}")
	output, _ := checkCmd.Output()
	if strings.Contains(string(output), containerName) {
		fmt.Printf("\033[33m⟳\033[0m Removing existing container %s\n", containerName)
		exec.Command("podman", "rm", "-f", containerName).Run()
	}

	args := []string{
		"run", "-it",
		"--name", containerName,
		"--hostname", "viber00t",
		"--userns=keep-id:uid=0,gid=0",
		"-v", fmt.Sprintf("%s:/c0de/project", cwd),
	}

	// Mount Claude config directory if it exists
	claudeDir := filepath.Join(os.Getenv("HOME"), ".claude")
	if _, err := os.Stat(claudeDir); err == nil {
		args = append(args, "-v", fmt.Sprintf("%s:/root/.claude:rw", claudeDir))
	}

	// Mount claude.json config file if it exists
	claudeJSON := filepath.Join(os.Getenv("HOME"), ".claude.json")
	if _, err := os.Stat(claudeJSON); err == nil {
		args = append(args, "-v", fmt.Sprintf("%s:/root/.claude.json:rw", claudeJSON))
	}

	// Mount git config
	gitConfig := filepath.Join(os.Getenv("HOME"), ".gitconfig")
	if _, err := os.Stat(gitConfig); err == nil {
		args = append(args, "-v", fmt.Sprintf("%s:/root/.gitconfig:ro", gitConfig))
	}

	// Mount git credentials
	gitCreds := filepath.Join(os.Getenv("HOME"), ".git-credentials")
	if _, err := os.Stat(gitCreds); err == nil {
		args = append(args, "-v", fmt.Sprintf("%s:/root/.git-credentials:ro", gitCreds))
	}

	// Mount SSH keys for git
	sshDir := filepath.Join(os.Getenv("HOME"), ".ssh")
	if _, err := os.Stat(sshDir); err == nil {
		args = append(args, "-v", fmt.Sprintf("%s:/root/.ssh:ro", sshDir))
	}

	// Add privileged mode if requested
	if config.Project.Privileged {
		args = append(args, "--privileged", "--security-opt", "label=disable")
		// Mount docker socket if it exists
		if _, err := os.Stat("/var/run/docker.sock"); err == nil {
			args = append(args, "-v", "/var/run/docker.sock:/var/run/docker.sock")
		}
	}

	// Add volumes
	for _, vol := range config.Volumes {
		if vol.Source != "" && vol.Target != "" {
			source := expandPath(vol.Source)
			args = append(args, "-v", fmt.Sprintf("%s:%s:Z", source, vol.Target))
		}
	}

	// Add ports
	for _, port := range config.Ports {
		if port.Host != 0 && port.Container != 0 {
			args = append(args, "-p", fmt.Sprintf("%d:%d", port.Host, port.Container))
		}
	}

	// Environment variables
	args = append(args, "-e", "TERM=xterm-256color")
	args = append(args, "-e", "VIBER00T_PROJECT="+config.Project.Name)
	args = append(args, "-e", "IS_SANDBOX=true")

	// Use project-specific image
	imageName := getProjectImageName(config)
	args = append(args, imageName)

	// Override with bash
	args = append(args, "/bin/bash")

	fmt.Printf("\033[35m◉\033[0m Starting shell for \033[36m%s\033[0m...\n", config.Project.Name)
	fmt.Println("\033[90m───────────────────────────────────\033[0m")

	cmd := exec.Command("podman", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal("\033[31m✗\033[0m Shell failed:", err)
	}
}

func cleanImages() {
	fmt.Println("\033[35m◉\033[0m Cleaning viber00t images and cache...")

	// Remove all viber00t images
	cmd := exec.Command("podman", "images", "--format", "{{.Repository}}:{{.Tag}}", "--filter", "reference=viber00t/*")
	output, _ := cmd.Output()

	images := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, img := range images {
		if img != "" {
			fmt.Printf("\033[33m⟳\033[0m Removing image: %s\n", img)
			exec.Command("podman", "rmi", img).Run()
		}
	}

	// Clean cache directory
	cacheDir := filepath.Join(getXDGCacheHome(), "viber00t")
	if err := os.RemoveAll(cacheDir); err != nil {
		fmt.Printf("\033[33m⚠\033[0m  Failed to clean cache: %v\n", err)
	}

	// Clean state directory
	stateDir := filepath.Join(getXDGStateHome(), "viber00t")
	if err := os.RemoveAll(stateDir); err != nil {
		fmt.Printf("\033[33m⚠\033[0m  Failed to clean state: %v\n", err)
	}

	fmt.Println("\033[32m✓\033[0m Cleanup complete!")
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home := os.Getenv("HOME")
		return filepath.Join(home, path[2:])
	}
	return path
}
