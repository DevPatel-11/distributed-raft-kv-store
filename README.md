# Distributed Raft KV Store

A production-ready, distributed key-value store implementation using the Raft consensus algorithm. This project demonstrates comprehensive understanding of distributed systems, consensus mechanisms, and fault tolerance.

## 🎯 Features

- **Raft Consensus Algorithm**: Complete implementation of leader election, log replication, and commit logic
- **Distributed Key-Value Store**: Thread-safe KV operations across multiple nodes
- **Multi-Node Cluster**: Support for 3+ node clusters with automatic failover
- **REST API Gateway**: HTTP interface for client interactions
- **Persistence**: Log persistence with snapshot management (planned)
- **Fault Tolerance**: Continues operation during node failures and network partitions
- **Docker Support**: Multi-container deployment with Docker Compose

## 🏗️ Architecture

### Project Structure

```
distributed-raft-kv-store/
├── raft/                    # Raft consensus implementation
│   ├── types.go            # Core data types (Node, LogEntry, RPC types)
│   ├── node.go             # Node lifecycle and election logic
│   └── rpc.go              # RPC handlers (RequestVote, AppendEntries)
├── kvstore/                # Key-value store state machine
│   └── store.go            # In-memory KV store with metadata
├── cmd/
│   ├── raft-node/          # Raft node CLI
│   │   └── main.go         # Node entry point
│   └── kv-gateway/         # REST API gateway
│       └── main.go         # Gateway entry point with HTTP handlers
├── internal/               # Shared utilities
├── Makefile               # Build and deployment commands
├── docker-compose.yml     # Multi-node cluster configuration
├── .env.example          # Configuration template
├── go.mod                # Go module definition
└── README.md             # This file
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make

### Local Development

```bash
# Build binaries
make build

# Run 3-node cluster with Docker Compose
make run

# Test KV operations
curl http://localhost:8080/health
curl -X PUT http://localhost:8080/kv/mykey -H "Content-Type: application/json" -d '{"value": "myvalue"}'
curl http://localhost:8080/kv/mykey

# Stop services
make stop
```

## 📡 REST API Endpoints

### Health Check
```bash
GET /health
```
Response: `{"status": "healthy", "keys": "0"}`

### Set Key-Value
```bash
PUT /kv/{key}
Content-Type: application/json
{"value": "your-value"}
```

### Get Value
```bash
GET /kv/{key}
```
Response: `{"key": "mykey", "value": "myvalue"}`

### Delete Key
```bash
DELETE /kv/{key}
```
Response: `{"key": "mykey", "status": "deleted"}`

## 🔄 Raft Consensus Explained

### States
- **Follower**: Default state, receives RPCs from leader
- **Candidate**: Requests votes during election timeout
- **Leader**: Sends heartbeats and replicates log entries

### Key Mechanisms
- **Leader Election**: Random election timeouts (150-300ms) prevent split votes
- **Log Replication**: AppendEntries RPC ensures consistency across nodes
- **Safety**: Only entries committed to majority are applied to state machine

## 📚 Configuration

Edit `.env.example` and rename to `.env`:

```env
# Node Configuration
NODE_ID=1
NODE_ADDR=localhost:50051
CLUSTER_NODES=localhost:50051,localhost:50052,localhost:50053

# Timing Parameters
ELECTION_TIMEOUT_MIN=150        # milliseconds
ELECTION_TIMEOUT_MAX=300        # milliseconds
HEARTBEAT_INTERVAL=50           # milliseconds

# Gateway
GATEWAY_ADDR=0.0.0.0:8080
```

## 🐳 Docker Deployment

### Start 3-Node Cluster
```bash
docker-compose up -d

# Monitor logs
docker-compose logs -f
```

### Ports
- Raft Node 1: 50051
- Raft Node 2: 50052
- Raft Node 3: 50053
- REST API Gateway: 8080

## 🧪 Testing

### Unit Tests
```bash
make test
```

### Manual Testing
```bash
# Check health
curl http://localhost:8080/health

# Create keys
for i in {1..10}; do
  curl -X PUT http://localhost:8080/kv/key$i \
    -H "Content-Type: application/json" \
    -d '{"value": "value'$i'"}'
done

# Retrieve keys
curl http://localhost:8080/kv/key1

# Delete keys
curl -X DELETE http://localhost:8080/kv/key1
```

## 📈 Scalability

- Horizontal scaling by adding more Raft nodes
- Log compaction through snapshots (planned)
- Async replication for high throughput

## 🔐 Security Considerations

- TLS support for inter-node communication (planned)
- Authentication and authorization (planned)
- Rate limiting on API endpoints (planned)

## 🚧 Future Enhancements

- [ ] Snapshot management and log compaction
- [ ] TLS encryption for RPC
- [ ] gRPC for inter-node communication
- [ ] Persistence layer (RocksDB)
- [ ] Metrics and monitoring (Prometheus)
- [ ] Web UI for cluster visualization
- [ ] Benchmarking suite
- [ ] Configuration hot-reload

## 📖 Learning Resources

- [Raft Paper](https://raft.github.io/raft.pdf)
- [Raft Visualization](https://raft.github.io/raftscope/index.html)
- [etcd Raft Implementation](https://github.com/etcd-io/etcd/tree/main/raft)

## 📄 License

MIT License - See LICENSE file for details

## 👤 Author

DevPatel-11

---

**Status**: MVP Complete - Core Raft consensus and KV store operational
**Last Updated**: December 2025