# Absolute bearings

Every shell sits in exactly one directory at all times - its working
directory. `pwd` prints it, and `cd` moves it. That is the whole toolkit
of this path; the skill is using it without thinking.

An **absolute path** starts with `/` and names a location from the root
of the filesystem, so it works no matter where you currently are. This
module is all about jumping straight to places by their full address.
