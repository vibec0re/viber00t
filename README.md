# viber00t 💀⚡

> fuck docker desktop. fuck vagrant. fuck your 47-step k8s setup.
> 
> **we're shipping code, not writing dissertations.**

```bash
./viber00t  # boom. you're in.
```

## wtf is this

containerized dev environments that don't suck. built on podman because docker is bloat. 

**one command. zero bullshit.**

## quickstart (30 seconds or your money back)

```bash
# clone it
git clone https://github.com/yourusername/viber00t
cd viber00t && go build -o viber00t

# init your project
cd ~/your-sick-project
~/viber00t/viber00t init

# fucking send it
~/viber00t/viber00t
```

congrats you're now coding inside a container like it's 2099 🎉

## config (optional, for control freaks)

`viber00t.toml` if you must:

```toml
[project]
name = "my-chaos-engine"
agent = "claude-code"  # or whatever ai overlord you worship
privileged = true      # docker-in-podman bc recursion is fun

[[install]]
packages = ["rust", "neovim", "chaos"]  # your weapons of choice

[[volumes]]
source = "~/secrets"
target = "/c0de/secrets"  # mount whatever tf you want

[[ports]]
host = 6969
container = 6969  # nice
```

## features that actually matter

- **instant containers** - no 10GB docker desktop eating your ram
- **auto-mounts everything** - project, ssh keys, ai creds, your soul
- **language templates** - rust/go/python/node/whatever
- **docker-in-podman** - because inception
- **zero config** - but configurable if you're into that
- **claude integration** - ai pair programming from the matrix

## why viber00t

- **docker desktop**: 4GB RAM for a whale icon? nah
- **vagrant**: it's 2024 not 2014
- **devcontainers**: microsoft complexity? pass
- **nix**: great if you have a PhD in functional programming
- **viber00t**: `./viber00t` and you're coding. that's it.

## philosophy

```
ship > plan
code > meetings  
vibes > process
chaos > order
```

## commands (all 3 of them)

```bash
./viber00t         # enter the matrix (run container)
./viber00t init    # create config (optional)
./viber00t clean   # nuke cached images
```

## requirements

- podman (not docker)
- go (to build)
- a pulse

## building

```bash
go build -o viber00t
# congrats you're a 10x developer now
```

## troubleshooting

**"it doesn't work"**
- try `./viber00t` harder

**"permission denied"**
- `chmod +x viber00t` genius

**"i need more features"**
- no you don't. ship your code.

**"the docs are unclear"**
- the code is 600 lines. read it.

## contributing

1. fork it
2. break it  
3. fix it
4. PR it
5. ship it

no 47-page contributing guidelines. no CLA. no bullshit.

if your code works and doesn't add complexity, it's probably getting merged.

## license

MIT because GPL is a novel and Apache needs a law degree

## credits

built with 🖤 by the vibec0re collective

powered by:
- too much caffeine ☕
- questionable life choices 🎲
- pure fucking willpower ⚡

---

*remember: every second you spend configuring your dev environment is a second you're not shipping code*

**fuck it, let's fucking gooooo** 🚀🚀🚀