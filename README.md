# Banyan - Bamboo Implementation

## What is Banyan?
Banyan is the first chained-BFT protocol that allows blocks to be confirmed in just a single round-trip time. It is integrated into the Internet Computer Consensus (ICC) protocol, without requiring any communication overhead. Crucially, even if the fast path is not effective, no penalties are incurred. 

This repository is a fork of bamboo: https://github.com/gitferry/bamboo

> **Bamboo** is a prototyping and evaluation framework that studies the next generation BFT (Byzantine fault-tolerant) protocols specific for blockchains, namely chained-BFT, or cBFT.
By leveraging Bamboo, developers can prototype a brand new cBFT protocol in around 300 LoC and evaluate using rich benchmark facilities.

> Bamboo details can be found in this [technical report](https://arxiv.org/abs/2103.00777). The paper appeared at [ICDCS 2021](https://icdcs2021.us/).

## What is included?

Protocols:
- [x] [HotStuff](https://dl.acm.org/doi/10.1145/3293611.3331591)
- [ ] [Two-chain HotStuff](https://dl.acm.org/doi/10.1145/3293611.3331591)
- [x] [Streamlet](https://dl.acm.org/doi/10.1145/3419614.3423256)
- [ ] [Fast-HotStuff](https://arxiv.org/abs/2010.11454)
- [ ] [LBFT](https://arxiv.org/abs/2012.01636)
- [ ] [SFT](https://arxiv.org/abs/2101.03715)
- [x] [Internet Computer Consensus](https://dl.acm.org/doi/abs/10.1145/3519270.3538430)
- [x] [Banyan](https://arxiv.org/html/2312.05869v1)

Features:
- [x] Benchmarking
- [x] Fault injection

## File Structure

```bash
aws/             # AWS configuration and deployment scripts
bin/deploy/      # Deployment scripts and utilities
bin/logs/        # Logs location
blockchain/      # Core blockchain implementation and logic
blockchain_view/ # View-change counterpart
config/          # config
crypto/          # Cryptographic utilities
election/        # Leader election mechanisms and algorithms
identity/        # Identity management
local_timeout/   # Logic for managing and handling local timeouts
log/             # log
message/         # Message structures
node/            # Core node functionality
pacemaker/       # Pacemaker and heartbeat mechanisms
protocol/        # Core protocol definitions and interactions
replica/         # Replica management and synchronization logic
server/          # Server-side logic and network handling
socket/          # Socket communication utilities
transport/       # Data transport mechanisms and utilities
types/           # Type definitions and shared data structures
utils/           # General utility functions and helpers
```


## How to build

1. Install [Go](https://golang.org/dl/).

2. Download the source code.

3. Build `server`.
```
cd bamboo/bin
go build ../server
```

# How to run

Users can run Bamboo-based cBFT protocols locally or on cloud infrastructure.

## Local

1. ```cd bamboo/bin```.
2. Put the name of the protocol you are going to run in `run_local.sh` (banyan, icc, hotstuff, streamlet).
3. Modify `ips.txt` with a set of IPs of each node. The number of IPs equals to the number of nodes. Here, the local IP is `127.0.0.1`. Each node will be assigned by an increasing port from `8070`.
4. Modify configuration parameters in `config.json`.
5. ```bash run_local.sh```.
6. In the browser go to `http://127.0.0.1:8070/query` to start the local run.
7. After the duration specified in `config.json`, go to `http://127.0.0.1:8070/query` to see the measured results.

Logs are produced in the local directory with the name `server.xxx.log` where `xxx` is the pid of the process.

## Cloud (AWS)
Bamboo can be deployed in a cloud network. We fork [hashrand-rs](https://github.com/akhilsb/hashrand-rs) (which itself is a fork of the Narwal benchmarking suite) to provide simple interaction with AWS.

1. To set up the AWS testbed follow the instructions at ```/bin/deploy/README.md```.
2. If needed adapt the settings by adjusting the `bin/deploy/config.json` file.
3. Run the `bash benchmarkRemote.sh` file. An example usage is described at the start of the file. To run banyan for the first time use: ```bash benchmarkRemote.sh banyan 0 1 1```.