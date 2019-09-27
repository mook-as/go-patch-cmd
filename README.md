# go-patch-cmd

This is a simple command line client for [go-patch], to test out ops files.

[go-patch]: https://github.com/cppforlife/go-patch

## Usage

```sh
go-patch-cmd target.yaml ops.yaml ops.yaml...
```

The first argument is the YAML file to apply the patch to.

The rest of the arguments are ops files to apply.

Optionally, `--output <file>` may be given to emit to the output to the given
file, rather than standard output.
