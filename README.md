# PeerDAS custody benchmark

You may want to restart your beacon node with a new node id, how fast can you generate a new node ID that has the same [custody columns](https://github.com/ethereum/consensus-specs/blob/dev/specs/_features/eip7594/das-core.md#get_custody_columns) as your old node ID?

## Get started

Environment variables:
- `CUSTODY_SUBNET_COUNT` - the number of custody subnets to use, default is `1`, max is `DATA_COLUMN_SIDECAR_SUBNET_COUNT`
- `DATA_COLUMN_SIDECAR_SUBNET_COUNT` - [spec](https://github.com/ethereum/consensus-specs/blob/dev/specs/_features/eip7594/das-core.md#networking)
- `NUMBER_OF_COLUMNS` - [spec](https://github.com/ethereum/consensus-specs/blob/dev/specs/_features/eip7594/das-core.md#data-size)

```bash
# randomly generate a source node id and try and find a matching node id with the same custody columns
go run main.go

# find a matching node id when the source node CUSTODY_SUBNET_COUNT set to 4
CUSTODY_SUBNET_COUNT=4 go run main.go

# find a matching node id from a given node id
go run main.go 41d2dc1948d5b4eb2ddeebed81e22dda93a1fa866f90f490856cd7ccf9ed45a1
```

## Benchmarking

```bash
# default CUSTODY_SUBNET_COUNT=1
go test -bench=.

# CUSTODY_SUBNET_COUNT=4
CUSTODY_SUBNET_COUNT=4 go test -bench=.

# run 10 times
CUSTODY_SUBNET_COUNT=2 go test -bench=. -benchtime=10x
```

## My results

> **Note:** The benchmarks are run on a 32 thread cpu

Values at time of benchmarking;
- `DATA_COLUMN_SIDECAR_SUBNET_COUNT` - 32
- `NUMBER_OF_COLUMNS` - 128

Average time taken to find a matching node id;
- `CUSTODY_SUBNET_COUNT=1` - **~0.188ms**
- `CUSTODY_SUBNET_COUNT=2` - **~3.840ms**
- `CUSTODY_SUBNET_COUNT=3` - **~128.178ms**
- `CUSTODY_SUBNET_COUNT=4` - **~3.818s**
- `CUSTODY_SUBNET_COUNT=5` - **~136.136s**

### `CUSTODY_SUBNET_COUNT=1`

Average `0.188ms`

```bash
CUSTODY_SUBNET_COUNT=1 go test -bench=. -benchtime=10000x

goos: linux
goarch: amd64
pkg: github.com/savid/peerdas-custody-benchmark
cpu: AMD Ryzen 9 5950X 16-Core Processor            
BenchmarkFindMatchingNodeID-32    	   10000	    188186 ns/op
PASS
ok  	github.com/savid/peerdas-custody-benchmark	1.899s
```

### `CUSTODY_SUBNET_COUNT=2`

Average `3.840ms`

```bash
CUSTODY_SUBNET_COUNT=2 go test -bench=. -benchtime=10200x

goos: linux
goarch: amd64
pkg: github.com/savid/peerdas-custody-benchmark
cpu: AMD Ryzen 9 5950X 16-Core Processor            
BenchmarkFindMatchingNodeID-32    	    1000	   3839904 ns/op
PASS
ok  	github.com/savid/peerdas-custody-benchmark	3.854s
```

### `CUSTODY_SUBNET_COUNT=3`

Average `128.178ms`

```bash
CUSTODY_SUBNET_COUNT=3 go test -bench=. -benchtime=100x

goos: linux
goarch: amd64
pkg: github.com/savid/peerdas-custody-benchmark
cpu: AMD Ryzen 9 5950X 16-Core Processor            
BenchmarkFindMatchingNodeID-32    	     100	 128178191 ns/op
PASS
ok  	github.com/savid/peerdas-custody-benchmark	13.064s
```

### `CUSTODY_SUBNET_COUNT=4`

Average `3.818s`

```bash
CUSTODY_SUBNET_COUNT=4 go test -bench=. -benchtime=100x -timeout 60m

goos: linux
goarch: amd64
pkg: github.com/savid/peerdas-custody-benchmark
cpu: AMD Ryzen 9 5950X 16-Core Processor            
BenchmarkFindMatchingNodeID-32    	     100	3818105941 ns/op
PASS
ok  	github.com/savid/peerdas-custody-benchmark	387.489s
```

### `CUSTODY_SUBNET_COUNT=5`

Average `136.136s`

```bash
CUSTODY_SUBNET_COUNT=5 go test -bench=. -benchtime=10x -timeout 60m

goos: linux
goarch: amd64
pkg: github.com/savid/peerdas-custody-benchmark
cpu: AMD Ryzen 9 5950X 16-Core Processor            
BenchmarkFindMatchingNodeID-32    	      10	136136034750 ns/op
PASS
ok  	github.com/savid/peerdas-custody-benchmark	1396.396s
```

