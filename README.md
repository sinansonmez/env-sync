# env-sync

env-sync is a small CLI intended to compare and generate a destination `.env` file from a
source template. The source file (like file from test or uat environment file) provides the canonical key order; the
destination file (like prod environment file) contributes existing values. Optional flags control whether to
keep destination-only keys, fill missing or empty values, and dry-run output.

Status: ready (CLI flags and sync logic implemented).

## Install

Homebrew:

```bash
brew tap sinansonmez/tools
brew install env-sync
env-sync -h
```

Go:

```bash
go install ./...
```

## Usage

```bash
env-sync -source .env.uat -dest .env.prod
```

## Flags

- `-source` (default: `.env.uat`) path to source env file (uat/test/dev)
- `-dest` (default: `.env.prod`) path to destination env file (prod)
- `-dry-run` print result; do not write destination
- `-keep-unused` append keys found only in destination to the end of output
- `-use-source-defaults` when a key is missing in destination, keep the default
  value from source instead of blank
- `-fill-empty` if a key exists in destination but value is empty, fill from
  source (still does not overwrite non-empty)

## Examples

Print the result without writing:

```bash
env-sync -source .env.uat -dest .env.prod -dry-run
```

Preserve destination-only keys, but fill missing values from the source:

```bash
env-sync -use-source-defaults -fill-empty
```
