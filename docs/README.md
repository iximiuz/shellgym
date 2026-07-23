# Shell Gym Documentation

- [Student Guide](student-guide.md) - using the gym: the screen, doing
  reps, hints, navigation, progress and resuming.
- [Authoring Guide](authoring-guide.md) - writing learning paths: the
  content format, frontmatter reference, vars, init scripts, tasks,
  hints, solve scripts, guidelines, and the test workflow.
- [Built-in Checks Reference](checks.md) - every `wait_*`/helper
  command available to task scripts: semantics, flags, examples, and
  pitfalls.
- [Detection Mechanisms](detection.md) - how the daemon observes the
  student without instrumenting their shell: procfs shell discovery, the
  kernel proc connector, direct state polling, and the check API socket.
- [Design](design.md) - architecture: subsystems, CLI, state layout,
  web UI, HTTP API, deployment, and testing strategy.
