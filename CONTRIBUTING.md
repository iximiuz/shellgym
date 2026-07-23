# Contributing to Shell Gym

Contributions are welcome - bug fixes, new checks, docs improvements, and learning paths alike.
For anything non-trivial, please open an issue first to discuss the idea before spending time on
a pull request.

Before submitting a change, make sure the project still builds and the bundled learning path
still validates:

```sh
make build
./bin/shellgym validate --path paths/sample-linux-101
```

## Licensing of contributions

Shell Gym is distributed under the [PolyForm Noncommercial License 1.0.0](LICENSE.md), which
does not permit commercial use. To keep commercial licensing of Shell Gym possible, contributions
need a broader inbound grant than the license itself provides:

By submitting a contribution (a pull request, patch, or any other material), you license it to
Ivan Velichko (iximiuz Labs) under any terms, including commercial licensing and relicensing of
Shell Gym as a whole. You retain the copyright to your contribution and the right to use it for
any other purpose.

If you cannot or do not want to agree to this, please do not submit contributions - opening
issues with ideas and bug reports is still very much appreciated.
