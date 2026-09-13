# Take Back Control

Sooner or later every terminal user meets a command that does not give
the prompt back. It runs too long, or it waits for input that never
comes, or you simply changed your mind. Closing the terminal window
with the "stuck" command is one way out but it's definitely not the most elegant one.

The shell has keys for exactly this situation. **Ctrl-C** interrupts
the foreground command. **Ctrl-Z** suspends it and keeps it around.
The `fg` command brings a suspended command back. In this module you will break out
of stuck commands on purpose until the key strokes come without thinking.
