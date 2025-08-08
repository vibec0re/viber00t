# [[install]] config - the complete fucking guide 🔥

> because apparently we need docs for 2 fields. here we are.

## wtf is [[install]]

it's how you tell viber00t what packages to yeet into your container. that's it.

## the two sacred fields

### `packages` - raw apt packages

straight to `apt-get install`. no magic. no bullshit.

```toml
[[install]]
packages = ["neovim", "zsh", "fish", "cowsay", "sl", "fortune"]
```

want postgres? redis? some obscure lib from 2003? throw it in:

```toml
[[install]]
packages = [
  "postgresql-client",
  "redis-tools", 
  "libncurses5-dev",
  "libffi-dev",
  "libgmp-dev",
  "figlet",
  "lolcat"
]
```

### `envs` - language environment templates

pre-built package collections for languages. because typing 8 python packages is annoying.

```toml
[[install]]
envs = ["python", "rust", "node"]  # boom. full stack.
```

## available language templates

here's what each env actually installs:

### `python`
```
python3, python3-dev, python3-pip, python3-venv, 
pipx, poetry, pyenv, python3-setuptools
```

### `rust` 
```
rustc, cargo, rustfmt, rust-src, 
pkg-config, libssl-dev
```

### `node`
```
nodejs, npm, yarn, n
```

### `go`
```
golang, gopls
```

### `ruby`
```
ruby-full, ruby-dev, bundler, rbenv
```

### `java`
```
openjdk-17-jdk, maven, gradle
```

### `cpp`
```
clang, clang-tools, clang-format, cmake, 
ninja-build, ccache, gdb, valgrind
```

### `php`
```
php, php-cli, php-mbstring, php-xml, composer
```

### `dotnet`
```
dotnet-sdk-8.0, nuget
```

## mix and match (chaos mode)

you can have multiple `[[install]]` blocks. they stack. go wild:

```toml
# first block: languages
[[install]]
envs = ["rust", "python"]

# second block: databases
[[install]]
packages = ["postgresql-14", "mongodb", "redis"]

# third block: chaos tools
[[install]]
packages = ["nmap", "htop", "ncdu", "speedtest-cli"]

# fourth block: more languages why not
[[install]]
envs = ["go", "node"]
packages = ["yarn", "pnpm"]  # override node defaults
```

## real world examples

### web dev chad
```toml
[[install]]
envs = ["node", "python"]
packages = ["nginx", "certbot", "pm2"]
```

### rust maximalist
```toml
[[install]]
envs = ["rust"]
packages = ["llvm", "clang", "mold", "sccache", "cargo-watch", "cargo-edit"]
```

### data science nerd
```toml
[[install]]
envs = ["python"]
packages = ["jupyter", "pandoc", "texlive", "r-base", "julia"]
```

### devops masochist
```toml
[[install]]
packages = [
  "terraform",
  "ansible", 
  "kubectl",
  "helm",
  "vault",
  "consul",
  "nomad",
  "packer"
]
```

### full chaos (not recommended but respected)
```toml
[[install]]
envs = ["python", "rust", "node", "go", "ruby", "java", "cpp", "php", "dotnet"]
packages = ["everything", "in", "the", "ubuntu", "repos"]
```

## pro tips

1. **envs first, packages second** - templates give you the basics, packages add specifics
2. **don't overthink it** - if you need something, add it. if you don't, don't.
3. **global config exists** - check `~/.config/viber00t/config.toml` for base packages
4. **rebuilds are cached** - same config = same image hash = instant startup
5. **docker included by default** - it's in base_packages, don't re-add it

## common gotchas

**"my package isn't installing"**
- check if it exists: `apt-cache search <package>`
- check spelling (it's `nodejs` not `node`)

**"builds are slow"**
- yeah first build takes a minute. cached after that. deal with it.

**"i want conda/nvm/rvm"**
- those are package managers not packages. install them differently.

**"can i use alpine/arch packages?"**
- no. ubuntu/debian only. fork it if you want alpine.

## the philosophy

remember: you're configuring a container, not planning a wedding. 

- need postgres? add it.
- need rust? add it.  
- need both? add both.
- need nothing? leave it empty.

**ship code, not configs.**

## still confused?

the entire install logic is in `main.go:379-452`. it's like 70 lines. read it.

---

*if you needed this doc, you're overthinking it. just add packages and move on.* 🚀