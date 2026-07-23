# Redirection and pipes

Commands read from an input stream and write to output streams. The shell
lets you re-plumb those streams: send output into files (`>`, `>>`) or feed
one command's output into another command (`|`). This composability is the
core idea of the Unix command line, and this module trains it.
