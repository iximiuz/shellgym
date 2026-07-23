# Shell Gym - project rules

## Build locally, test in a remote Linux playground

Never run `shellgym serve` (or `make run`, `shellgym solve`, or any e2e flow that
starts the daemon) on this machine. The daemon runs **as root** and the reps
mutate live system state - creating/removing files, killing processes, freeing
ports - so it must only run on a disposable host (see the CAUTION note in
README.md). Compile-time work stays local; runtime testing happens in an
[iximiuz Labs playground](https://labs.iximiuz.com/playgrounds) via `labctl`.

Safe to run locally: `make build`, `make test`, `make vet`, `make validate`
(validate only lints/renders a path, it doesn't execute anything).

Typical loop:

```sh
# 1. Build locally (CGO_ENABLED=0, so the binary is static and portable;
#    playground machines are linux/amd64 - same as this box, no cross-compile needed)
make build

# 2. Start a disposable playground once per session (prints the playground ID)
labctl playground start ubuntu-26-04 --quiet

# 3. Copy the binary and content over
labctl cp bin/shellgym <playground-id>:~/shellgym
labctl cp -r paths <playground-id>:~/paths

# 4. Run remotely (playground login user is "laborant")
labctl ssh <playground-id> -- sudo /home/laborant/shellgym serve --path /home/laborant/paths/sample-linux-101 --user laborant
# ...or open an interactive shell with: labctl ssh <playground-id>

# 5. Reach the web UI (daemon listens on 63636) from the local machine
labctl port-forward <playground-id> -L 63636:63636
```

Playgrounds are ephemeral; reuse one playground ID across iterations within a
session (just re-`cp` the rebuilt binary), and `labctl playground destroy <id>`
when done. `labctl playground list` shows recent sessions if the ID is lost.
