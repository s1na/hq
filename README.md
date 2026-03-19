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

### List clients and aliases

```
hq clients
```

The `--client` flag accepts shorthand aliases anywhere. For example, `--client geth` resolves to `go-ethereum`, `--client nimbus` to `nimbus-el`.

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
# Diffs for besu failures in the most recent rpc-compat run
hq diff --sim rpc-compat --client besu

# Narrow down to a specific test
hq diff --sim rpc-compat --client besu --test "debug_getRawHeader/*"

# Show full output (including raw request/response) instead of compact diffs
hq diff --sim rpc-compat --client besu --full

# Diff a specific run file
hq diff 1710000000-xxx.json --client geth
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

Find out which rpc-compat tests besu is failing and inspect the diffs:

```bash
# 1. See how each client is faring
hq stats --sim rpc-compat

# 2. List besu's failures
hq failures --sim rpc-compat --client besu

# 3. See the diffs for all failures
hq diff --sim rpc-compat --client besu

# 4. Narrow down to a specific test
hq diff --sim rpc-compat --client besu --test "debug_getRawHeader/*"
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
