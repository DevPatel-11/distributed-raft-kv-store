package main

import (
	"flag"
	"fmt"

	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/DevPatel-11/distributed-raft-kv-store/raft"
)

func main() {
	nodeID := flag.Int("id", 1, "Node ID")
	clusterAddr := flag.String("cluster", "localhost:50051,localhost:50052,localhost:50053", "Cluster addresses")
	flag.Parse()

	// Parse cluster addresses
	addresses := strings.Split(*clusterAddr, ",")
	clusterNodes := make(map[int]string)

	for i, addr := range addresses {
		clusterNodes[i+1] = strings.TrimSpace(addr)
	}

	fmt.Printf("Starting Raft node %d...\n", *nodeID)
	fmt.Printf("Cluster nodes: %v\n", clusterNodes)

	// Create a new Raft node
	node := raft.NewRaftNode(*nodeID, clusterNodes)

	// Start the Raft node
	node.Start()
	fmt.Printf("[Node %d] Raft node started\n", node.NodeID)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Printf("\n[Node %d] Shutting down...\n", node.NodeID)
}
