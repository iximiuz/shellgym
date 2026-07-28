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

Its unique role is to build command-line fluency through repeated practice while giving the student an initial practical understanding of how Linux behaves.

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
- Explicitly state its prerequisites.
- Explicitly state the knowledge it assumes.
- Avoid task descriptions that require knowledge not listed in the assumed-knowledge section.

The first part of the curriculum should be mostly linear. Once the student has basic terminal, filesystem, process, permission, package, service, and logging skills, the roadmap may branch into:

- Networking.
- Storage.
- Kernel and host internals.
- Performance.
- Remote operations.
- Automation.
- Troubleshooting.

The learning paths can later be connected in a visual roadmap with prerequisites and optional branches.

---

# General assumptions

The curriculum may assume that the student:

- Knows what a file is.
- Knows what a directory or folder is.
- Understands that files are stored on disks or other storage devices.
- Can use a keyboard and enter text.
- Can open a terminal application when shown where it is.

The curriculum must not initially assume that the student:

- Understands the Linux single-directory filesystem hierarchy.
- Knows what `/`, `/home`, `/tmp`, `/etc`, or `/var` mean.
- Understands absolute and relative paths.
- Has used a command-line shell.
- Has written software.
- Has used Git.
- Understands source code, repositories, deployments, services, packages, processes, ports, or system administration.
- Recognizes development or operations terminology unless an earlier path has introduced it.

Knowledge should be introduced before it appears in task descriptions. The terminology and context used in tasks should become more technical as the student progresses through the curriculum.

---

# Scenario continuity

The learning paths should form several connected sequences rather than presenting unrelated task collections.

The continuity should make later tasks feel like a progression from earlier work. It must not create undeclared state dependencies between units or require the student to remember irrelevant fictional details.

No path or task should require prior familiarity with a fictional application merely because that application appeared earlier in the curriculum.

Whenever a path uses an application:

- Reintroduce the application briefly in the path introduction.
- Explain only the parts relevant to the current path.
- Make each task understandable without knowledge of earlier application tasks.
- Keep application-specific behaviour simple and explicit.
- Ensure that the difficulty comes from the Linux operation being practised, not from understanding the application.

## Paths 1–6: Basic Linux terminal and filesystem work

The student begins using a Linux shell for ordinary file and system operations.

Tasks should use contexts that apply to both workstations and servers:

- Inspecting the current user, host, and working directory.
- Moving through the user's home directory.
- Working with temporary files under `/tmp`.
- Inspecting common system directories at an introductory level.
- Creating work directories under the home directory.
- Backing up files.
- Renaming and organizing configuration or data files.
- Reading text files.
- Inspecting generated output.

Do not use development or system-administration plots that require concepts not yet introduced.

Avoid terms such as:

- Repository.
- Checkout.
- Deployment.
- Build artifact.
- Service configuration.
- Production logs.
- Source tree.

These concepts have not yet been introduced.

## Paths 7–10: Command-line work with a simple inventory tool

After learning basic filesystem operations, the student begins working with a small command-line inventory application.

The working application name is **Stockroom**.

Stockroom is a simple tool that reads and writes plain-text inventory records. Example records may describe:

- Asset IDs.
- Item names.
- Locations.
- Quantities.
- Owners.
- Status values.
- Update timestamps.

The command-line tool should remain deliberately simple. A task may state, for example:

> `stockroom list` prints one inventory record per line.

The student should not need to infer command syntax, data formats, or application behaviour from earlier tasks.

Tasks may involve:

- Saving Stockroom output.
- Filtering inventory records.
- Finding files containing an asset ID.
- Sorting records by quantity.
- Setting an output directory through an environment variable.

Use ordinary paths such as:

- `~/stockroom`
- `~/stockroom/data`
- `~/stockroom/reports`
- `/tmp/stockroom`
- `/opt/stockroom`
- `/var/lib/stockroom`
- `/etc/stockroom`

Do not add `gym`, `shellgym`, `gymtrack`, or similar prefixes to ordinary project, user, process, service, or filesystem names unless the task is specifically about the Shell Gym software itself.

## Paths 11–17: Operating Stockroom Server

Starting with Path 11, introduce **Stockroom Server** as a standalone application context.

Stockroom Server may be described as a background service that exposes inventory data over HTTP and periodically imports records from files.

It can be related to the earlier Stockroom command-line tool, but no path may assume that the learner completed or remembers those earlier tasks.

Every path using Stockroom Server should briefly explain:

- What the service does.
- Which file, process, socket, user, or endpoint matters for the current tasks.
- Any application command or data format needed to complete the tasks.

Tasks may involve:

- Starting and stopping the server process.
- Managing its files and permissions.
- Installing required packages.
- Operating a systemd service.
- Inspecting logs.
- Checking network configuration.
- Sending HTTP requests.
- Inspecting sockets.

At this stage, task descriptions may use common development and operations terms that have already been introduced.

## Paths 18–23: Administering a Linux application host

The student takes responsibility for a Linux host running Stockroom Server or another similarly simple prepared service.

Each path must still reintroduce the relevant service or process independently.

Tasks may involve:

- Routing and firewall rules.
- Network namespaces.
- Disks and mounts.
- Filesystems and LVM.
- Kernel interfaces.
- Namespaces and cgroups.
- Host performance investigation.

These paths may use realistic system-administration terminology because the earlier paths have established the required foundation.

## Paths 24–26: Remote operation and automation

The student operates one or more Linux hosts remotely.

The tasks may use Stockroom Server, but they must not require prior Stockroom knowledge. A task should state all relevant hostnames, paths, commands, endpoints, and expected behaviour.

Tasks may involve:

- SSH access.
- Secure authentication.
- Remote file transfer.
- Bastions and tunnels.
- Small maintenance scripts.
- Safe repeatable automation.
- Structured JSON and YAML data.

## Paths 27–28: Operational troubleshooting

The student diagnoses constrained failures involving commands, files, services, storage, and network components used in earlier paths.

A troubleshooting task may use Stockroom Server or another prepared application. Each task must provide enough application context to understand the expected state without requiring knowledge of previous paths.

The scenarios should remain small enough for Shell Gym. They should not require a long investigation or unfamiliar domain knowledge.

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

## Prerequisites

None.

## Assumed knowledge

The student:

- Knows what files and directories are.
- Can use a keyboard and enter text.
- Does not need prior terminal or Linux experience.
- Is not expected to understand the Linux filesystem hierarchy.

## Scenario context

The student opens a Linux terminal for the first time. The tasks introduce command execution through simple actions that do not require development or system-administration knowledge.

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

Not every command listed here needs to count as a separately taught command. Some are simple tools for creating observable reps.

## Helpful task-page tips

Some terminal interaction is difficult to verify directly. Relevant tips should appear alongside tasks rather than being omitted from the curriculum:

- Use the Up and Down arrow keys to revisit command history.
- Use Left and Right arrows to edit a command.
- Use `Ctrl-A` and `Ctrl-E` to jump to the beginning and end of the line.
- Use `Ctrl-U` and `Ctrl-K` to remove parts of the current line.
- Use `Ctrl-L` to clear the visible terminal.
- Use Tab completion instead of typing long paths manually.

These tips can be introduced gradually and repeated where they become useful.

## Rep design

Use simple actions with immediately visible results:

- Print a short message.
- Print two words as separate arguments.
- Print one argument containing spaces.
- Run a command that succeeds.
- Run a command that fails.
- Run a fallback command after a failure.
- Print the current username.
- Print the host name.
- Print the current directory.
- Run a command that succeeds without printing output.

Do not introduce project, source-code, service, package, or system-administration terminology in this path.

---

# 2. Control the Shell and Ask for Help

## Prerequisites

- Meet the Linux Shell.

## Assumed knowledge

The student:

- Can enter commands.
- Recognizes commands, arguments, and options.
- Understands successful and unsuccessful command completion at a basic level.
- Can read the previous exit status.
- Is not expected to know how Linux documentation is organized.

## Scenario context

The student continues learning to operate the terminal and needs to recover from mistakes, stop commands, and find command documentation without leaving the shell.

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
- Recovering from a command that appears not to finish.
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

## Rep design

Use ordinary terminal situations:

- Check whether a command is available.
- Determine whether a name refers to a builtin or executable.
- Open the documentation for a command used in the previous path.
- Find the option that changes a command's output.
- Stop a command that continues running.
- Suspend and resume a command.
- Search a manual page for a word.
- Exit a manual page and return to the prompt.

Do not require the student to understand services, source code, or package management.

---

# 3. Navigate the Filesystem

## Prerequisites

- Meet the Linux Shell.
- Control the Shell and Ask for Help.

## Assumed knowledge

The student:

- Knows what files and directories are.
- Can enter commands and read their output.
- Can request command help.
- Does not yet know the Linux single-directory hierarchy.
- Does not yet know absolute or relative path syntax.

## Scenario context

The student begins accessing files and system locations from the terminal. The tasks use the home directory, temporary directories, and common Linux directories while introducing the Linux filesystem hierarchy.

## Goal

Build strong familiarity with moving through the Linux filesystem and introduce the Linux single-directory hierarchy, including the purpose of important locations such as `/`, `/home`, and `/tmp`.

This path may contain more than 30 reps because most navigation operations are quick.

## Topics

### Linux filesystem hierarchy

- Understanding that Linux exposes one directory hierarchy beginning at `/`.
- Recognizing `/` as the root directory.
- Understanding that mounted storage appears at locations inside the same hierarchy.
- Recognizing `/home` as the usual parent directory for user home directories.
- Recognizing the student's home directory.
- Recognizing `/tmp` as a location for temporary files.
- Recognizing `/etc`, `/var`, `/usr`, and `/opt` at a superficial level without requiring administration knowledge.
- Distinguishing the root directory `/` from the root user's home directory `/root`.

### Filesystem navigation

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
- Locating a target file by inspecting several directories.
- Returning to a known location after checking another directory.

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

Use varied tasks that apply to both workstation and server environments:

- Jump to the user's home directory.
- Jump to a prepared work directory under the home directory.
- Move from a data directory to its backup directory.
- Jump to `/tmp`.
- Inspect the location of a prepared temporary file.
- Jump from `/var` to `/var/tmp`.
- Visit `/etc` and return to the previous directory.
- Navigate directories with spaces in their names.
- Return to the previous directory after checking another location.
- Start from an unknown directory and return home.
- Compare the current directory in two separate terminal sessions.
- Follow a provided relative path to a target file.

Do not use tasks about checked-out repositories, source trees, deployments, or service configuration.

---

# 4. List and Create Files

## Prerequisites

- Navigate the Filesystem.

## Assumed knowledge

The student:

- Understands the Linux single-directory hierarchy at an introductory level.
- Can navigate with absolute and relative paths.
- Knows the purpose of the home directory and `/tmp`.
- Understands quoting for paths containing spaces at a basic level.
- Is not expected to understand Unix permissions beyond recognizing that permission fields exist.

## Scenario context

The student uses the terminal to inspect directories and create files and directory structures under the home directory and `/tmp`.

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

Use small tasks such as:

- List files in a work directory.
- Find the newest backup file.
- Find the largest file in a directory listing.
- Show hidden files in the home directory.
- Create a directory for backups.
- Create nested directories for data and reports.
- Create an empty status file.
- Create a filename that contains spaces.
- Create a hidden configuration directory.
- List `/tmp` and the home directory without changing the working directory.
- Preview which files a glob matches.
- Work with a filename beginning with `-`.

Do not require project, programming, version-control, package-management, or system-service knowledge.

---

# 5. Copy, Move, Remove, and Link Files

## Prerequisites

- List and Create Files.

## Assumed knowledge

The student:

- Can navigate the filesystem.
- Can list files and inspect basic listing metadata.
- Can create files and directories.
- Understands basic glob expansion and quoting.
- Is not expected to understand inodes before this path introduces hard links.

## Scenario context

The student organizes files under the home directory, creates backups, moves generated data, removes temporary files, and learns how links differ from copies.

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

Use file-management tasks that apply to both workstations and servers:

- Copy a configuration file into a backup directory.
- Copy several data files into an archive.
- Rename a generated file to a clearer name.
- Move temporary output from `/tmp` into the home directory.
- Move matching data files into a dated directory.
- Remove a duplicate backup.
- Remove an empty temporary directory.
- Remove a generated directory tree.
- Preview matching files before removing them.
- Create a convenient symbolic link to a frequently used directory.
- Repair a broken symbolic link.
- Compare a copied file with a hard-linked file.

Include multiple exercises where destination semantics matter. Many students can run `cp` but still hesitate over whether a destination path means:

- A new filename.
- An existing directory.
- A directory to be copied.
- The contents of a directory.

This distinction deserves deliberate repetition.

---

# 6. Read and Inspect File Contents

## Prerequisites

- Copy, Move, Remove, and Link Files.

## Assumed knowledge

The student:

- Can navigate, list, create, copy, move, and remove files.
- Understands basic paths, quoting, and globbing.
- Knows that some files contain text and others may contain binary data.
- Is not expected to understand application logs or structured operational records yet.

## Scenario context

The student reads configuration fragments, generated reports, status files, and other text files from the terminal. The path introduces methods appropriate for small, large, and growing files.

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
- Following a growing text file.
- Recognizing newly appended records.
- Stopping a follow operation.
- Counting lines.
- Counting words.
- Counting bytes.
- Numbering lines.
- Comparing two versions of a file.
- Recognizing added, removed, and changed lines.
- Inspecting an unfamiliar collection of files.
- Choosing between full output, partial output, and interactive browsing.
- Inspecting a file's type before attempting to read it.
- Avoiding printing a very large or binary file directly into the terminal.
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

## Rep design

Use tasks such as:

- Read a short configuration fragment.
- Combine several report parts in order.
- Check the beginning of an exported data file.
- Check the latest entries in a growing status file.
- Browse a large list without printing all of it.
- Search for a value inside a large text file.
- Count records in a file.
- Compare an original configuration with a backup.
- Check a file's type before opening it.
- Avoid printing a prepared binary file to the terminal.

The final reps may introduce Stockroom as a simple command-line tool that produces plain-text inventory records. The task must explain the command and output without assuming prior familiarity.

---

# 7. Redirect Streams and Build Pipelines

## Prerequisites

- Read and Inspect File Contents.

## Assumed knowledge

The student:

- Can inspect files and command output.
- Can use basic shell operators such as `;`, `&&`, and `||`.
- Understands that commands may succeed or fail.
- Does not yet need a formal understanding of file descriptors.
- Does not need prior familiarity with Stockroom.

## Scenario context

This path uses a small command-line tool named Stockroom.

Stockroom prints simple inventory records. Each task states the relevant command and what it prints. For example, a task may explain that `stockroom list` prints one record per line.

The student needs to save, combine, and filter this output.

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

- Save the inventory list printed by `stockroom list`.
- Append a later inventory snapshot without deleting the earlier one.
- Run a prepared Stockroom command that emits diagnostics and save only its errors.
- Count the records produced by a stated command.
- Display and save a filtered list.
- Feed a prepared records file into a command through stdin.
- Keep useful output and diagnostics in separate files.

Each task must explain any Stockroom command it uses. The student should not need to remember Stockroom syntax from another rep.

Weak:

- “Run a command using `2>`.”
- “Use a pipe.”

---

# 8. Search Files and Text

## Prerequisites

- Redirect Streams and Build Pipelines.

## Assumed knowledge

The student:

- Can navigate and manipulate files.
- Can inspect file contents.
- Can use redirection and pipelines.
- Does not need prior regular-expression knowledge.
- Does not need prior familiarity with Stockroom.

## Scenario context

This path uses prepared Stockroom inventory files.

A Stockroom record is a plain-text line containing fields such as an asset ID, item name, location, quantity, or status. Each task explains the relevant record format and target value.

The student needs to find specific files and records.

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

- Find a stated asset ID in a prepared inventory file.
- Exclude records whose status field is `retired`.
- Show context around a malformed record.
- Search several archived inventory files.
- Search item names case-insensitively.
- Count records matching a specified location.
- Find which configuration file contains a stated value.
- Locate recently modified report files.
- Find unexpectedly executable files.
- Remove only matching temporary files.
- Find files and then inspect their contents.
- Find files and move them to an archive directory.

Each task must define the relevant record format or file layout. No Stockroom knowledge should carry over implicitly.

---

# 9. Transform Text and Records

## Prerequisites

- Search Files and Text.

## Assumed knowledge

The student:

- Can search files and text.
- Can build pipelines and redirect output.
- Understands plain-text records and delimited fields at an introductory level.
- Has not previously used `sed` or `awk`.
- Does not need prior familiarity with Stockroom.

## Scenario context

This path uses plain-text inventory records produced for Stockroom.

Each task defines the field separator and meaning of the relevant fields. The student needs to reorder, extract, normalize, count, and summarize the records.

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

A path may technically contain more than five command names because several are small filters. The author should still minimize simultaneous novelty and reuse commands across many reps.

## Rep design

Use inventory-data tasks:

- Sort records by item name.
- Sort records numerically by quantity.
- Count repeated location values.
- Extract the asset ID field from a stated format.
- Convert status values to a consistent case.
- Remove empty or comment lines.
- Replace an obsolete location code.
- Filter records whose quantity exceeds a stated threshold.
- Calculate a total quantity.
- Preserve a header while sorting the remaining records.
- Produce a summary report from several pipeline stages.

Every task must include enough format information to be completed independently.

---

# 10. Control the Shell Environment

## Prerequisites

- Transform Text and Records.

## Assumed knowledge

The student:

- Can compose commands with pipelines and redirects.
- Understands basic quoting for filenames.
- Has seen `$?` and `$HOME`.
- Does not yet understand shell-variable scope, environment inheritance, or command lookup in detail.
- Does not need prior familiarity with Stockroom.

## Scenario context

This path uses a standalone Stockroom command-line tool.

The tool can read values such as the data directory and output format from environment variables. Each task identifies the relevant variable and expected behaviour.

The student needs to control these values through shell variables and environment variables.

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
- A Stockroom data directory supplied through a stated environment variable.
- A one-command output-format override.
- A command installed in a prepared executable directory.
- A temporary `PATH` extension.

Each task must explain any Stockroom-specific variable or behaviour it uses.

---

# 11. Work with Processes and Jobs

## Prerequisites

- Control the Shell Environment.

## Assumed knowledge

The student:

- Can run foreground commands.
- Has interrupted and suspended commands.
- Can use redirection.
- Understands shell variables and command lookup.
- Does not need prior knowledge of process trees, signals, file descriptors, or scheduling.
- Does not need prior familiarity with Stockroom.

## Scenario context

This path introduces **Stockroom Server** as a standalone application.

Stockroom Server is a prepared program that can run in the foreground or background and serve inventory data. The tasks provide the exact start command or describe how the prepared process was launched.

No knowledge of the Stockroom command-line tool is required.

The student needs to inspect server processes, keep jobs running, stop incorrect processes, and determine which files or sockets they use.

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

## Rep design

Use small local-process scenarios:

- Find a running Stockroom Server process.
- Distinguish two prepared processes with similar names.
- Inspect the command line of a process.
- Stop a process normally.
- Stop a process that ignores SIGTERM.
- Start a prepared report generator in the background.
- Bring a background job to the foreground.
- Keep a command running after the shell exits.
- Find the process using a prepared file.
- Inspect a process's open file descriptors.
- Adjust the priority of a disposable process.

Each task must state what the relevant process does and how it can be identified.

---

# 12. Work with Users, Groups, and Permissions

## Prerequisites

- Work with Processes and Jobs.

## Assumed knowledge

The student:

- Can inspect processes and their owners.
- Can read long-format file listings.
- Recognizes owner, group, and permission fields but has not yet studied their complete meaning.
- Understands that several user accounts can exist on one Linux host.
- Does not yet understand Unix permission evaluation.
- Does not need prior familiarity with Stockroom Server.

## Scenario context

This path uses a prepared installation of Stockroom Server.

The service reads inventory data from `/var/lib/stockroom` and writes reports to a separate directory. Each task explains which user or group needs access and what operation must succeed.

The student needs to configure who can read, modify, or run the relevant files and directories.

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

Use realistic access tasks:

- Give the stated Stockroom service user read access to an inventory file.
- Allow a prepared operator group to update a data directory.
- Prevent other users from reading a configuration file.
- Repair incorrect ownership.
- Allow directory traversal without allowing directory listing.
- Make a maintenance script executable without changing unrelated permissions.
- Create files with private default permissions.
- Create files with group-shared default permissions.

Each task must identify the relevant user, group, file, and required operation. No application knowledge should be assumed.

Directory permissions need dedicated reps. They are commonly misunderstood even by students who can decode `rwx` on regular files.

---

# 13. Use Advanced Linux Access Controls

## Prerequisites

- Work with Users, Groups, and Permissions.

## Assumed knowledge

The student:

- Understands UIDs, GIDs, ownership, mode bits, and `umask`.
- Can use `sudo` in a controlled environment.
- Understands file and directory permission differences.
- Does not yet understand ACLs, capabilities, or mandatory access controls.
- Does not need prior familiarity with Stockroom Server.

## Scenario context

This path uses prepared Stockroom files and helper programs.

Ordinary owner, group, and mode permissions are insufficient for several stated access requirements. Each task describes the required access or operation directly.

The student uses additional Linux access-control mechanisms and diagnoses which mechanism blocks an operation.

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

## Rep design

Use contained access scenarios:

- Configure a shared import directory.
- Ensure new files inherit the intended group.
- Grant one additional user access without changing the owning group.
- Configure a default ACL for new report files.
- Find files with special permission bits.
- Grant a prepared helper executable one narrow capability.
- Remove an unnecessary capability.
- Determine whether an access failure comes from mode bits, an ACL, AppArmor, or SELinux.

Every task must explicitly state the expected operation and the relevant file or process.

---

# 14. Install and Inspect Software Packages

## Prerequisites

- Use Advanced Linux Access Controls.

## Assumed knowledge

The student:

- Can navigate system directories.
- Can inspect files, permissions, and processes.
- Can use `sudo`.
- Understands commands and executables but does not yet understand how Linux distributions package software.
- Is not expected to know distribution families before this path.
- Does not need prior familiarity with Stockroom Server.

## Scenario context

This path uses a prepared Stockroom Server installation that depends on ordinary command-line utilities.

Each task states which command or capability is needed. The student identifies the Linux distribution and installs, removes, and inspects packages safely.

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

## Rep design

Use concrete package-management tasks:

- Identify the distribution.
- Find which package provides a stated command.
- Inspect a package before installing it.
- Install a non-essential utility required by a prepared maintenance command.
- Find the files installed by a package.
- Find which package owns an executable.
- Compare the installed and available versions.
- Perform a dry run before removal.
- Remove a disposable package.
- Distinguish package removal from configuration purging.

The application context should remain incidental. The task must state the missing command or required package property directly.

---

# 15. Operate Services, Logs, and Scheduled Jobs

## Prerequisites

- Install and Inspect Software Packages.

## Assumed knowledge

The student:

- Can inspect and control processes.
- Can manage files, permissions, and packages.
- Understands that a program can run in the background.
- Has not yet operated systemd services, the journal, cron, or timers.
- Does not need prior familiarity with Stockroom Server.

## Scenario context

This path uses a prepared systemd service named `stockroom.service`.

The service exposes inventory data and writes diagnostic messages to the journal. A separate scheduled job imports records periodically.

Each task explains the relevant service, file, or scheduled operation independently.

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

## Rep design

Use a dedicated Stockroom service and timer:

- Check whether the service is running.
- Start and stop it.
- Reload a stated configuration change where supported.
- Enable it without starting it.
- Start it without enabling it.
- Find the failed unit.
- Read the unit file.
- Inspect a drop-in override.
- Find the error that caused startup failure.
- Follow new service log entries.
- Inspect logs from the current boot.
- Find the next scheduled import.
- Run the scheduled service manually.
- Diagnose a cron job with a missing environment variable.

Every task must state the relevant unit name and expected behaviour.

---

# 16. Inspect Network Configuration and Connectivity

## Prerequisites

- Operate Services, Logs, and Scheduled Jobs.

## Assumed knowledge

The student:

- Can operate a service and inspect its logs.
- Understands that hosts communicate over networks.
- Recognizes hostnames at a basic level.
- Is not expected to understand interfaces, routes, ARP, DNS internals, TCP, or ports before they are introduced.
- Does not need prior familiarity with Stockroom Server.

## Scenario context

This path uses a prepared Stockroom Server host that needs to communicate with a separate inventory database or import host.

Each task states the relevant hostname, address, or expected connection.

The student inspects the host's network configuration and diagnoses connectivity in a consistent order.

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
- Changing an interface MTU in a safe scenario.
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

## Rep design

Use small topologies and clear objectives:

- Identify the interface used for outbound traffic.
- Find the host's IPv4 and IPv6 addresses.
- Determine which route reaches a stated destination.
- Inspect a failed neighbor entry.
- Resolve a stated hostname through the system resolver.
- Compare system resolution with a direct DNS query.
- Add a temporary `/etc/hosts` override.
- Determine whether a failure is caused by DNS or routing.
- Verify that a host is reachable even when ping is blocked.
- Test whether the expected TCP port is reachable.

Every task must state the expected destination or service independently.

---

# 17. Work with Sockets and Network Services

## Prerequisites

- Inspect Network Configuration and Connectivity.

## Assumed knowledge

The student:

- Understands interfaces, addresses, routes, hostname resolution, and basic connectivity checks.
- Can operate a systemd service.
- Understands that services listen on network addresses and ports at an introductory level.
- Has not yet used HTTP clients, manual TCP connections, socket inspection, or packet capture.
- Does not need prior familiarity with Stockroom Server.

## Scenario context

This path uses a prepared Stockroom Server instance.

The service exposes a small HTTP API. Each task states the endpoint, expected response, listener address, or traffic pattern being investigated.

The student connects to the service, inspects its sockets, and captures selected traffic.

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

## Rep design

Use Stockroom Server and dedicated disposable network services:

- Request a stated status endpoint.
- Inspect headers without downloading the response body.
- Follow a redirect.
- Send a small JSON request whose format is provided.
- Diagnose a certificate-name mismatch.
- Connect to a prepared TCP service manually.
- Create a temporary listener and send data to it.
- Find which process owns a stated port.
- Diagnose a service listening only on loopback.
- Compare IPv4 and IPv6 listeners.
- Capture one HTTP request.
- Filter a packet capture by host and port.

Every task must provide the relevant endpoint, request format, or expected listener state.

---

# 18. Route, Filter, and Isolate Network Traffic

## Prerequisites

- Work with Sockets and Network Services.

## Assumed knowledge

The student:

- Understands interfaces, addresses, routes, sockets, and ports.
- Can generate and inspect network traffic.
- Can use `sudo`.
- Has not yet configured packet forwarding, nftables, NAT, network namespaces, or virtual Ethernet links.
- Does not need prior familiarity with Stockroom Server.

## Scenario context

This path uses isolated networks containing a prepared HTTP service such as Stockroom Server.

Each task describes the topology, addresses, ports, and expected connectivity. The application itself is not part of the problem.

All routing, firewall, NAT, and namespace changes occur in isolated topologies that cannot interrupt Shell Gym connectivity.

## Goal

Introduce host routing, firewalling, NAT, and Linux virtual networking through small isolated scenarios.

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

## Rep design

Use contained network topologies:

- Add a route between two prepared subnets.
- Remove an incorrect route.
- Enable forwarding for a namespace router.
- Add an allow rule for one stated service port.
- Add and verify a drop rule for another port.
- Inspect rule counters after generating traffic.
- Configure masquerading for a private namespace.
- Redirect one local port to another.
- Create and address a veth pair.
- Connect two namespaces through a bridge.
- Repair an interface or route removed from a prepared topology.
- Delete all temporary network objects cleanly.

The tasks must specify the topology directly and must not rely on application knowledge.

---

# 19. Inspect Disk Usage, Devices, and Mounts

## Prerequisites

- Operate Services, Logs, and Scheduled Jobs.
- Recommended: Search Files and Text.

## Assumed knowledge

The student:

- Can navigate, search, and inspect files.
- Can inspect processes and open files.
- Understands in ordinary terms that storage devices hold files.
- Does not yet understand block devices, filesystems, mount points, or inode exhaustion.
- Does not need prior familiarity with Stockroom Server.

## Scenario context

This path uses prepared application data under `/var/lib/stockroom` and several disposable filesystems.

Each task states which directory, mount, or device is relevant.

The student determines where space is used and how the host's devices and mounted filesystems relate to directories.

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

## Rep design

Use dedicated application data and disposable devices:

- Find which filesystem contains a stated data directory.
- Compare filesystem free space with directory usage.
- Find the largest subdirectory.
- Find a large generated file.
- Identify a deleted file that remains open.
- Diagnose inode exhaustion caused by many small files.
- Identify the block device backing a mount.
- Read its filesystem type and UUID.
- Mount a prepared filesystem.
- Diagnose a busy unmount.
- Inspect a read-only or bind mount.
- Determine why data appears at more than one mount point.

Every task must identify the relevant path or device directly.

---

# 20. Manage Filesystems, Partitions, and Logical Storage

## Prerequisites

- Inspect Disk Usage, Devices, and Mounts.

## Assumed knowledge

The student:

- Understands block devices, partitions, filesystems, mount points, and UUIDs at an introductory level.
- Can mount and unmount prepared filesystems.
- Can use `sudo`.
- Has not yet created partitions, filesystems, swap, or LVM structures.

## Scenario context

The host needs additional storage for application data.

The student performs all destructive work on disposable virtual disks and loop devices. Each task identifies the exact device that is safe to modify.

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

## Rep design

Use only explicitly identified disposable devices:

- Create a file-backed disk.
- Attach it to a loop device.
- Create and mount an ext4 filesystem.
- Add a persistent mount by UUID.
- Detect and repair an invalid `fstab` entry.
- Create a partition on an attached disposable disk.
- Create and activate swap.
- Create a small LVM volume for application data.
- Extend the logical volume and filesystem.
- Inspect each storage layer.
- Remove the temporary storage stack safely.

---

# 21. Inspect the Linux Host and Kernel Interfaces

## Prerequisites

- Inspect Disk Usage, Devices, and Mounts.
- Work with Processes and Jobs.
- Operate Services, Logs, and Scheduled Jobs.

## Assumed knowledge

The student:

- Can inspect processes, files, devices, mounts, services, and logs.
- Understands that the Linux kernel manages processes, memory, devices, and filesystems.
- Has used `/proc/<pid>/fd` but has not systematically studied `/proc`, `/sys`, `/dev`, sysctls, or modules.
- Does not need prior familiarity with Stockroom Server.

## Scenario context

This path uses prepared processes, devices, and a Stockroom Server instance where useful.

Each task identifies the relevant PID, interface, block device, or kernel setting. No application behaviour needs to be inferred.

The student needs more detailed host information than ordinary commands provide.

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

## Rep design

Use direct host-inspection tasks:

- Find the current kernel version and architecture.
- Determine the host's boot time.
- Read the environment of a specified process.
- Find its current directory and executable.
- Compare its open descriptors through `/proc` and `lsof`.
- Read CPU and memory information directly from procfs.
- Find the sysfs entry for a network interface.
- Connect a block device in `/dev` with its sysfs entry.
- Read and temporarily change a safe sysctl.
- Inspect a loaded module.
- Distinguish a loaded module from built-in kernel functionality.

---

# 22. Work with Isolation and Resource Controls

## Prerequisites

- Inspect the Linux Host and Kernel Interfaces.
- Use Advanced Linux Access Controls.

## Assumed knowledge

The student:

- Understands processes, `/proc`, capabilities, services, and basic kernel interfaces.
- Understands that one host can run several processes under different users.
- Has not yet used namespaces, cgroup v2 controls, or seccomp directly.
- Does not need prior familiarity with Stockroom Server.

## Scenario context

This path uses several prepared worker processes.

Some may be described as Stockroom Server workers, but each task states what the process does and which isolation or resource property matters.

The student inspects and changes their isolation and resource-control settings.

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

## Rep design

Use isolated disposable processes:

- Compare the namespaces of two specified processes.
- Run a shell with a different hostname namespace.
- Observe different PID values inside and outside a PID namespace.
- Enter a prepared network namespace.
- Find the cgroup of a stated service.
- Create a child cgroup for a test process.
- Apply and observe a memory limit.
- Apply and observe a CPU limit.
- Compare process capabilities.
- Inspect seccomp state in `/proc`.
- Trigger a prepared capability or seccomp failure.
- Clean up the experimental process and cgroup.

---

# 23. Observe Performance and Triage a Host

## Prerequisites

- Inspect the Linux Host and Kernel Interfaces.
- Recommended: Work with Isolation and Resource Controls.
- Recommended: Inspect Disk Usage, Devices, and Mounts.

## Assumed knowledge

The student:

- Can inspect processes, services, logs, `/proc`, cgroups, disks, and mounts.
- Understands CPU, memory, storage, and processes in ordinary technical terms.
- Has not yet learned a structured Linux performance-investigation workflow.
- Does not need prior familiarity with Stockroom Server.

## Scenario context

This path uses prepared processes that create controlled CPU, memory, and storage pressure.

A process may be described as a Stockroom worker, but each task states what it is expected to do and which symptom must be investigated.

The student collects evidence about CPU, memory, storage, processes, and system calls.

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

### Final triage workflow

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

## Rep design

Use controlled performance conditions:

- Identify a CPU-heavy process.
- Distinguish high load caused by CPU work from blocked tasks.
- Find a process consuming excessive memory.
- Identify swap activity.
- Find an OOM kill in the journal.
- Inspect memory pressure through PSI.
- Identify a busy block device.
- Connect storage activity to a process.
- Trace a failed file operation.
- Trace a failed network connection.
- Determine whether pressure is host-wide or limited to one cgroup.

Every task must explain the expected process role and symptom directly.

---

# 24. Operate Remote Systems with SSH

## Prerequisites

- Inspect Network Configuration and Connectivity.
- Work with Users, Groups, and Permissions.
- Control the Shell Environment.

## Assumed knowledge

The student:

- Can operate a local Linux shell.
- Understands users, permissions, hostnames, addresses, ports, and basic network connectivity.
- Can manipulate files and inspect command exit status.
- Has not yet used SSH keys, host verification, remote commands, file transfer, bastions, or tunnels.
- Does not need prior familiarity with Stockroom Server.

## Scenario context

This path uses several prepared Linux hosts.

One host may run Stockroom Server, but every task states the remote hostname, user, relevant path, command, or service endpoint.

The student connects securely, runs commands, transfers files, and reaches internal services through an SSH bastion.

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

## Rep design

Use a small remote environment:

- Connect to a specified host.
- Run one remote status command.
- Preserve the intended local or remote variable expansion.
- Inspect and verify the host key.
- Repair one changed `known_hosts` entry.
- Generate and install a key.
- Configure a host alias.
- Copy a configuration file.
- Synchronize a report directory with a dry run first.
- Connect through a bastion.
- Forward a local port to a stated remote service.
- Expose a prepared local service through remote forwarding.
- Create and stop a SOCKS proxy.

Every task must be self-contained with respect to the remote application or service.

---

# 25. Write Small Shell Scripts

## Prerequisites

- Operate Remote Systems with SSH.
- Control the Shell Environment.
- Search Files and Text.
- Transform Text and Records.

## Assumed knowledge

The student:

- Can compose shell commands, pipelines, conditions with `&&` and `||`, variables, and quoted expansions.
- Can operate on files and remote hosts.
- Has not yet written complete shell scripts with arguments, tests, loops, or functions.
- Does not need prior familiarity with Stockroom Server.

## Scenario context

This path uses small maintenance tasks involving files, prepared services, and remote hosts.

Some scripts may process Stockroom inventory files or query Stockroom Server, but each task defines the input, output, command, and data format required.

The student places repeated command sequences into reusable scripts.

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

## Rep design

Use small maintenance scripts:

- Print a status summary from stated input files.
- Accept a report filename as an argument.
- Verify that an input file exists.
- Return a non-zero status for invalid input.
- Process every filename in `"$@"`.
- Filter records using commands learned earlier.
- Copy selected files to an archive.
- Print diagnostics to stderr.
- Define a function that validates a directory.
- Run one stated command remotely for each prepared host.
- Create one complete script that remains short and directly testable.

Any application command or data format must be explained inside the task.

---

# 26. Build Safer Command-Line Automation

## Prerequisites

- Write Small Shell Scripts.

## Assumed knowledge

The student:

- Can write and run small shell scripts.
- Understands arguments, conditions, loops, functions, exit statuses, quoting, pipelines, and stderr.
- Can work with local and remote files.
- Has not yet systematically used temporary directories, traps, stricter shell modes, null delimiters, or structured-data tools.
- Does not need prior familiarity with Stockroom Server.

## Scenario context

This path uses prepared maintenance scripts that modify files or consume JSON and YAML data.

Some data may come from a Stockroom Server endpoint, but every task states the endpoint, fields, and expected output.

The student makes the scripts safer, repeatable, and suitable for structured input.

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

## Rep design

Use realistic automation requirements:

- Create a unique temporary directory.
- Remove it after success.
- Remove it after failure or interruption.
- Preserve the script's original exit status.
- Detect an unset required variable.
- Detect a failed pipeline stage.
- Process filenames safely with null delimiters.
- Add a dry-run mode.
- Make a directory-creation operation idempotent.
- Extract a stated endpoint from JSON.
- Filter a JSON array of hosts or assets.
- Update a prepared YAML value with `yq`.
- Combine `curl` and `jq` without parsing human-formatted output.

Every structured-data task must identify the relevant fields and expected result.

---

# 27. Troubleshoot Commands, Permissions, Services, and Networking

## Prerequisites

- Operate Services, Logs, and Scheduled Jobs.
- Work with Sockets and Network Services.
- Operate Remote Systems with SSH.
- Recommended: Build Safer Command-Line Automation.

## Assumed knowledge

The student:

- Can inspect commands, files, permissions, processes, services, logs, network configuration, sockets, and SSH connections.
- Is expected to choose among commands introduced in earlier paths.
- Is not expected to use unfamiliar debugging tools or perform long investigations.
- Does not need prior familiarity with Stockroom Server.

## Scenario context

This path uses short incidents involving a prepared service such as Stockroom Server.

Every task explains:

- What the service or command is expected to do.
- Which file, port, user, hostname, or endpoint matters.
- What symptom is currently visible.

The student does not need knowledge from earlier Stockroom tasks.

## Goal

Reuse earlier skills in short incident-style drills without introducing much new syntax.

These exercises should require diagnosis but remain faster and more guided than the platform's larger challenge format.

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

Each rep should provide enough context to keep the task fast but not directly prescribe the command sequence.

For example:

> The `stockroom.service` unit should serve HTTP on port 8080. A client on the same host cannot connect. Restore access.

This is better than:

> Run `ss`, then change the bind address, then restart the service.

A task using Stockroom Server must explain the expected listener, endpoint, file, or user directly.

Tasks should not require the learner to recall earlier fictional application behaviour.

---

# 28. Troubleshoot Storage, Resources, and Host State

## Prerequisites

- Inspect Disk Usage, Devices, and Mounts.
- Manage Filesystems, Partitions, and Logical Storage.
- Inspect the Linux Host and Kernel Interfaces.
- Observe Performance and Triage a Host.
- Recommended: Work with Isolation and Resource Controls.

## Assumed knowledge

The student:

- Can inspect filesystems, devices, mounts, processes, services, logs, kernel state, cgroups, CPU, memory, and disk activity.
- Is expected to apply workflows from earlier paths.
- Is not expected to investigate large distributed systems or unfamiliar applications.
- Does not need prior familiarity with Stockroom Server.

## Scenario context

This path uses short host incidents involving a prepared application service.

Every task states the expected service behaviour, relevant data path, process role, and visible symptom.

The application provides a realistic workload but should not add diagnostic complexity.

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

### Final host-triage workflow

The student should learn to use this sequence:

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

## Rep design

Use one or a few consistent host layouts where practical:

- Keep common service names, paths, users, groups, ports, and storage locations consistent.
- Explain all relevant application details in each task.
- Change the Linux failure being diagnosed rather than making the application itself difficult to understand.
- Provide enough context to keep the investigation within a few commands.
- Avoid hidden causes that require reading large amounts of unfamiliar configuration.
- Require verification after every repair.
- Prefer faults that correspond directly to the mechanisms taught in earlier paths.

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

This is enough to provide a meaningful beginner experience based on ordinary file and system operations that apply to both workstation and server environments.

### Release 2: Command-line composition

Paths:

7. Redirect Streams and Build Pipelines
8. Search Files and Text
9. Transform Text and Records
10. Control the Shell Environment
11. Work with Processes and Jobs

Paths 7–10 may use the simple Stockroom command-line tool.

Path 11 introduces Stockroom Server independently. The learner must not need knowledge of the earlier command-line tool.

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

Stockroom Server may provide a consistent application context, but each path and task must reintroduce the relevant application behaviour.

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

## 3. Use standard technical language

Task descriptions, hints, path introductions, and supporting text should use direct technical language.

Do not use metaphors, analogies, fictional physical spaces, or other figures of speech to explain Linux operations.

Prefer:

- Jump to the directory containing the data files.
- Find the configuration file.
- Inspect the process.
- Stop the service.
- Follow the log.
- Pipe the command's output into another command.
- The directory is not accessible because execute permission is missing.

Avoid:

- Enter the room containing the files.
- Explore a maze of directories.
- Follow breadcrumbs to the file.
- Fight the process.
- The service is asleep.
- The pipeline is a conveyor belt.
- Permissions are keys to locked rooms.

Widely accepted technical slang is acceptable when it makes the writing less formal and remains unambiguous. Examples include:

- Jump to a directory.
- Kill a process.
- Tail a log.
- Pipe output into another command.
- The host is swapping.
- The port is already taken.
- The process is stuck.

---

## 4. Use tasks applicable to workstation and server environments

Early tasks should not focus on consumer-oriented personal-computer activities such as:

- Organizing a photo collection.
- Browsing a downloads folder.
- Managing music or video files.
- Cleaning up desktop shortcuts.

Prefer tasks that make sense in both workstation and server contexts:

- Working with files under the home directory.
- Creating a work directory.
- Backing up a file.
- Moving generated output.
- Cleaning up temporary files.
- Inspecting `/tmp`.
- Reading a configuration fragment.
- Comparing a file with its backup.
- Inspecting common directories such as `/etc`, `/var`, and `/opt` at an appropriate level.
- Creating and organizing data, reports, configuration, or archive directories.

Later tasks should reflect realistic technical work:

- Filtering application output.
- Finding a value in a configuration file.
- Repairing permissions.
- Restarting a service.
- Finding which process owns a port.
- Checking disk usage.
- Connecting to a remote host.
- Diagnosing resource pressure.

---

## 5. Match task context to curriculum depth

The terminology used in a task must be appropriate for the path's position in the curriculum.

Early paths must not assume familiarity with:

- Git.
- Repositories.
- Checking out code.
- Source trees.
- Builds.
- Deployments.
- Services.
- Containers.
- Packages.
- Production systems.
- Incident response.

For example, a basic navigation task should say:

> Find the directory containing the prepared backup and jump to it.

It should not say:

> A repository has been checked out under `/srv`; jump to its source directory.

Technical terminology may be introduced after the curriculum has taught the required concepts.

When introducing a new term:

- Explain it briefly.
- Use it consistently.
- Do not require additional unlisted knowledge to understand the task.
- Add it to the path's assumed knowledge only after an earlier module has introduced it.

---

## 6. Maintain scenario continuity without creating application prerequisites

Paths should form a small number of connected sequences.

Use the curriculum-wide progression:

1. Basic Linux file and shell operations.
2. Command-line processing with a simple inventory tool.
3. Operating a standalone inventory service.
4. Administering the Linux host.
5. Operating remote hosts.
6. Automating and troubleshooting Linux operations.

Continuity should reduce unnecessary context, not create a prerequisite to understand a fictional application.

A path using Stockroom or Stockroom Server must reintroduce it briefly.

A task using Stockroom or Stockroom Server must be understandable without:

- Completing an earlier Stockroom task.
- Remembering an earlier Stockroom command.
- Knowing an earlier file format.
- Knowing an earlier endpoint.
- Remembering fictional business context.

The task should contain only enough background to explain:

- What the application component does.
- Which command, file, process, endpoint, or service matters.
- What state is expected.
- What result the student must produce.

Do not require the student to remember character names, company history, business goals, or unrelated plot details.

---

## 7. Keep application-specific complexity low

The complexity of a rep should come from the Linux capability being trained.

The application used in the task should be simple enough that the student does not need to investigate it.

For example, a task about redirection may say:

> `stockroom list` prints one inventory record per line. Save its output to `${REPORT}`.

A task about HTTP may say:

> Stockroom Server listens on port `${PORT}` and returns its status at `/health`. Request that endpoint.

A task about permissions may say:

> The `stockroom` user must read `${DATA_FILE}`. Change the file's ownership or permissions so that this succeeds.

Avoid tasks that require the student to:

- Discover undocumented application commands.
- Infer an unknown file format.
- Understand application-specific configuration syntax that has not been provided.
- Reverse engineer an unfamiliar API.
- Diagnose application logic before applying the Linux operation.

---

## 8. Allow learners to enter the curriculum at later paths

Some learners will skip early paths and start from Path 10, 11, or later.

Each path must therefore be understandable when opened directly from the catalog, subject only to its explicitly listed Linux prerequisites.

When a path introduces a new application or a new application mode:

- Describe it independently.
- Do not assume familiarity with an earlier mode.
- State the relevant paths, commands, files, endpoints, and users.
- Explain the expected state.

Stockroom Server may be an extension of the Stockroom command-line tool, but Path 11 and later paths must not require knowledge of the command-line tool.

The same principle applies to optional branches. A storage path must not assume that the learner remembers a networking-path scenario, and a networking path must not depend on a storage-path application setup unless declared explicitly.

---

## 9. Avoid unnecessary Shell Gym naming

The Shell Gym daemon runs on a fresh, disposable Linux machine whose purpose is to host the learning environment and allow the learner to modify or break it safely.

Ordinary example resources should use ordinary names.

Prefer:

- `/opt/stockroom`
- `/etc/stockroom`
- `/var/lib/stockroom`
- `/var/log/stockroom`
- `stockroom.service`
- `stockroom-worker`
- `operator`
- `reports`
- `archive`
- `imports`

Avoid:

- `/opt/gymtrack/stockroom`
- `/tmp/gym-project`
- `gym-user`
- `gym-service`
- `shellgym-app`
- `gymnet`
- `gymdisk`

Use `Shell Gym`, `shellgym`, `$GYM_USER`, and related names only when referring to the actual product, daemon, authoring interface, or runtime variable.

---

## 10. Repeat the same skill through varied situations

A single successful rep does not establish reliable command recall.

At the same time, repetition should not feel like the exact same exercise with a different filename.

Weak repetition:

- Find `ERROR` in file A.
- Find `ERROR` in file B.
- Find `ERROR` in file C.

Stronger repetition:

- Find an asset ID in an inventory file.
- Exclude retired records.
- Show context around a malformed entry.
- Search several archived files.
- Search case-insensitively.
- Count matching records.
- Find which configuration file contains a value.
- Pipe matching records into another command.
- Save the result without losing diagnostics.

The underlying skill repeats, but its application changes.

---

## 11. Introduce few commands and revisit many

Each path should normally introduce no more than 3–5 commands.

Prefer 2–3 when possible.

A path introducing `find` should naturally reuse:

- `cd`
- `ls`
- Quoting
- Globs
- `grep`
- `rm`
- Redirects

Later paths should continue requiring earlier skills without announcing a separate review section.

For example:

- A service-management rep may require navigating to a configuration directory.
- A networking rep may require filtering `ss` output.
- A storage rep may require finding and sorting large files.
- A troubleshooting rep may require checking permissions before restarting a service.

---

## 12. Use 15–30 reps as a default, not an absolute rule

A normal path should contain 15–30 reps.

Use more reps when:

- Each action takes only a few seconds.
- The skill requires substantial repetition.
- The scenario can vary naturally.
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

## 13. Keep most reps under one minute

Most individual reps should be solvable in less than a minute once the student understands the operation.

Longer investigation belongs in the platform's challenge format.

A Shell Gym rep may require choosing among a few commands, but it should not normally require:

- Reading several pages of documentation.
- Reverse engineering an unfamiliar application.
- Discovering a hidden multi-step root cause.
- Designing a complete solution.
- Writing a large script.
- Making architectural decisions.

Troubleshooting paths may be somewhat harder but should still provide constrained, small scenarios.

---

## 14. Use realistic but small scenarios

Good scenarios:

- A work directory under the user's home directory.
- A backup directory.
- A collection of data files.
- A configuration file and its backup.
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

Avoid large simulated production systems.

The student should recognize the operational pattern without spending several minutes reading background material.

---

## 15. Randomize enough to keep repetition honest

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
  pick: [archive, cache, reports, imports]
```

Use random tokens where uniqueness is part of the task:

```yaml
TOKEN:
  shell: "head -c4 /dev/urandom | od -An -tx1 | tr -d ' \n'"
```

Randomized values should not weaken scenario continuity. A Stockroom task should still use recognizable Stockroom paths and terminology even when filenames, ports, users, or tokens vary.

---

## 16. Do not put the exact solution in the task text

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

Providing an application command that is not itself being taught is acceptable.

For example, in a redirection task:

> `stockroom list` prints the current inventory. Save its output to `${REPORT}`.

The task may provide `stockroom list` because the learning objective is redirection, not discovering the application command.

---

## 17. Use helpful tips for interaction that is hard to verify

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

## 18. Prefer state checks, but use command checks when appropriate

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

## 19. Guard negative checks with a baseline

A task checking that something disappeared must first establish that it existed.

Bad:

```bash
wait_file_gone "$TARGET"
```

This may pass immediately before the scenario is ready.

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

## 20. Keep init scripts idempotent

Activation may retry init after failure or reset.

Init scripts should:

- Safely recreate the scenario.
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

## 21. Use dependencies only when state genuinely carries forward

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

For standalone practice, rebuild the scenario independently instead of creating unnecessary chains.

Curriculum narrative continuity does not itself justify a state dependency.

Two units may both refer to Stockroom while independently recreating the required state.

No unit should depend on a learner remembering an earlier application's command syntax or data format.

---

## 22. Use task dependencies for meaningful multi-stage reps

Most units should teach one small action.

A multi-task unit is appropriate when the actions form one coherent operational sequence, such as:

1. Find the process holding a port.
2. Stop that process.
3. Verify that the port is free.

Use task-level `needs:` so the UI presents the sequence clearly.

Do not create long procedural units merely to simulate a challenge. If a unit becomes a substantial investigation, it probably belongs in the challenge format instead.

---

## 23. Distinguish edge and level tasks carefully

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

## 24. Use the student's home directory for persistent cross-unit state

Choose scenario location according to lifetime.

Use `/tmp` or another volatile location for:

- Independent units.
- Disposable files.
- Short-lived sockets.
- Temporary processes.
- Loop-device scenarios that are recreated on activation.

Use the student's home directory when:

- A later unit depends on the state.
- The scenario should survive a reboot.
- The task concerns shell configuration.
- The task concerns SSH configuration.
- The task concerns user-owned scripts.
- The task concerns persistent work files.

Use system paths only when the learning objective is specifically about those paths.

Avoid creating a generic `/tmp/gym`, `/opt/gym`, or similarly named hierarchy merely because the task runs inside Shell Gym.

---

## 25. Make destructive operations safe by construction

The student should practice potentially dangerous commands, but only inside contained scenarios.

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

The scenario should make the correct operation safe and an overly broad operation detectable where possible.

Because the machine is disposable and dedicated to Shell Gym, realistic system paths and service names may be used when appropriate. Safety should come from controlled setup and verification rather than from adding `gym` to every resource name.

---

## 26. State prerequisites and assumed knowledge explicitly

Every learning path must define:

- Prerequisites.
- Assumed knowledge.
- Scenario context.
- What the student will learn.

The first module of every implemented learning path must present the same information to the student.

The first module should include:

### Prerequisites

List required earlier paths by name.

Use `None` for the first path or for a genuinely independent extension path.

Distinguish required prerequisites from recommended prerequisites.

### Assumed knowledge

State the concepts, commands, terminology, and general computer knowledge that task descriptions may rely on.

Also make important non-assumptions explicit where useful. For example:

- No prior Linux experience is required.
- Knowledge of files and directories is assumed.
- Knowledge of Git is not assumed.
- The Linux filesystem hierarchy will be introduced in this path.
- Basic systemd knowledge is assumed from an earlier path.
- No prior familiarity with Stockroom is required.

### What the student will learn

Provide a short, direct list of operational outcomes.

For example:

- Navigate with absolute and relative paths.
- Return to the home and previous directories.
- Work with paths containing spaces.
- Recognize the purpose of `/`, `/home`, and `/tmp`.

### Scenario context

Explain the continuing scenario in one or two short paragraphs.

The context must:

- Use terminology appropriate for the current curriculum depth.
- Avoid metaphors and analogies.
- Avoid unnecessary fictional details.
- Explain why the operations are useful.
- Remain understandable when the path is opened directly from the catalog.
- Reintroduce any fictional application or service used by the path.

Task authors must verify that every term used in a task is either:

- Common knowledge listed in the general assumptions.
- Introduced earlier in the same path.
- Included in the path's assumed knowledge.
- Explained directly in the task.

---

## 27. Make application-based tasks self-contained

Any task that uses Stockroom, Stockroom Server, or another fictional application must include the application information needed for that task.

Depending on the task, this may include:

- The command to run.
- What the command prints.
- The input file format.
- The relevant field separator.
- The configuration path.
- The service unit name.
- The service user.
- The listening address and port.
- The HTTP endpoint.
- The expected response.
- The process name.
- The location of data or logs.

Do not require the learner to retrieve these details from an earlier task merely to practise a Linux capability.

Application details may be repeated across tasks when repetition keeps those tasks self-contained.

The task should not include the target Linux solution, but it may include application-specific syntax that is outside the learning objective.

---

## 28. Acceptance-test every unit through a real shell

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

## 29. Keep theory short and operational

Each path should have a concise introduction explaining:

- Its prerequisites.
- The knowledge it assumes.
- What the student will be able to do.
- Why the skill matters in Linux operations.
- The minimum mental model needed to avoid using commands without understanding their effect.
- How the path fits into the current curriculum sequence.
- Any fictional application or service used by the path.

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

Shell Gym should make the associated operations familiar through repeated execution.

---

# Path authoring checklist

Before accepting a learning path, confirm:

- The path represents one coherent operational skill area.
- It normally takes 15–30 minutes.
- It contains roughly 15–30 reps, adjusted for rep complexity.
- It introduces no more than 3–5 new commands.
- It reuses previously learned commands.
- Its prerequisites are explicitly listed.
- Its assumed knowledge is explicitly listed.
- Its first module presents prerequisites, assumed knowledge, scenario context, and learning outcomes.
- Every term used in task descriptions is either assumed, previously introduced, or explained.
- The task context is appropriate for the path's position in the curriculum.
- Early paths do not assume development or system-administration knowledge.
- Early tasks apply to both workstation and server environments.
- Early tasks do not focus on consumer-oriented file-management activities.
- The Linux single-directory hierarchy is introduced before tasks rely on it.
- Scenario continuity with nearby paths is maintained where useful.
- Scenario continuity does not create undeclared state dependencies.
- Each path reintroduces any fictional application it uses.
- Each application-based task includes the application details needed to complete it.
- Application knowledge from earlier paths is not required.
- A later application mode, such as Stockroom Server, is introduced independently.
- The difficulty comes from the Linux capability rather than application discovery.
- Ordinary resources do not receive unnecessary `gym`-based names.
- Task descriptions use direct, standard technical language.
- Task descriptions do not use metaphors or unrelated real-world analogies.
- Tasks represent realistic computer or Linux operations.
- Most reps complete in less than a minute.
- Repetition is varied rather than cosmetic.
- The exact Linux solution is not present in the task statement.
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
- The path provides terminal practice rather than a quiz, textbook chapter, or large challenge.
