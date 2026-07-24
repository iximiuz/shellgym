# Shell Gym Linux Curriculum

## Purpose

Shell Gym should offer a long-term catalog of short, hands-on Linux learning paths.

The curriculum starts with students who have essentially no Linux command-line experience and gradually prepares them for day-to-day work as DevOps engineers, SREs, platform engineers, and Linux system operators.

Shell Gym is not intended to replace:

- Explanatory tutorials.
- Longer exploratory labs.
- Open-ended troubleshooting challenges.
- Architecture or internals courses.
- Tasks that require substantial investigation before the student knows what to type.

Its unique role is to build command-line muscle memory while giving the student an initial practical understanding of how Linux behaves.

The core learning loop is:

1. Receive a small operational assignment.
2. Perform the action in an ordinary Linux shell.
3. Get immediate feedback when the expected system state appears.
4. Repeat the same underlying skill in varied, realistic situations.
5. Reuse the skill later without being explicitly told that it is being reviewed.

---

# Curriculum structure

The initial curriculum should contain approximately **20–30 learning paths**.

This proposal contains **28 paths**.

Each path should normally:

- Take approximately 15–30 minutes to complete.
- Contain approximately 15–30 reps.
- Focus on one coherent operational skill area.
- Introduce no more than 3–5 new commands.
- Prefer introducing only 2–3 commands when possible.
- Reuse commands and concepts introduced in earlier paths.
- Contain more reps when the operations are especially simple.
- Contain fewer reps when the operations require more reasoning or take longer to execute.
- End with several reps that combine the new skill with earlier skills.

The first part of the curriculum should be mostly linear. Once the student has basic terminal, filesystem, process, permission, package, service, and logging skills, the roadmap may branch into:

- Networking.
- Storage.
- Kernel and host internals.
- Performance.
- Remote operations.
- Automation.
- Troubleshooting.

The learning paths can later be stitched into a visual roadmap with prerequisites and optional branches.

---

# Recommended roadmap

## Linear foundation

The following paths should generally be completed in order:

1. Meet the Linux Shell
2. Control the Shell and Ask for Help
3. Navigate the Filesystem
4. List and Create Files
5. Copy, Move, Remove, and Link Files
6. Read and Inspect File Contents
7. Redirect Streams and Build Pipelines
8. Search Files and Text
9. Transform Text and Records
10. Control the Shell Environment
11. Work with Processes and Jobs
12. Work with Users, Groups, and Permissions
13. Use Advanced Linux Access Controls
14. Install and Inspect Software Packages
15. Operate Services, Logs, and Scheduled Jobs

After these foundations, the curriculum may branch.

## Networking branch

16. Inspect Network Configuration and Connectivity  
17. Work with Sockets and Network Services  
18. Route, Filter, and Isolate Network Traffic  

## Storage branch

19. Inspect Disk Usage, Devices, and Mounts  
20. Manage Filesystems, Partitions, and Logical Storage  

## Host internals and performance branch

21. Inspect the Linux Host and Kernel Interfaces  
22. Work with Isolation and Resource Controls  
23. Observe Performance and Triage a Host  

## Remote operations and automation branch

24. Operate Remote Systems with SSH  
25. Write Small Shell Scripts  
26. Build Safer Command-Line Automation  

## Applied operations branch

27. Troubleshoot Commands, Permissions, Services, and Networking  
28. Troubleshoot Storage, Resources, and Host State  

---

# 1. Meet the Linux Shell

## Goal

Introduce the terminal, the command execution model, and the basic structure of shell commands.

The student should finish the path comfortable entering commands and interpreting the relationship between the prompt, command, arguments, options, output, and exit status.

## Topics

- Recognizing the shell prompt.
- Entering a command and waiting for it to finish.
- Distinguishing command output from the next prompt.
- Commands, positional arguments, and options.
- Short options and long options.
- Combining compatible short options.
- Passing multiple arguments.
- Using `--` to mark the end of options.
- Understanding that whitespace separates arguments.
- Recognizing that quoting can preserve whitespace inside one argument.
- Running multiple commands with `;`.
- Running a command only after another succeeds with `&&`.
- Running a fallback command with `||`.
- Reading the previous command's exit status through `$?`.
- Recognizing successful and unsuccessful commands.
- Understanding that no output does not necessarily mean failure.
- Identifying the current user, host, terminal, and directory.
- Distinguishing shell builtins from external programs at a high level.

## Likely commands and shell features

- `echo`
- `printf`
- `true`
- `false`
- `whoami`
- `hostname`
- `tty`
- `pwd`
- `$?`
- `;`
- `&&`
- `||`
- `--`

Not every command listed here needs to count as a separately taught command. Some are merely simple tools for creating observable reps.

## Helpful task-page tips

Some terminal interaction is difficult to verify directly. Relevant tips should appear alongside tasks rather than being omitted from the curriculum:

- Use the Up and Down arrow keys to revisit command history.
- Use Left and Right arrows to edit a command.
- Use `Ctrl-A` and `Ctrl-E` to jump to the beginning and end of the line.
- Use `Ctrl-U` and `Ctrl-K` to remove parts of the current line.
- Use `Ctrl-L` to clear the visible terminal.
- Use Tab completion instead of typing long paths manually.

These tips can be introduced gradually and repeated where they become useful.

---

# 2. Control the Shell and Ask for Help

## Goal

Teach the student how to inspect commands, obtain documentation, interrupt work, and manage simple foreground and suspended commands.

## Topics

- Finding whether a command exists.
- Finding which executable or builtin will run.
- Distinguishing aliases, functions, builtins, and external programs.
- Using command-specific `--help`.
- Using shell builtin help.
- Opening manual pages.
- Searching manual page names and descriptions.
- Understanding manual page sections at a basic level.
- Searching inside a manual page.
- Quitting a pager.
- Interrupting a foreground command.
- Suspending a foreground command.
- Resuming a suspended command in the foreground.
- Recognizing that `Ctrl-C` sends an interrupt rather than closing the terminal.
- Recognizing that `Ctrl-Z` suspends rather than terminates.
- Recovering from a command that appears to hang.
- Checking whether the previous command succeeded.
- Finding basic information without being given exact syntax.

## Likely commands and shell features

- `type`
- `command -v`
- `help`
- `man`
- `apropos`
- `fg`
- `Ctrl-C`
- `Ctrl-Z`

## Helpful task-page tips

- `/pattern` searches inside many pagers and manual pages.
- `n` usually moves to the next search result.
- `q` exits `less` and most manual-page views.
- `Ctrl-R` searches shell history in reverse.
- Tab completion can also complete command names.

The path does not need to verify every pager or line-editing keystroke. It should verify the resulting command execution where possible and present the interaction details as contextual tips.

---

# 3. Navigate the Filesystem

## Goal

Build strong muscle memory for moving through the Linux filesystem.

This path may contain more than 30 reps because most navigation operations are quick.

## Topics

- Printing the current working directory.
- Recognizing an absolute path.
- Recognizing a relative path.
- Changing to an absolute path.
- Changing to a relative child directory.
- Moving through several nested directories.
- Referring to the current directory with `.`.
- Referring to the parent directory with `..`.
- Moving through multiple parent directories.
- Returning to the home directory with no argument.
- Referring to the home directory with `~`.
- Using `$HOME`.
- Returning to the previous working directory with `cd -`.
- Moving repeatedly between two locations.
- Navigating to paths containing spaces.
- Navigating to paths containing shell metacharacters.
- Using quoted paths.
- Using escaped pathnames.
- Using tab completion for long or awkward names.
- Understanding that the working directory belongs to a shell process.
- Recognizing that separate shells may have different working directories.
- Navigating from an unknown starting location.
- Following a chain of filesystem clues.
- Returning to a known location after exploring elsewhere.

## Likely commands and shell features

- `pwd`
- `cd`
- `.`
- `..`
- `~`
- `$HOME`
- `cd -`

## Rep design

Avoid a sequence of nearly identical tasks such as “change into directory A,” “change into directory B,” and “change into directory C.”

Use varied miniature scenarios:

- Enter a project source directory.
- Move from an application directory to its logs.
- Return to the user's home directory.
- Jump between a configuration directory and a data directory.
- Follow relative paths provided by another command.
- Navigate directories with awkward names.
- Recover after intentionally entering the wrong directory.

---

# 4. List and Create Files

## Goal

Teach the student to inspect directory contents and create basic filesystem structures.

## Topics

- Listing the current directory.
- Listing another directory without entering it.
- Listing multiple paths.
- Showing hidden files.
- Showing long-format metadata.
- Interpreting file type indicators.
- Reading file owner and group.
- Reading permission fields at a superficial level.
- Reading sizes and timestamps.
- Using human-readable sizes.
- Sorting listings by name, size, or modification time.
- Distinguishing files from directories.
- Creating an empty file.
- Updating a file's timestamp.
- Creating one directory.
- Creating nested directories.
- Creating several paths.
- Creating files in the current directory.
- Creating files using absolute paths.
- Creating files with spaces in their names.
- Creating hidden files and directories.
- Understanding pathname expansion by the shell.
- Using `*`.
- Using `?`.
- Using character classes.
- Seeing what a glob expands to before using it destructively.
- Understanding that ordinary globs do not normally include hidden files.
- Handling names that begin with `-`.

## Likely commands

- `ls`
- `touch`
- `mkdir`

## Rep design

Use small scenes such as:

- Preparing an application directory.
- Creating a configuration tree.
- Creating a directory for daily reports.
- Finding the newest or largest file in a listing.
- Revealing an unexpected hidden file.
- Creating a filename that requires quoting.

---

# 5. Copy, Move, Remove, and Link Files

## Goal

Build confidence manipulating files and directories while avoiding common destructive mistakes.

## Topics

- Copying a file to another filename.
- Copying a file into a directory.
- Copying several files.
- Copying a directory tree recursively.
- Understanding the difference between copying a directory and copying its contents.
- Overwriting an existing destination.
- Preserving metadata where relevant.
- Renaming a file.
- Moving a file into another directory.
- Moving a directory tree.
- Moving several matching files.
- Removing one file.
- Removing an empty directory.
- Removing a non-empty directory tree.
- Removing several files selected by a pattern.
- Previewing matches before removing them.
- Handling a filename beginning with `-`.
- Avoiding accidental recursive removal of the wrong location.
- Identifying a file's type.
- Inspecting detailed file metadata.
- Creating symbolic links.
- Reading symbolic-link targets.
- Recognizing broken symbolic links.
- Creating hard links.
- Recognizing that hard links share an inode.
- Understanding that deleting one hard-link name does not remove the underlying file while another link remains.
- Understanding the difference between a copied file and a linked file.

## Likely commands

- `cp`
- `mv`
- `rm`
- `rmdir`
- `ln`

Supporting commands such as `file` and `stat` may be introduced here or reused from task hints.

## Rep design

Include multiple exercises where destination semantics matter. Many students can run `cp`, but still hesitate over whether a destination path means:

- A new filename.
- An existing directory.
- A directory to be copied.
- The contents of a directory.

This distinction deserves deliberate repetition.

---

# 6. Read and Inspect File Contents

## Goal

Teach the student to choose an appropriate way to inspect files based on their size and purpose.

## Topics

- Printing a small file.
- Printing several files.
- Concatenating files in a specific order.
- Reading only the first lines.
- Reading only the last lines.
- Selecting a specific number of lines.
- Browsing a large file interactively.
- Searching while browsing a large file.
- Following a growing log.
- Recognizing newly appended records.
- Stopping a follow operation.
- Counting lines.
- Counting words.
- Counting bytes.
- Numbering lines.
- Comparing two versions of a file.
- Recognizing added, removed, and changed lines.
- Inspecting an unknown collection of files.
- Choosing between full output, partial output, and interactive browsing.
- Inspecting a file's type before attempting to read it.
- Avoiding dumping a very large or binary file directly into the terminal.
- Combining navigation, listing, and content inspection.

## Likely commands

- `cat`
- `head`
- `tail`
- `less`
- `wc`

Supporting commands:

- `nl`
- `diff`
- `file`

## Helpful task-page tips

Interactive `less` behavior can be introduced through tips:

- `/text` searches forward.
- `n` finds the next match.
- `g` and `G` move to the beginning and end.
- `q` exits.
- `F` follows a growing file in many `less` versions.

Shell Gym may verify that the correct file was opened with the right tool even if every pager keystroke cannot be observed.

---

# 7. Redirect Streams and Build Pipelines

## Goal

Develop an operational understanding of stdin, stdout, stderr, redirection, and pipelines.

## Topics

- Saving stdout to a file.
- Replacing an existing file with redirected output.
- Appending output to an existing file.
- Feeding a file into stdin.
- Understanding the difference between a filename argument and stdin.
- Redirecting stderr separately.
- Capturing stdout and stderr in different files.
- Combining stderr with stdout.
- Discarding selected output.
- Recognizing that redirection is performed by the shell.
- Connecting stdout from one command to stdin of another.
- Building a two-command pipeline.
- Building a longer pipeline.
- Counting records produced by another command.
- Filtering command output before saving it.
- Saving output while still displaying it.
- Appending through `tee`.
- Recognizing that a pipeline represents several concurrently connected processes.
- Recognizing failure inside a pipeline.
- Introducing `pipefail` at a practical level.
- Avoiding unnecessary temporary files.
- Combining file input, pipes, and output redirection.
- Separating useful output from diagnostics.

## Likely commands and shell features

- `>`
- `>>`
- `<`
- `2>`
- `2>&1`
- `|`
- `tee`

Supporting commands such as `cat`, `grep`, `sort`, and `wc` should be reused rather than retaught.

## Rep design

Prefer effects over syntax-recitation tasks.

Good:

- Save only the error messages from a failing command.
- Count matching records from a generated stream.
- Display and save a filtered list.
- Append a new diagnostic result without losing the previous one.

Weak:

- “Run a command using `2>`.”
- “Use a pipe.”

---

# 8. Search Files and Text

## Goal

Teach the student to locate information in file contents and locate files by their metadata.

## Topics

### Searching text

- Finding literal text.
- Searching case-insensitively.
- Searching several files.
- Searching recursively.
- Printing filenames containing matches.
- Printing line numbers.
- Inverting a match.
- Counting matches.
- Showing lines before and after a match.
- Searching with basic regular expressions.
- Searching with extended regular expressions.
- Matching the beginning or end of a line.
- Avoiding accidental interpretation of a literal string as a regular expression.

### Searching the filesystem

- Searching from a chosen root.
- Finding files by exact name.
- Finding files using name patterns.
- Finding only regular files.
- Finding only directories.
- Finding symbolic links.
- Finding files by size.
- Finding files by modification time.
- Finding files newer than a reference file.
- Finding files by permissions.
- Combining several predicates.
- Running an action for each result.
- Handling filenames containing whitespace.
- Handling filenames containing newlines using null delimiters.
- Passing results safely to another command.
- Understanding when `find -exec` is preferable to `xargs`.

## Likely commands

- `grep`
- `find`
- `xargs`

## Rep design

Use scenarios such as:

- Find a request ID in rotated logs.
- Exclude health-check requests.
- Find the configuration file that defines a port.
- Locate recently modified files.
- Find unexpectedly executable files.
- Remove only matching temporary files.
- Find files and then inspect their contents.
- Find files and move them to an archive directory.

---

# 9. Transform Text and Records

## Goal

Teach small, composable transformations commonly used in operational command lines.

This is not intended to become a complete `sed` or `awk` programming course.

## Topics

- Sorting lines alphabetically.
- Sorting numerically.
- Sorting by a selected field.
- Reversing sort order.
- Removing adjacent duplicates.
- Counting repeated records.
- Extracting delimited fields.
- Extracting character ranges.
- Replacing or deleting individual characters.
- Converting case.
- Joining corresponding lines.
- Selecting lines with `sed`.
- Replacing a simple pattern with `sed`.
- Removing selected lines.
- Printing selected fields with `awk`.
- Filtering records with a simple `awk` condition.
- Calculating a small aggregate with `awk`.
- Handling different field separators.
- Combining extraction, filtering, sorting, and counting.
- Recognizing malformed or unexpected records.
- Preserving headers while transforming data.
- Producing a small operational report from text input.

## Likely commands

- `sort`
- `uniq`
- `cut`
- `tr`
- `paste`

Introductory use of:

- `sed`
- `awk`

A path may technically contain more than five command names because several are tiny filters. The author should still minimize simultaneous novelty and reuse commands across many reps.

---

# 10. Control the Shell Environment

## Goal

Teach variables, environment inheritance, quoting, expansion, command lookup, and basic interactive shell configuration.

Because this area contains many subtle ideas, it may use more reps than an average path or be split into two paths during implementation.

## Topics

### Variables and environment

- Assigning a shell variable.
- Reading a shell variable.
- Using `${NAME}` to delimit a variable reference.
- Distinguishing shell variables from environment variables.
- Exporting a variable.
- Observing inheritance by child processes.
- Inspecting the environment.
- Temporarily setting a variable for one command.
- Unsetting a variable.

### Quoting and expansion

- Understanding unquoted expansion.
- Using single quotes.
- Using double quotes.
- Escaping one character.
- Preserving spaces inside one argument.
- Preventing pathname expansion.
- Allowing variable expansion while preserving whitespace.
- Using command substitution with `$(...)`.
- Using brace expansion.
- Combining variable, command, and pathname expansion.
- Avoiding obsolete backtick syntax.
- Understanding why unquoted command substitution is dangerous.

### Command grouping and lookup

- Running commands in a subshell.
- Grouping commands in the current shell.
- Observing cwd isolation in a subshell.
- Observing variable isolation in a subshell.
- Inspecting `PATH`.
- Finding the command selected by `PATH`.
- Temporarily extending `PATH`.
- Understanding command shadowing.
- Recognizing the risks of placing the current directory early in `PATH`.

### Interactive configuration

- Reading a shell startup file.
- Adding a simple alias.
- Adding a small shell function.
- Reloading configuration.
- Understanding login versus interactive startup at a high level.

## Likely commands and features

- Variable assignment
- `export`
- `env`
- `unset`
- Single and double quotes
- `$(...)`
- Parentheses and command groups
- `PATH`
- `source`

## Rep design

Quoting deserves heavy repetition. Reps should use realistic awkward values:

- Filenames with spaces.
- Variables containing several words.
- Literal `*` characters.
- Dollar signs that must not expand.
- Command output containing whitespace.
- Paths assembled from several variables.

---

# 11. Work with Processes and Jobs

## Goal

Teach the student to inspect, identify, control, and reason about running processes.

## Topics

- Listing processes.
- Reading PIDs.
- Reading process owners.
- Reading process state.
- Reading command lines.
- Filtering processes by user.
- Finding a process by name.
- Finding a process by a command-line pattern.
- Inspecting parent and child relationships.
- Displaying a process tree.
- Understanding PPID.
- Sending a termination signal.
- Sending an interrupt signal.
- Sending a hangup signal.
- Using SIGKILL only when a process cannot be stopped normally.
- Understanding that signal delivery does not always imply immediate termination.
- Starting a background job.
- Listing shell jobs.
- Bringing a job to the foreground.
- Resuming a suspended job in the background.
- Referring to jobs by job specifier.
- Keeping work running after the shell exits.
- Redirecting output for detached work.
- Using `nohup`.
- Using `disown`.
- Starting a process with adjusted niceness.
- Changing the niceness of a running process.
- Inspecting open files and sockets.
- Finding which process uses a file.
- Inspecting file descriptors through `/proc/<pid>/fd`.
- Inspecting basic resource limits.
- Connecting shell jobs to operating-system processes.

## Likely commands

- `ps`
- `pgrep`
- `kill`
- `jobs`
- `fg` and `bg`

Supporting commands:

- `pstree`
- `pkill`
- `nohup`
- `nice`
- `renice`
- `lsof`
- `fuser`
- `ulimit`

---

# 12. Work with Users, Groups, and Permissions

## Goal

Teach the ordinary Unix identity and permission model until the student can apply it without guessing.

## Topics

### Identity

- Inspecting the current UID and GIDs.
- Inspecting another user's identity.
- Listing group membership.
- Reading `/etc/passwd`.
- Reading `/etc/group`.
- Understanding usernames versus numeric IDs.
- Understanding the primary group and supplementary groups.
- Running a command with `sudo`.
- Running a command as another user.
- Starting a login shell as another user.
- Understanding the difference between `sudo`, `su`, and `runuser`.

### Ownership

- Reading file ownership.
- Changing the owner.
- Changing the group.
- Changing both owner and group.
- Recursively changing ownership.
- Avoiding unintended recursive ownership changes.

### Permission modes

- Reading owner, group, and other permission bits.
- Understanding read, write, and execute on regular files.
- Understanding read, write, and execute on directories.
- Distinguishing directory listing from directory traversal.
- Understanding why deleting a file is controlled by its parent directory.
- Changing permissions symbolically.
- Changing permissions numerically.
- Adding and removing one permission without disturbing others.
- Applying permissions recursively.
- Avoiding unsafe broad modes such as `777`.
- Understanding effective access based on identity and mode bits.

### Default permissions

- Inspecting `umask`.
- Predicting default file and directory permissions.
- Temporarily changing `umask`.
- Creating private files.
- Creating group-shared files.

## Likely commands

- `id`
- `groups`
- `sudo`
- `chown`
- `chmod`

Supporting commands:

- `su`
- `runuser`
- `chgrp`
- `umask`

## Rep design

Directory permissions need dedicated reps. They are commonly misunderstood even by students who can decode `rwx` on regular files.

---

# 13. Use Advanced Linux Access Controls

## Goal

Extend the ordinary permission model with mechanisms commonly encountered in real systems.

The path should remain operational. Deep historical and kernel-level explanations belong in tutorials.

## Topics

- Sticky-bit behavior in shared directories.
- Setgid behavior on shared directories.
- Inherited group ownership.
- Setuid executables.
- Setgid executables.
- Recognizing the security implications of setuid programs.
- Finding files with special mode bits.
- Reading POSIX ACLs.
- Adding an ACL entry.
- Removing an ACL entry.
- Understanding ACL masks.
- Creating default directory ACLs.
- Recognizing that ACLs extend rather than replace ordinary mode bits.
- Inspecting file capabilities.
- Granting a narrowly scoped file capability in an isolated environment.
- Removing file capabilities.
- Understanding capabilities as divisions of traditional root privilege.
- Inspecting the capabilities of a process.
- Recognizing the presence of AppArmor.
- Inspecting an AppArmor profile or denial.
- Recognizing the presence and enforcement state of SELinux.
- Inspecting labels and denials at an introductory level.
- Distinguishing discretionary and mandatory access control.
- Diagnosing whether ordinary permissions, ACLs, capabilities, AppArmor, or SELinux are blocking an operation.

## Likely commands

- `chmod`
- `getfacl`
- `setfacl`
- `getcap`
- `setcap`

Supporting tools may include:

- `find`
- `capsh`
- `aa-status`
- `getenforce`
- `ls -Z`

Distro-specific units should use `labels` and `requires` where needed.

---

# 14. Install and Inspect Software Packages

## Goal

Teach the student to identify the host distribution and manage installed software.

Package-manager-specific units should have distro-filtered variants.

## Topics

- Reading `/etc/os-release`.
- Identifying distribution ID and family.
- Identifying machine architecture.
- Identifying kernel architecture separately from distribution version.
- Refreshing package metadata.
- Searching for a package.
- Inspecting package information.
- Installing a package.
- Removing a package.
- Recognizing installed versus available versions.
- Finding which files a package installed.
- Finding which package owns a particular file.
- Verifying whether a package is installed.
- Listing installed packages matching a pattern.
- Inspecting package dependencies.
- Finding the executable installed by a package.
- Distinguishing a package name from a command name.
- Recognizing repository and package-manager errors.
- Performing a safe package-manager dry run where supported.
- Understanding the difference between removing a package and purging configuration.
- Recognizing when a package upgrade requires a service restart.

## Likely command families

### Debian and Ubuntu

- `apt`
- `apt-cache`
- `dpkg`
- `dpkg-query`

### Fedora, Rocky, RHEL, and similar systems

- `dnf`
- `rpm`

The learning outcomes should stay consistent even when the commands differ.

---

# 15. Operate Services, Logs, and Scheduled Jobs

## Goal

Teach day-to-day service management, log inspection, and recurring task scheduling on systemd-based hosts.

## Topics

### Services

- Checking whether a service is active.
- Reading detailed service status.
- Starting a service.
- Stopping a service.
- Restarting a service.
- Reloading a service.
- Understanding the difference between restart and reload.
- Enabling a service at boot.
- Disabling a service at boot.
- Distinguishing enabled from currently active.
- Finding failed units.
- Resetting failed state where appropriate.
- Reading a unit file.
- Inspecting unit properties.
- Inspecting unit dependencies.
- Recognizing service ordering.
- Creating or inspecting a drop-in override.
- Reloading the systemd manager after changing unit configuration.
- Validating whether a service uses socket activation.
- Recognizing when a service starts manually but fails under systemd.

### Logs and boot state

- Reading logs for one unit.
- Following new journal entries.
- Filtering by priority.
- Filtering by time range.
- Reading logs from the current boot.
- Reading logs from a previous boot.
- Reading kernel messages.
- Finding the cause of a failed unit.
- Finding startup-ordering problems.
- Distinguishing application output from systemd diagnostics.
- Recognizing rate-limited or rotated logs.
- Using logs to validate a service restart.

### Scheduled work

- Listing cron jobs.
- Creating a basic cron entry.
- Understanding cron's reduced environment.
- Listing systemd timers.
- Inspecting the next timer run.
- Connecting a timer to its service.
- Running a scheduled service manually.
- Recognizing missed or failed scheduled executions.

## Likely commands

- `systemctl`
- `journalctl`
- `dmesg`
- `crontab`
- `systemd-analyze`

Supporting commands may include `systemctl list-timers`.

---

# 16. Inspect Network Configuration and Connectivity

## Goal

Teach the student to inspect a host's network identity and diagnose connectivity layer by layer.

## Topics

### Interfaces and addresses

- Listing network interfaces.
- Reading interface state.
- Reading MAC addresses.
- Reading IPv4 and IPv6 addresses.
- Distinguishing loopback from external interfaces.
- Bringing an isolated interface up or down.
- Changing an interface MTU in a safe scene.
- Recognizing address scope.
- Identifying the likely outbound interface.

### Routes and neighbors

- Reading the routing table.
- Identifying a connected route.
- Identifying the default route.
- Determining which route will be used for a destination.
- Reading route metrics.
- Inspecting the neighbor table.
- Recognizing reachable, stale, and failed neighbors.
- Understanding ARP and NDP at an operational level.

### Name resolution

- Resolving a hostname through the system resolver.
- Querying DNS directly.
- Inspecting resolver configuration.
- Recognizing multiple resolver sources.
- Using `/etc/hosts`.
- Creating a temporary local hostname override.
- Understanding lookup order at a practical level.
- Separating DNS failure from connectivity failure.

### Connectivity tests

- Testing basic reachability.
- Recognizing that ping is not a complete service test.
- Tracing the network path.
- Recognizing the effect of missing ICMP responses.
- Testing whether a TCP service is reachable.
- Distinguishing link, address, route, DNS, transport, and application failures.
- Following a consistent layer-by-layer diagnostic workflow.

## Likely commands

- `ip`
- `getent`
- `dig` or `host`
- `ping`
- `tracepath` or `traceroute`

---

# 17. Work with Sockets and Network Services

## Goal

Connect application-level network operations to processes, sockets, addresses, and ports.

## Topics

### HTTP and TLS

- Making a basic HTTP request.
- Requesting headers only.
- Inspecting response headers.
- Reading status codes.
- Following redirects.
- Sending a different HTTP method.
- Sending a request body.
- Sending request headers.
- Inspecting connection errors.
- Inspecting TLS certificate information.
- Recognizing hostname or certificate validation failures.

### TCP and UDP

- Connecting to a TCP service manually.
- Creating a temporary TCP listener.
- Sending data between two terminals.
- Recognizing client and server endpoints.
- Sending and receiving UDP datagrams.
- Understanding connectionless behavior.
- Handling timeout behavior.

### Socket inspection

- Listing listening sockets.
- Listing connected sockets.
- Reading local and remote addresses.
- Reading socket states.
- Finding which process owns a port.
- Connecting a process PID to its socket.
- Distinguishing loopback, wildcard, and specific-address listeners.
- Diagnosing a service bound to the wrong address.
- Recognizing IPv4 versus IPv6 listeners.
- Recognizing socket activation.

### Packet capture

- Capturing traffic on a selected interface.
- Filtering by host.
- Filtering by port.
- Filtering by protocol.
- Watching a TCP handshake.
- Connecting a packet capture to a generated request.
- Avoiding unbounded packet captures.

## Likely commands

- `curl`
- `nc`
- `ss`
- `lsof`
- `tcpdump`

---

# 18. Route, Filter, and Isolate Network Traffic

## Goal

Introduce host routing, firewalling, NAT, and Linux virtual networking through small isolated scenes.

## Topics

### Routing and forwarding

- Adding a temporary route.
- Removing a route.
- Adding a host route.
- Adding a route through a gateway.
- Inspecting route selection.
- Enabling IP forwarding.
- Reading forwarding-related sysctls.
- Distinguishing local delivery from forwarding.

### Firewalling

- Listing nftables tables and chains.
- Reading chain hooks and priorities.
- Reading rule counters.
- Understanding rule order.
- Adding a narrowly scoped allow rule.
- Adding a narrowly scoped drop rule.
- Matching addresses and ports.
- Validating a rule through generated traffic.
- Removing a rule safely.
- Avoiding accidental host lockout.
- Understanding stateful filtering at an introductory level.

### NAT

- Recognizing source NAT.
- Configuring masquerading in an isolated topology.
- Recognizing destination NAT.
- Redirecting traffic to another local port.
- Inspecting whether translated traffic reaches its destination.
- Connecting forwarding, filtering, and NAT.

### Network namespaces and virtual links

- Creating a network namespace.
- Running a command inside it.
- Listing namespace-local interfaces.
- Creating a veth pair.
- Moving one endpoint into a namespace.
- Assigning addresses.
- Bringing links up.
- Connecting namespaces through a bridge.
- Adding routes between isolated networks.
- Deleting the topology cleanly.
- Repairing a deliberately broken virtual network.

## Likely commands

- `ip`
- `sysctl`
- `nft`

Supporting tools:

- `ping`
- `ss`
- `tcpdump`

This path may contain fewer than 15 reps if each topology-building rep is substantial. Alternatively, it can be divided into a routing/firewall path and a namespaces/virtual-networking path.

---

# 19. Inspect Disk Usage, Devices, and Mounts

## Goal

Teach the student to understand where storage is consumed and how block devices become mounted filesystems.

## Topics

### Space usage

- Inspecting filesystem capacity.
- Reading used and available space.
- Using human-readable units.
- Measuring a directory tree.
- Comparing `df` and `du`.
- Recognizing mount boundaries.
- Finding large directories.
- Finding large files.
- Recognizing deleted-but-open files.
- Inspecting inode usage.
- Diagnosing inode exhaustion.
- Recognizing many-small-file problems.

### Block devices

- Listing disks and partitions.
- Reading device size and type.
- Understanding parent-child device relationships.
- Reading filesystem type and UUID.
- Distinguishing a block device from a mounted filesystem.
- Recognizing loop devices.
- Recognizing device-mapper devices.
- Finding which device backs a mount point.

### Mounts

- Listing mounted filesystems.
- Inspecting one mount point.
- Mounting a prepared filesystem.
- Unmounting a filesystem.
- Diagnosing a busy mount.
- Recognizing bind mounts.
- Recognizing read-only mounts.
- Inspecting mount options.
- Understanding that the same filesystem may appear at several mount points.

## Likely commands

- `df`
- `du`
- `lsblk`
- `findmnt`
- `mount`

Supporting commands:

- `blkid`
- `lsof`
- `umount`

---

# 20. Manage Filesystems, Partitions, and Logical Storage

## Goal

Teach basic storage administration using disposable virtual disks and loop devices.

All destructive operations must be confined to dedicated lab devices.

## Topics

### Filesystems and loop devices

- Creating a file-backed virtual disk.
- Attaching a loop device.
- Creating a filesystem.
- Mounting the new filesystem.
- Writing and reading data from it.
- Unmounting and detaching it.
- Inspecting filesystem metadata.
- Running a safe filesystem check.
- Recognizing filesystem-specific tools.
- Comparing basic ext4 and XFS operational differences.

### Persistent mounts

- Reading `/etc/fstab`.
- Using filesystem UUIDs.
- Adding a persistent mount.
- Testing configuration with `mount -a`.
- Diagnosing an invalid mount option.
- Recognizing boot risks from invalid `fstab` entries.
- Using bind mounts.
- Using read-only and `noexec` mount options.

### Partitions and swap

- Inspecting a partition table.
- Creating a partition on a disposable disk.
- Informing the kernel of partition changes.
- Formatting a partition.
- Creating swap space.
- Enabling swap.
- Disabling swap.
- Inspecting active swap.

### LVM

- Creating a physical volume.
- Creating a volume group.
- Creating a logical volume.
- Formatting and mounting it.
- Extending a logical volume.
- Growing the filesystem.
- Inspecting the PV, VG, and LV relationship.
- Removing the temporary LVM stack safely.

### Additional storage topics

These can appear as harder reps, optional modules inside the path, or later extension paths:

- Software RAID with `mdadm`.
- Repairing a damaged filesystem.
- OverlayFS lower, upper, work, and merged directories.
- LUKS encryption and unlocking.
- ZFS pools, datasets, snapshots, and clones.
- Recognizing device-mapper layering.
- Understanding how filesystems, LVM, encryption, and RAID stack.

## Likely commands

- `losetup`
- `mkfs`
- `fdisk` or `parted`
- `swapon` and `swapoff`
- LVM commands such as `pvcreate`, `vgcreate`, and `lvcreate`

Because this path contains several tool families, it may be split during implementation into:

1. Filesystems and persistent mounts.
2. Partitions, swap, and LVM.

---

# 21. Inspect the Linux Host and Kernel Interfaces

## Goal

Teach the student where Linux exposes host, process, device, and kernel state.

## Topics

### Host identity and state

- Reading the kernel version.
- Reading machine architecture.
- Reading distribution identity.
- Reading hostname.
- Reading uptime.
- Reading system time and time zone.
- Recognizing boot time.

### Process information through `/proc`

- Inspecting `/proc/<pid>/status`.
- Reading a process command line.
- Reading a process environment.
- Inspecting its current working directory.
- Inspecting its executable link.
- Inspecting its open file descriptors.
- Inspecting its mounts.
- Reading process memory summaries.
- Connecting ordinary commands such as `ps` and `lsof` to `/proc`.

### System information through `/proc`

- Inspecting CPU information.
- Inspecting memory information.
- Inspecting system mounts.
- Inspecting load and uptime.
- Inspecting networking summaries.
- Understanding that many system tools read from procfs.

### Devices and sysfs

- Exploring `/sys`.
- Finding a network interface in sysfs.
- Finding a block device in sysfs.
- Reading device attributes.
- Understanding `/dev` as the device-node view.
- Connecting `/sys`, `/dev`, and `lsblk`.
- Recognizing major and minor device numbers.

### Kernel parameters and modules

- Reading a kernel parameter.
- Temporarily changing a safe sysctl.
- Making a sysctl persistent.
- Understanding `/proc/sys`.
- Reading kernel messages.
- Listing loaded modules.
- Inspecting module information.
- Loading a safe module.
- Unloading a safe module.
- Recognizing built-in functionality that is not represented by a loaded module.

## Likely commands

- `uname`
- `sysctl`
- `dmesg`
- `lsmod`
- `modinfo`

Most inspection should also directly use `/proc`, `/sys`, and `/dev`.

---

# 22. Work with Isolation and Resource Controls

## Goal

Introduce Linux namespaces, cgroups, capabilities, and seccomp as the mechanisms underlying process isolation and containers.

## Topics

### Namespaces

- Inspecting the namespaces of a process.
- Comparing namespace identifiers between processes.
- Recognizing PID, mount, network, user, IPC, UTS, cgroup, and time namespaces.
- Running a command in a new namespace.
- Entering an existing namespace.
- Observing PID differences inside and outside a namespace.
- Observing a namespace-specific hostname.
- Observing mount isolation.
- Observing network isolation.
- Understanding namespace membership as a process property.

### cgroup v2

- Finding the cgroup v2 mount.
- Inspecting the cgroup hierarchy.
- Finding the cgroup of a process.
- Creating a child cgroup.
- Moving a process into a cgroup.
- Applying a memory limit.
- Observing a memory limit.
- Applying a CPU control.
- Reading resource usage.
- Recognizing systemd-created cgroups.
- Connecting service units to cgroups.
- Cleaning up an experimental cgroup.

### Capabilities and seccomp

- Inspecting process capabilities.
- Comparing root with a process that has a limited capability set.
- Recognizing a capability-related failure.
- Inspecting seccomp state through `/proc`.
- Recognizing a syscall blocked by seccomp.
- Understanding at a practical level that namespaces do not themselves limit resource consumption.
- Understanding at a practical level that cgroups do not themselves isolate filesystem or network views.
- Understanding that containers combine several independent kernel mechanisms.

## Likely commands

- `lsns`
- `unshare`
- `nsenter`
- `systemd-run`
- Direct cgroup filesystem operations

Supporting tools:

- `capsh`
- `prlimit`
- `/proc`

---

# 23. Observe Performance and Triage a Host

## Goal

Teach a repeatable first-pass workflow for investigating CPU, memory, disk, process, and system-call activity.

## Topics

### Load and CPU

- Reading uptime and load average.
- Understanding load as runnable and uninterruptible work rather than CPU percentage.
- Inspecting overall CPU usage.
- Inspecting per-CPU usage.
- Recognizing user, system, idle, I/O wait, and steal time.
- Finding CPU-heavy processes.
- Observing process scheduling over time.

### Memory

- Reading total, used, free, and available memory.
- Understanding page cache at an operational level.
- Inspecting swap usage.
- Reading `/proc/meminfo`.
- Finding memory-heavy processes.
- Recognizing an OOM kill in logs.
- Inspecting process memory mappings at an introductory level.
- Reading `/proc/<pid>/smaps` where useful.

### System activity

- Reading `vmstat`.
- Interpreting runnable and blocked tasks.
- Interpreting paging and swapping.
- Interpreting CPU-state columns.
- Recognizing sustained pressure rather than one instantaneous sample.
- Inspecting Pressure Stall Information.

### Disk I/O

- Reading device throughput.
- Reading IOPS.
- Reading utilization.
- Recognizing queueing and latency indicators.
- Connecting a busy device to a process.
- Distinguishing filesystem capacity from storage performance.

### Per-process observation

- Observing per-process CPU.
- Observing per-process memory.
- Observing per-process disk I/O.
- Tracing file-related system calls.
- Tracing network-related system calls.
- Tracing process creation.
- Finding a failed syscall.
- Attaching a tracer to a running process.

### Extended observation topics

- Basic `perf` recording and reporting.
- Core dumps.
- cgroup-level resource observation.
- Slab and page-cache inspection.
- Introductory eBPF-based observation tools where available.

### Final triage circuit

The student should repeatedly answer questions such as:

- Is the host CPU-bound?
- Is it short on memory?
- Is it swapping?
- Is storage saturated?
- Which process is responsible?
- Is the process blocked on a syscall?
- Did the kernel or OOM killer terminate something?
- Is the problem global or limited to one cgroup?

## Likely commands

- `uptime`
- `top`
- `vmstat`
- `iostat`
- `pidstat`

Supporting tools:

- `free`
- `mpstat`
- `strace`
- `perf`
- `/proc`
- PSI files

---

# 24. Operate Remote Systems with SSH

## Goal

Teach routine and secure remote Linux operations.

## Topics

### Connections and commands

- Connecting interactively.
- Running one remote command.
- Running a remote pipeline safely.
- Observing remote exit status.
- Distinguishing local shell expansion from remote shell expansion.
- Selecting a username.
- Selecting a port.
- Handling connection failures.

### Host identity

- Understanding host-key verification.
- Inspecting a host fingerprint.
- Reading `known_hosts`.
- Handling a changed host key safely.
- Removing one obsolete host-key entry.
- Avoiding indiscriminate deletion of host-key records.

### Authentication

- Generating an SSH key.
- Protecting private-key permissions.
- Installing a public key.
- Reading `authorized_keys`.
- Selecting an identity explicitly.
- Using an SSH agent.
- Adding and listing agent identities.
- Recognizing an authentication-method failure.

### Configuration

- Creating a host alias.
- Configuring hostname, user, port, and identity.
- Inspecting effective SSH configuration.
- Using per-host options.
- Avoiding unnecessarily global settings.

### File transfer

- Copying a file to a remote host.
- Copying a remote file locally.
- Copying a directory recursively.
- Using SFTP interactively.
- Synchronizing a directory with rsync.
- Understanding rsync trailing-slash semantics.
- Performing an rsync dry run.
- Deleting remote files only when explicitly intended.

### Bastions and tunnels

- Connecting through a jump host.
- Configuring `ProxyJump`.
- Creating local port forwarding.
- Creating remote port forwarding.
- Creating a SOCKS proxy with dynamic forwarding.
- Finding which side of the SSH connection listens.
- Verifying forwarded connectivity.
- Stopping a tunnel.

### Security extensions

- Restricting `authorized_keys`.
- Understanding `PermitOpen`.
- Understanding `PermitListen`.
- Restricting shell access.
- Diagnosing agent-forwarding risks.
- Operating through a hardened bastion.

## Likely commands

- `ssh`
- `ssh-keygen`
- `ssh-add`
- `scp`
- `rsync`

Supporting tools:

- `sftp`
- `ssh-keyscan`
- `ssh-keygen -R`

---

# 25. Write Small Shell Scripts

## Goal

Teach enough shell scripting to automate small operational tasks without turning the path into a general programming course.

## Topics

- Creating a script file.
- Adding a shebang.
- Running a script through the interpreter.
- Making a script executable.
- Running it through its pathname.
- Understanding why the current directory is not normally searched automatically.
- Accepting positional arguments.
- Reading `$0`, `$1`, and `$#`.
- Iterating over `"$@"`.
- Returning a deliberate exit status.
- Propagating command failure.
- Testing file existence.
- Testing directory existence.
- Comparing strings.
- Comparing integers.
- Writing an `if` statement.
- Writing `if`/`else`.
- Combining test conditions.
- Writing a `for` loop.
- Writing a `while` loop.
- Reading lines from input.
- Defining a function.
- Passing function arguments.
- Using local variables.
- Returning status from a function.
- Calling earlier shell tools from a script.
- Quoting variable expansions.
- Producing useful diagnostics on stderr.
- Creating a small reusable operational script.

## Likely shell constructs

- Shebang
- Positional arguments
- `test` or `[ ]`
- `if`
- `for`
- Functions

Commands introduced earlier should be reused inside scripts.

---

# 26. Build Safer Command-Line Automation

## Goal

Teach defensive practices for temporary files, cleanup, pipelines, structured data, and repeatable automation.

## Topics

### Temporary state and cleanup

- Creating a unique temporary file.
- Creating a temporary directory.
- Avoiding predictable temporary filenames.
- Cleaning up after success.
- Cleaning up after failure.
- Installing a trap.
- Cleaning up after interruption.
- Preserving the original exit status during cleanup.

### Safer shell behavior

- Using `set -u`.
- Using `pipefail`.
- Understanding the limitations and surprises of `set -e`.
- Checking failures explicitly where necessary.
- Quoting expansions.
- Iterating over `"$@"`.
- Reading input without destroying backslashes.
- Handling filenames with whitespace.
- Handling filenames with newlines using null delimiters.
- Avoiding parsing human-formatted output when a machine-readable form exists.
- Writing idempotent operations.
- Performing dry runs before destructive changes.
- Logging what an automation changed.
- Sending diagnostics to stderr.

### Structured data

- Pretty-printing JSON.
- Selecting one JSON field.
- Iterating over JSON arrays.
- Filtering JSON objects.
- Producing raw text output.
- Constructing a small JSON object.
- Reading a known YAML value where `yq` is available.
- Editing structured data with an appropriate tool rather than regular-expression replacement.
- Combining `curl` and `jq`.
- Extracting values for use in later commands.

## Likely commands and features

- `mktemp`
- `trap`
- `set`
- `jq`
- `yq` where explicitly available

This path should prefer short, operational scripts rather than large programs.

---

# 27. Troubleshoot Commands, Permissions, Services, and Networking

## Goal

Reuse earlier skills in short incident-style drills without introducing much new syntax.

These exercises should require diagnosis, but remain faster and more guided than the platform's larger challenge format.

## Scenarios and topics

### “Command not found”

- Typographical error.
- Missing package.
- Incorrect `PATH`.
- Executable file lacking execute permission.
- Script invoked without a pathname.
- Command shadowed by an alias, function, or earlier `PATH` entry.
- Shell command hashing after an executable moves.

### “Permission denied”

- Wrong owner.
- Wrong group.
- Missing file permission.
- Missing directory traversal permission.
- Unwritable parent directory.
- Read-only mount.
- `noexec` mount.
- ACL denial.
- AppArmor or SELinux denial.
- Service running under a different user than expected.

### Failed service

- Invalid configuration.
- Missing file.
- Incorrect ownership.
- Port already in use.
- Missing dependency.
- Wrong environment.
- Incorrect working directory.
- Startup-order problem.
- Restart required after package or configuration changes.
- Service starts manually but fails under systemd.

### Port and socket failures

- Process already owns the port.
- Process listens only on loopback.
- Process listens on the wrong interface.
- Service uses IPv6 while the client assumes IPv4.
- Socket-activated listener owns the port.
- Stale process survives a failed restart.
- Firewall blocks access.
- Application accepts TCP but returns an application-level error.

### DNS and connectivity failures

- Wrong hostname.
- Missing `/etc/hosts` entry.
- Broken resolver configuration.
- DNS works but routing fails.
- Host is reachable but the port is closed.
- TCP connects but TLS validation fails.
- HTTP redirects unexpectedly.
- Ping is blocked even though the service works.
- Default route points to the wrong gateway.

## Rep design

Each rep should provide enough context to keep the task fast, but not directly prescribe the command sequence.

For example:

> The application should be reachable on port 8080 from the client host, but the connection fails. Restore access.

This is better than:

> Run `ss`, then change the bind address, then restart the service.

---

# 28. Troubleshoot Storage, Resources, and Host State

## Goal

Build a compact and repeatable Linux host-troubleshooting workflow.

## Scenarios and topics

### Full filesystem

- One unexpectedly large file.
- One unexpectedly large directory.
- Log growth.
- Deleted-but-open file.
- Inode exhaustion.
- Data written under a mount point while the filesystem was unmounted.
- Wrong filesystem inspected with `du`.
- Application cache growth.

### Broken mounts

- Invalid `/etc/fstab` syntax.
- Wrong UUID.
- Missing device.
- Unsupported filesystem type.
- Busy mount.
- Incorrect mount options.
- Read-only mount.
- Filesystem requiring repair.
- LVM logical volume not active.
- Encrypted device not unlocked.

### CPU pressure

- One CPU-heavy process.
- Many small competing processes.
- Runaway loop.
- Process with unexpectedly high system time.
- Load caused by blocked tasks rather than CPU saturation.
- One hot CPU while others remain mostly idle.

### Memory pressure

- Process consuming excessive memory.
- Swap activity.
- Memory constrained by a cgroup.
- OOM kill.
- Cache mistaken for unavailable memory.
- Process leaking memory.
- Large shared or mapped regions.

### Disk I/O pressure

- One process producing excessive writes.
- Saturated block device.
- High latency despite modest throughput.
- Log or database workload filling a queue.
- Activity occurring on a different device than expected.
- Filesystem full versus storage device slow.

### Process and kernel state

- Process ignoring SIGTERM.
- Zombie process.
- Process stuck in uninterruptible sleep.
- Exhausted file-descriptor limit.
- Missing kernel module.
- Incorrect sysctl.
- Clock or time-zone problem.
- Core dump after an application crash.
- Kernel or OOM messages explaining an incident.

### Final host-triage circuit

The student should learn to move through a compact workflow:

1. Establish the scope of the failure.
2. Check recent logs.
3. Check processes and services.
4. Check CPU and load.
5. Check memory and swap.
6. Check storage capacity and inodes.
7. Check disk activity.
8. Check sockets and network state.
9. Inspect kernel evidence.
10. Apply a small fix and verify the result.

---

# Curriculum implementation guidance

## Recommended release order

The curriculum should be released incrementally rather than authored in full before students can use it.

### Release 1: Basic terminal fluency

Paths:

1. Meet the Linux Shell
2. Control the Shell and Ask for Help
3. Navigate the Filesystem
4. List and Create Files
5. Copy, Move, Remove, and Link Files
6. Read and Inspect File Contents

This is enough to provide a meaningful beginner experience.

### Release 2: Command-line composition

Paths:

7. Redirect Streams and Build Pipelines
8. Search Files and Text
9. Transform Text and Records
10. Control the Shell Environment
11. Work with Processes and Jobs

At this point, students can perform meaningful command-line work rather than execute isolated commands.

### Release 3: Junior Linux operator

Paths:

12. Work with Users, Groups, and Permissions
13. Use Advanced Linux Access Controls
14. Install and Inspect Software Packages
15. Operate Services, Logs, and Scheduled Jobs
16. Inspect Network Configuration and Connectivity
17. Work with Sockets and Network Services

This is the first strong DevOps- and SRE-oriented milestone.

### Release 4: Host operator branches

Paths:

18. Route, Filter, and Isolate Network Traffic
19. Inspect Disk Usage, Devices, and Mounts
20. Manage Filesystems, Partitions, and Logical Storage
21. Inspect the Linux Host and Kernel Interfaces
22. Work with Isolation and Resource Controls
23. Observe Performance and Triage a Host

### Release 5: Remote work, automation, and synthesis

Paths:

24. Operate Remote Systems with SSH
25. Write Small Shell Scripts
26. Build Safer Command-Line Automation
27. Troubleshoot Commands, Permissions, Services, and Networking
28. Troubleshoot Storage, Resources, and Host State

---

# Curriculum rules specific to Shell Gym

## 1. Prefer observable outcomes over exact commands

A rep should normally verify the effect the student produced.

Good outcomes:

- The shell entered the requested directory.
- A file exists at the requested location.
- The file contains the requested data.
- A process is running.
- A process is no longer running.
- A service listens on the requested port.
- The port is no longer occupied.
- A route makes the destination reachable.
- A variable appears in the environment of a launched command.
- A file has the intended ownership and permissions.

Weaker tasks:

- Run `chmod`.
- Run `grep`.
- Run `systemctl restart`.
- Run `ip route add`.

The named command may fail, target the wrong object, or leave the intended state unchanged.

Use `wait_exec` mainly for operations that intentionally leave no persistent state, such as:

- Listing files.
- Printing file contents.
- Making a one-off HTTP request.
- Querying DNS.
- Inspecting a manual page.
- Running a diagnostic command whose use is itself the skill.

Even then, keep the regular expression permissive enough to accept valid command variants.

---

## 2. Train operations, not command names

A path should be framed around an operational ability.

Prefer:

- Navigate an unfamiliar directory tree.
- Find information in logs.
- Identify and stop a process.
- Determine which process owns a port.
- Repair file access.
- Inspect why a service failed.
- Find where disk space went.

Avoid paths framed only as:

- Learn `cd`.
- Learn `grep`.
- Learn `ps`.
- Learn `chmod`.
- Learn `systemctl`.

The commands remain important, but they should appear as tools for accomplishing practical work.

---

## 3. Repeat the same skill through varied situations

A single successful rep does not build muscle memory.

At the same time, repetition should not feel like the exact same exercise with a different filename.

Weak repetition:

- Find `ERROR` in file A.
- Find `ERROR` in file B.
- Find `ERROR` in file C.

Stronger repetition:

- Find a request ID in an application log.
- Exclude health-check requests.
- Show context around a crash.
- Search several rotated files.
- Search case-insensitively.
- Count matching failures.
- Find which configuration file contains a value.
- Feed matching records into another command.
- Save the result without losing diagnostics.

The underlying skill repeats, but its application changes.

---

## 4. Introduce few commands and revisit many

Each path should normally introduce no more than 3–5 commands.

Prefer 2–3 when possible.

A path introducing `find` should naturally reuse:

- `cd`
- `ls`
- quoting
- globs
- `grep`
- `rm`
- redirects

Later paths should keep earlier skills alive without announcing a formal review session.

For example:

- A service-management rep may require navigating to a configuration directory.
- A networking rep may require filtering `ss` output.
- A storage rep may require finding and sorting large files.
- A troubleshooting rep may require checking permissions before restarting a service.

---

## 5. Use 15–30 reps as a default, not an absolute rule

A normal path should contain 15–30 reps.

Use more reps when:

- Each action takes only a few seconds.
- The skill requires substantial repetition.
- The scene can vary naturally.
- The student benefits from building speed and confidence.

Examples:

- Filesystem navigation.
- Listing files.
- Basic copying and moving.
- Simple redirection.
- Quoting.

Use fewer reps when:

- Each rep requires a larger setup.
- The operation takes longer to verify.
- The student must reason across several system components.
- The path works with networking topologies, filesystems, services, or performance incidents.

Examples:

- Network namespaces.
- LVM.
- Firewall rules.
- SSH tunnels.
- Host-troubleshooting scenarios.

---

## 6. Keep most reps under one minute

Most individual reps should be solvable in less than a minute once the student understands the operation.

Longer thinking belongs in the platform's challenge format.

A Shell Gym rep may require choosing among a few commands, but it should not normally require:

- Reading several pages of documentation.
- Reverse engineering an unfamiliar application.
- Discovering a hidden multi-step root cause.
- Designing a complete solution.
- Writing a large script.
- Making architectural decisions.

Troubleshooting paths may be somewhat harder, but should still provide constrained, small scenes.

---

## 7. Use realistic but tiny scenes

Good scenes:

- A miniature project directory.
- Three rotated log files.
- A broken service.
- Two processes competing for one port.
- A temporary loop-backed filesystem.
- A pair of network namespaces.
- A home directory with incorrect ownership.
- A small HTTP API.
- A process with an open deleted file.
- A small JSON response.
- A temporary bastion and internal host.

Avoid enormous simulated production systems.

The student should recognize the operational pattern without spending ten minutes reading background material.

---

## 8. Randomize enough to keep repetition honest

Use variables to vary:

- Directory names.
- File names.
- Tokens.
- Usernames.
- Group names.
- Process names.
- Ports.
- Hostnames.
- IP addresses.
- Log values.
- Search patterns.
- File sizes.
- Permissions.

Randomization should prevent memorizing one literal solution, but it must not make instructions difficult to read.

Prefer meaningful randomized choices over opaque random strings where possible.

Good:

```yaml
DIRNAME:
  pick: [archive, cache, reports, uploads]
```

Use random tokens where uniqueness is part of the task:

```yaml
TOKEN:
  shell: "head -c4 /dev/urandom | od -An -tx1 | tr -d ' \n'"
```

---

## 9. Do not put the exact solution in the task text

The task should describe the desired result, not paste the command.

Good:

> Move all `.log` files from `${SOURCE}` into `${ARCHIVE}`.

Weak:

> Run `mv ${SOURCE}/*.log ${ARCHIVE}`.

Hints may identify:

- The type of command.
- A relevant concept.
- Why the previous attempt failed.
- A useful manual-page section.
- The current observed state.

Hints should not simply reveal the complete final command unless the learning experience explicitly supports a last-resort solution reveal outside the ordinary hint flow.

---

## 10. Use helpful tips for interaction that is hard to verify

Not every useful terminal behavior can be observed reliably from outside the shell.

Do not omit these skills. Add concise tips directly to relevant task pages.

Examples:

- Tab completion.
- Up-arrow history.
- `Ctrl-R` reverse history search.
- `Ctrl-A` and `Ctrl-E`.
- `Ctrl-U` and `Ctrl-K`.
- `Ctrl-C`.
- `Ctrl-Z`.
- Searching and quitting `less`.
- Searching manual pages.
- Selecting and pasting text.
- Opening another terminal when comparing two processes or network endpoints.

The task itself can verify the resulting system state while the tip teaches the interaction technique.

---

## 11. Prefer state checks, but use command checks when appropriate

Prefer:

- `wait_cwd`
- `wait_file`
- `wait_file_gone`
- `wait_file_contains`
- `wait_proc`
- `wait_proc_gone`
- `wait_port`
- `wait_port_free`
- Direct inspection in check scripts

Use `wait_exec` when:

- The command leaves no durable state.
- Running the diagnostic command is itself the intended skill.
- Several valid command forms can be accepted with a permissive pattern.

Remember:

- `wait_exec` proves that a command was executed.
- It does not prove that the command succeeded.
- It only observes commands executed after the unit's latest activation.
- It should not require one exact spelling when several correct forms exist.

Where possible, combine command observation with state verification.

---

## 12. Guard negative checks with a baseline

A task checking that something disappeared must first establish that it existed.

Bad:

```bash
wait_file_gone "$TARGET"
```

This may pass immediately before the scene is ready.

Better:

```bash
wait_file --timeout 15 "$TARGET" || exit 1
wait_file_gone "$TARGET"
```

The same principle applies to:

- Processes that must stop.
- Ports that must become free.
- Files that must be removed.
- Configuration lines that must disappear.
- Routes or firewall rules that must be removed.

---

## 13. Keep init scripts idempotent

Activation may retry init after failure or reset.

Init scripts should:

- Safely recreate the scene.
- Remove stale state from previous attempts.
- Avoid failing when an object already exists.
- Use dedicated names.
- Avoid changing unrelated host state.
- Clean up old processes.
- Clean up old mounts.
- Clean up temporary network namespaces.
- Set deterministic ownership and permissions.
- Avoid relying on an earlier unit unless declared through `needs:`.

Init runs as root, so student-owned files must be explicitly assigned to `$GYM_USER`.

---

## 14. Use dependencies only when state genuinely carries forward

Students may attempt units out of order.

Never assume an earlier unit was completed unless the current unit declares it in `needs:`.

Use dependencies when the current unit intentionally builds on earlier state, such as:

- A directory tree created earlier.
- A user or group configured earlier.
- A filesystem prepared earlier.
- A service installed earlier.
- A network topology constructed earlier.
- A randomized value that must remain consistent.

Keep dependency chains short, preferably fewer than five units.

For standalone practice, rebuild the scene independently instead of creating unnecessary chains.

---

## 15. Use task dependencies for meaningful multi-stage reps

Most units should teach one small action.

A multi-task unit is appropriate when the actions form one coherent operational sequence, such as:

1. Find the process holding a port.
2. Stop that process.
3. Verify that the port is free.

Use task-level `needs:` so the UI presents the sequence clearly.

Do not create long procedural units merely to simulate a challenge. If a unit becomes a substantial investigation, it probably belongs in the challenge format instead.

---

## 16. Distinguish edge and level tasks carefully

Use edge tasks for accomplishments:

- The student entered a directory.
- The student ran a diagnostic command.
- The student created a file.
- The student stopped a process.
- The student exported a variable and launched a command.

Use level tasks for conditions that must still be true when the unit completes:

- A service is running.
- A port is listening.
- A file has a particular mode.
- A route currently exists.
- A mount is active.
- A firewall rule currently allows traffic.

Remember:

- Level tasks may become false again.
- Unit completion occurs when all edge tasks have completed and all level tasks are simultaneously satisfied.
- Edge tasks may not depend on level tasks.

---

## 17. Use the student's home directory for persistent cross-unit state

Choose scene location according to lifetime.

Use `/tmp/gym` or another volatile location for:

- Independent units.
- Disposable files.
- Short-lived sockets.
- Temporary processes.
- Loop-device scenes that are recreated on activation.

Use the student's home directory when:

- A later unit depends on the state.
- The scene should survive a reboot.
- The task concerns shell configuration.
- The task concerns SSH configuration.
- The task concerns user-owned scripts.
- The task concerns persistent project files.

Use system paths only when the learning objective is specifically about those paths.

---

## 18. Make destructive operations safe by construction

The student should practice potentially dangerous commands, but only inside contained scenes.

Examples:

- `rm -r` should target a dedicated generated tree.
- Partitioning should use disposable virtual disks.
- Filesystem creation should use loop devices or dedicated attached disks.
- Firewall exercises should use isolated namespaces or an explicit recovery mechanism.
- Route changes should avoid breaking Shell Gym control connectivity.
- Permission exercises should avoid locking the daemon out of its own state.
- Process exercises should use clearly named disposable processes.
- Package removal should use non-essential packages.
- Service exercises should use dedicated lab services.

The scene should make the correct operation safe and an overly broad operation detectable where possible.

---

## 19. Acceptance-test every unit through a real shell

Every task requires a valid hidden `solve:` block.

Solve scripts must:

- Contain commands that can be typed one line at a time.
- Avoid heredocs.
- Avoid multiline constructs.
- Avoid line continuations.
- Avoid relying on shell state that was never established.
- Respect declared unit dependencies.
- Exercise the same observable path available to the student.

Run:

```sh
shellgym validate --path my-path
shellgym solve --path my-path
```

For focused iteration:

```sh
shellgym solve --path my-path --unit module/unit
```

If the unit has dependencies, solve the dependency chain first or run the whole path in order.

---

## 20. Keep theory short and operational

Each path may have a concise introduction explaining:

- What the student will be able to do.
- Why the skill matters in Linux operations.
- The minimum mental model needed to avoid cargo-culting commands.

Do not turn Shell Gym units into textbook chapters.

Long-form tutorials should explain topics such as:

- How the shell parses and expands a command.
- Why directory permissions behave differently from file permissions.
- Why load average includes uninterruptible tasks.
- How ext4 allocates inodes.
- How systemd builds dependency transactions.
- How conntrack interacts with NAT.
- How namespaces and cgroups are implemented.
- How filesystems and block layers are stacked.
- How DNS resolution works internally.

Shell Gym should make the associated operations automatic.

---

# Path authoring checklist

Before accepting a learning path, confirm:

- The path represents one coherent operational skill area.
- It normally takes 15–30 minutes.
- It contains roughly 15–30 reps, adjusted for rep complexity.
- It introduces no more than 3–5 new commands.
- It reuses previously learned commands.
- Most reps complete in less than a minute.
- Repetition is varied rather than cosmetic.
- The exact solution is not present in the task statement.
- Helpful terminal interaction tips are included where relevant.
- Outcomes are verified through state whenever possible.
- `wait_exec` patterns accept legitimate alternatives.
- Negative checks establish a baseline first.
- Init scripts are idempotent.
- Root-created student files are correctly owned.
- Destructive operations are isolated.
- Units do not depend on undeclared earlier state.
- Dependency chains are short.
- Variables are used to keep repetition honest.
- Every task has a valid `solve:` block.
- `shellgym validate` passes.
- `shellgym solve` passes.
- The final reps combine the new skill with earlier skills.
- The path feels like terminal practice rather than a quiz, tutorial, or large challenge.
