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
