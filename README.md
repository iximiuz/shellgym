# Shell Gym - an Interactive Linux Command-Line Trainer

Shell Gym is a background daemon with a built-in web UI that turns
any Linux box into an interactive command-line trainer.

"Learn the idea in a tutorial. Build the reflex in Shell Gym."

Tutorials explain concepts - Shell Gym drills them, helping you form the right muscle memory.
Open a split-screen: a completely ordinary terminal on the one side, and the Shell Gym UI
on the other. The UI shows small, fast-changing assignments - **reps** (in the traditional
gym sense). Each rep asks for one concrete action (enter a directory, create a file, kill a
process, free a port) and completes automatically the moment the system state changes.
There is no "check" button and no copy-paste: you will have to type real commands into a real shell,
and the gym trainer will observe your actions and guide you on the way.

## Why it exists

Reading about navigating the file tree, stdio redirection, or signals is not the same as being
able to perform these actions without thinking. That fluency comes only from hands-on practice
and repetition - and this is what Shell Gym provides.

## How it works

The student's shell is not modified in any way: no prompt hooks, no wrappers, no special shell functions.
All observation happens from the outside, meaning you work in the regular Linux terminal.

The Linux "magic" Shell Gym uses to achieve "zere instrumentation" observation:

- **procfs** - discovering the interactive shells, reading their working directories, scanning processes, files, and ports.
- **the kernel proc connector** - a netlink firehose of every `exec()` on the box, used to notice which commands the user runs.
- **plain state checks** - files existing, processes running, ports listening.

Because nothing is injected into the shell, the skills practiced in the
gym transfer one-to-one to any real terminal. See [detection.md](detection.md)
for how each mechanism works.

## Quick start

Shell Gym runs on a plain Linux host (a VM, a spare laptop, an EC2 instance, etc.) **as root**.
It should work on most (if not all) mainstream Linux distributions.

> [!CAUTION]
> Since reps will ask you to perform real actions on the live system, **use Shell Gym only with a disposable Linux host.** A few options:
> - Use a local VM (Lima, SlicerVM, VirtualBox, etc.)
> - Use an [iximiuz Labs Linux Playground](https://labs.iximiuz.com/playgrounds?category=linux&filter=official)
> - Use a DigitalOcean droplet, an EC2 instance, etc.

### Option 1: Download a release binary

Grab the latest release tarball (it bundles the `shellgym` binary and the
sample learning path) and unpack it:

```sh
arch=$(uname -m | sed 's/x86_64/amd64/; s/aarch64/arm64/')
curl -L "https://github.com/iximiuz/shellgym/releases/latest/download/shellgym_linux_${arch}.tar.gz" | tar xz
```

### Option 2: Build from source

Clone the repository and build the binary (requires Go):

```sh
make build
ln -s bin/shellgym shellgym
```

### Start the daemon

```sh
sudo ./shellgym serve --path "$PWD/paths/sample-linux-101" --user $USER
```

...or start the Shell Gym daemon in the background:

```sh
sudo systemd-run --unit=shellgym --collect \
    "$PWD/shellgym" serve --path "$PWD/paths/sample-linux-101" --user $USER
```

Once started, open the web UI in a browser and follow the learning path:

```sh
open http://127.0.0.1:63636
```

## Bring you own learning paths

Shell Gym defines a format, not a curriculum. The bundled [sample-linux-101](paths/sample-linux-101) path is the reference implementation.
A **learning path** is a directory tree that follows the following structure:

```
paths/<path>/              # path.yaml: id, title, user
    010.module-a/          # numeric prefix defines order
        module.md          # optional module intro (static)
        010.unit-x/
            unit.md        # a signle "rep" (tasks + checks)
    020.module-b/          # another module
        ...
```

- **Path** - the whole course (`path.yaml` + modules). One daemon serves one path.
- **Module** - a themed group of units with an optional intro scene (`module.md`).
- **Unit** - one rep: a markdown page (`unit.md`) with YAML frontmatter
  that defines setup scripts, verification tasks, hints, and a hidden
  reference solution (for testing).
- **Task** - one verifiable condition inside a unit. A unit completes when all of its tasks are met.

Units can be parametric (randomized directory names, tokens, ports),
depend on the state left behind by earlier units, and be filtered by
distro or host capabilities. Full format reference:
[authoring-guide.md](authoring-guide.md).

## Progress and resuming

Progress is persisted on disk (`/var/lib/shellgym` by default): the
user can stop any time and resume days later, surviving daemon
restarts and reboots of the box. Randomized parameters are sticky per
attempt, so a half-done unit looks the same after a resume.

However, since most learning paths will expect the student to modify
the state of the live system (e.g., create files, start processes, set env vars, etc.),
the successful resuming of the learning path may depend on the preservation
of system's state. For instance, if another process removes or modifes files
required by the learning path, it may break the restart, and the progress will be reset.

## Main CLI commands

| &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; Command &nbsp; &nbsp; &nbsp; &nbsp; | What it does |
|---|---|
| `shellgym serve` | The daemon itself: loads a path, runs the validation engine, serves the web UI |
| `shellgym validate` | Lints and renders a path without running it |
| `shellgym solve` | Auto-types reference solutions into a real pty shell, simulating a student pass |
| `shellgym skills` | Prints embedded authoring guides for AI-assisted content work |

Architecture details live in [design.md](design.md).
See also the [student guide](student-guide.md) for day-to-day usage and the [authoring guide](authoring-guide.md) for learning path creation.

## Documentation

- [Student Guide](docs/student-guide.md) - Using the gym: reps, hints, navigation, progress
- [Authoring Guide](docs/authoring-guide.md) - Writing learning paths: format, tasks, vars, testing
- [Checks Reference](docs/checks.md) - Every built-in `wait_*`/helper command in detail
- [Detection Mechanisms](docs/detection.md) - How the daemon observes the student's shell
- [Design](docs/design.md) - Architecture, subsystems, state, APIs, deployment

## Copyright

Copyright (c) 2026 Ivan Velichko ([iximiuz Labs](https://labs.iximiuz.com)).

Shell Gym is a part of the iximiuz Labs learning platform, licensed under the
[PolyForm Noncommercial License 1.0.0](LICENSE.md). You are welcome to use, modify, and share
Shell Gym for personal learning and other noncommercial purposes, but commercial use or
redistribution requires prior written permission. Commercial licenses are available on
request - contact ivan@iximiuz.com.

Contributions are welcome - see [CONTRIBUTING.md](CONTRIBUTING.md) for the ground rules.

Required Notice: Copyright (c) 2026 Ivan Velichko (https://labs.iximiuz.com)
