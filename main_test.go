package main

import (
	"testing"
)

func BenchmarkFindMatchingNodeID(b *testing.B) {
	initialNodeID := generateNodeID()

	for i := 0; i < b.N; i++ {
		if _, err := findMatchingNodeID(initialNodeID, CustodySubnetCount); err != nil {
			b.Fatal(err)
		}
	}
}
