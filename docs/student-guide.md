# Student Guide

Using Shell Gym takes two windows (e.g., a split-screen): a **terminal** on the training box and
the **gym page** in a browser at 127.0.0.1:63636. The page tells you what to do and
watches your progress; all actual work happens in the terminal.

## Doing a rep

1. **Read the exercise on the page.** It explains one small technique and gives a
   concrete goal.
2. **Switch to the terminal and do it for real.** Type the commands
   yourself - building that reflex is the whole point. Any interactive
   shell on the box is intercepted - there is nothing special to launch (beyond the `shellgym` daemon).
3. **Watch the check pass.** Verification is automatic and continuous.
   The moment the goal is reached, the task completes; when all tasks in
   an exercise are done, a check mark animation plays and (with
   auto-advance on) the next exercise slides in.

Some tasks verify a *current state* rather than a one-time action - for
example "the port is free" or "the process is gone". These can flip back
to unsatisfied if you undo the state before the whole exercise completes.

## Hints

- **Folded hints** - some exercises include collapsed hint boxes you can
  expand when stuck.
- **Dynamic hints** - task boxes can update with context-aware messages
  based on what the checks actually observe ("your shell is still in
  /home/laborant", "the file exists but is empty"). They appear
  automatically when the gym sees you struggling.

Hints point at what is wrong; they do not paste the solution.

## Moving around

When the browser page has the focus:

| Key | Action |
|---|---|
| `←` / `→` | previous / next exercise |
| `m` | open the path map (jump to any exercise) |
| `?` | help overlay |
| `Esc` | close overlays |

Auto-advance is on by default; toggle it in the toolbar if you prefer to
move manually. You can revisit completed exercises at any time - their
tasks stay completed.

## Progress and resuming

Progress is saved on the box continuously (at `/var/lib/shellgym` by default).
You can close the browser, log out, or stop for days - reopening the gym page
puts you back where you left off, with the same randomized values (directory
names, tokens, ports) as before.

To practice a path again from scratch, start it with a fresh state (e.g., `rm -rf /var/lib/shellgym`) -
the randomized parameters roll anew, so repeated runs stay honest instead of
becoming muscle memory for one specific answer.
