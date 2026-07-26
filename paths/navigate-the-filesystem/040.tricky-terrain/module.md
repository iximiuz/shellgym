# Tricky terrain

Real directories are not always called `src` and `docs`. Names contain
spaces, parentheses, dollar signs - characters the shell itself treats
specially. The fix is always the same: quote the name, escape the odd
character, or let Tab completion do the quoting for you.

This module runs you through the terrain where naive typing fails, and
ends with a circuit that puts the whole path together.
