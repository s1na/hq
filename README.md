# hq

`hq` (hive query) is a CLI tool for querying and investigating [Ethereum Hive](https://hive.ethpandaops.io) test results.

## Install

```
go install github.com/s1na/hq@latest
```

## Usage

### List available test suites

```
hq suites
```

### List recent runs

```
# All recent runs
hq runs

# Runs for a specific simulator
hq runs --sim rpc-compat

# Runs involving a specific client
hq runs --sim rpc-compat --client geth
```

### Show failures

```
# Geth failures in the most recent rpc-compat run
hq failures --sim rpc-compat --client geth

# Failures matching a test name pattern
hq failures --sim rpc-compat --client geth --test "eth_getBalance*"

# Failures for a specific run file
hq failures 1710000000-rpc-compat.json
```

### View result diffs

```
# Show diffs for all geth failures in a run
hq diff 1710000000-rpc-compat.json --client geth

# Show diff for a specific test
hq diff 1710000000-rpc-compat.json --client geth --test "eth_getStorageAt*"
```

### Pass/fail stats across runs

```
# Compare all clients on rpc-compat over the last 10 runs
hq stats --sim rpc-compat

# Track geth's rpc-compat pass rate over time
hq stats --sim rpc-compat --client geth

# Look at more history
hq stats --sim rpc-compat --client geth --last 20
```

### Example workflow

Find out which rpc-compat tests geth is failing and inspect the diffs:

```bash
# 1. List geth's failures in the most recent rpc-compat run
hq failures --sim rpc-compat --client geth

# 2. Note the run file from the output, then view the diffs
hq diff <run-file> --client geth
```

## Global flags

| Flag | Default | Description |
|------|---------|-------------|
| `--base-url` | `https://hive.ethpandaops.io` | Hive server base URL |
| `--suite` | `generic` | Test suite name |
| `--cache-dir` | `~/.cache/hq` | Cache directory |
| `--no-cache` | `false` | Bypass cache reads |
| `--no-color` | `false` | Disable colored output |

## License

MIT
