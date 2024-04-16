package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

var (
	DataColumnSidecarSubnetCount uint64 // https://github.com/ethereum/consensus-specs/blob/dev/specs/_features/eip7594/das-core.md#networking
	NumberOfColumns              uint64 // https://github.com/ethereum/consensus-specs/blob/dev/specs/_features/eip7594/das-core.md#data-size
	CustodySubnetCount           uint64 // between 0 and DataColumnSidecarSubnetCount
)

func init() {
	DataColumnSidecarSubnetCount = getEnvAsUint64("DATA_COLUMN_SIDECAR_SUBNET_COUNT", 32)
	NumberOfColumns = getEnvAsUint64("NUMBER_OF_COLUMNS", 128)
	CustodySubnetCount = getEnvAsUint64("CUSTODY_SUBNET_COUNT", 1)
}

func getEnvAsUint64(name string, defaultVal uint64) uint64 {
	if valueStr, exists := os.LookupEnv(name); exists {
		value, err := strconv.ParseUint(valueStr, 10, 64)
		if err == nil {
			return value
		}
		fmt.Fprintf(os.Stderr, "Error parsing environment variable %s: %v\nUsing default value %d\n", name, err, defaultVal)
	}
	return defaultVal
}

type ColumnIndex int

func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func uintToBytes(i uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func getHash(data []byte) []byte {
	hasher := sha256.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}

// https://github.com/ethereum/consensus-specs/blob/dev/specs/_features/eip7594/das-core.md#helper-functions
func getCustodyColumns(nodeID enode.ID, custodySubnetCount uint64) ([]ColumnIndex, error) {
	if custodySubnetCount > DataColumnSidecarSubnetCount {
		return nil, fmt.Errorf("custodySubnetCount must not exceed DataColumnSidecarSubnetCount")
	}

	var subnetIDs []uint64
	for i := uint64(0); uint64(len(subnetIDs)) < custodySubnetCount; i++ {
		nodeIDUint := bytesToUint64(nodeID[:]) + i
		subnetID := bytesToUint64(getHash(uintToBytes(nodeIDUint))[:8]) % DataColumnSidecarSubnetCount

		if unique := !contains(subnetIDs, subnetID); unique {
			subnetIDs = append(subnetIDs, subnetID)
		}
	}

	if len(subnetIDs) != int(custodySubnetCount) {
		return nil, fmt.Errorf("subnetIDs calculation error")
	}

	columnsPerSubnet := NumberOfColumns / DataColumnSidecarSubnetCount
	var result []ColumnIndex
	for _, subnetID := range subnetIDs {
		for j := 0; j < int(columnsPerSubnet); j++ {
			result = append(result, ColumnIndex(DataColumnSidecarSubnetCount*uint64(j)+subnetID))
		}
	}

	return result, nil
}

func contains(slice []uint64, value uint64) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func generateNodeID() enode.ID {
	privateKey, _ := crypto.GenerateKey()
	return enode.PubkeyToIDV4(&privateKey.PublicKey)
}

// findMatchingNodeID finds a node ID that has the same custody columns as the source node ID
func findMatchingNodeID(sourceNodeID enode.ID, custodySubnetCount uint64) (*enode.ID, error) {
	sourceColumns, err := getCustodyColumns(sourceNodeID, custodySubnetCount)
	if err != nil {
		return nil, err
	}

	for {
		privateKey, _ := crypto.GenerateKey()
		nodeID := enode.PubkeyToIDV4(&privateKey.PublicKey)

		testColumns, err := getCustodyColumns(nodeID, custodySubnetCount)
		if err != nil {
			return nil, err
		}
		if fmt.Sprint(sourceColumns) == fmt.Sprint(testColumns) {
			return &nodeID, nil
		}
	}
}

func main() {
	var sourceNodeID enode.ID
	if len(os.Args) > 1 {
		nodeIDStr := os.Args[1]
		var err error
		if sourceNodeID, err = enode.ParseID(nodeIDStr); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid node ID format: %v\n", err)
			os.Exit(1)
		}
	} else {
		sourceNodeID = generateNodeID()
	}

	fmt.Println("custody subnet count:", CustodySubnetCount)
	fmt.Println("Source node ID:", sourceNodeID)
	sourceColumns, err := getCustodyColumns(sourceNodeID, CustodySubnetCount)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting custody columns: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Source columns:", sourceColumns)
	fmt.Println("Starting search...")

	start := time.Now()
	if nodeID, err := findMatchingNodeID(sourceNodeID, CustodySubnetCount); err != nil {
		fmt.Fprintf(os.Stderr, "Error finding matching node ID: %v\n", err)
		return
	} else {
		elapsed := time.Since(start)
		fmt.Printf("Search took %s\n", elapsed)
		fmt.Printf("Found node ID: %v\n", *nodeID)
	}
}
