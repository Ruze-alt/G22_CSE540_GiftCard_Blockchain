# CSE 540: Gift Card Blockchain Tracker (Hyperledger Fabric)

**Group 22**: Mathm Alkaabi, Andrew Le, Dhanashree Somani, Shashank Singh, Shreya Marria

## Project Description

Hyperledger Fabric chaincode implementing a blockchain based gift card tracking system. 
Records the full lifecycle of each gift card: **creation** (issuer), **activation** (retailer), **redemption** (customer), and **audit history** (admin).

**Key Features:**
- Immutable provenance trail for each card
- Role based access control (issuer, retailer, customer, admin)
- Prevents duplicate redemption and unauthorized activation
- Queryable transaction history

## Dependencies 
- Docker
    - [Mac/Windows](https://docs.docker.com/desktop/) + Docker Compose (Comes With Desktop)
    - [Linux](https://docs.docker.com/engine/install/) + [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://go.dev/doc/install)
- [Git](https://git-scm.com/)
- [jq](https://jqlang.org/)
- [curl](https://curl.se/)

## How To Deploy

Use **two terminals** while running this project:

- **Terminal 1** = Hyperledger Fabric / chaincode commands
- **Terminal 2** = Web application / API server

Before starting, make sure **Docker is running** (for example, open Docker Desktop on Mac/Windows, or ensure Docker Engine is running on Linux) because the Fabric test network depends on Docker containers.

### Terminal 1 — Fabric and Chaincode Setup
#### 1. Clone or download this repository
```bash
git clone <repo-url>
cd <repo-folder>
```

#### 2. Install Hyperledger Fabric samples and binaries

Run:
```bash
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
./install-fabric.sh -h
```

The official Fabric install script downloads the Fabric samples, binaries, and Docker images needed to run the test network locally.

#### 3. Build the chaincode

From the project root:

```bash
cd chaincode
go clean -modcache
go mod tidy
go build
```

Running `go mod tidy` and `go build` first helps catch dependency or compile issues before deploying to the Fabric test network.

#### 4. Start the Fabric test network

From the `chaincode` directory, move into the test network directory:

```bash
cd .. && cd fabric-samples/test-network
./network.sh down
./network.sh up createChannel -ca
```

The Fabric test network supports creating a channel with Certificate Authorities using `up createChannel -ca`, which is the setup used by this project.

#### 5. Deploy the chaincode

```bash
./network.sh deployCC -ccn giftCard -ccp ../../chaincode -ccl go
```

The `deployCC` command packages, installs, approves, and commits the Go chaincode to the test network channel in one step.

#### 6. Load peer environment variables (Needed for command line invoking/querying, skip if only using web UI version)

```bash
source ./scripts/envVar.sh
setGlobals 1
```
This sets the peer CLI environment to use **Org1** by default for admin / issuer / retailer operations in the test network setup.

### Notes

- `setGlobals 1` = use **Org1MSP**
- `setGlobals 2` = use **Org2MSP**
- In this project, **Org1** acts as issuer / retailer / admin, and **Org2** acts as customer.

### Terminal 2 — Web Application Setup
In a second terminal, start the API server after the Fabric network and chaincode are already running.

```bash
export TEST_NETWORK_PATH=../fabric-samples/test-network

cd application
go mod tidy
go run ./cmd/server/
```

This starts the web application and REST API server locally. Once it is running, open:

```text
http://localhost:8080
```

## How To Use

### Web UI (Terminal 2)

Once Web UI running and opened you can interact with it and create, activate, transfer, redeem, suspend or reactivate gift cards, thus 
simulating and tracking the life cycle of a gift card throughout the chain. 

| Tab | Role | Available Actions |
|---|---|---|
| **Issuer** | Org1MSP | Create gift card |
| **Retailer** | Org1MSP | Activate card, Transfer card |
| **Customer** | Org2MSP | Redeem card |
| **Admin** | Org1MSP | Suspend card, Reactivate card |

All tabs share a **Look Up Gift Card** panel that displays card details, balance, and transaction history.

### REST API Examples

#### Create a card

```bash
curl -X POST http://localhost:8080/cards?role=issuer \
  -H 'Content-Type: application/json' \
  -d '{"cardID":"GC001","issuerID":"issuer1","balance":100}'
```

#### Get card details

```bash
curl http://localhost:8080/cards/GC001?role=issuer
```

#### Get card balance

```bash
curl http://localhost:8080/cards/GC001/balance?role=customer
```

#### Get card history

```bash
curl http://localhost:8080/cards/GC001/history?role=customer
```

### Direct Chaincode Commands (Terminal 1)

After deployment, keep **Terminal 1** open in:

```bash
fabric-samples/test-network
```

and make sure you already ran:

```bash
source ./scripts/envVar.sh
```

#### Use Org1 for issuer / retailer / admin actions

```bash
setGlobals 1
```

#### Create a gift card

```bash
peer chaincode invoke -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  -C mychannel -n giftCard \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
  -c '{"function":"CreateGiftCard","Args":["card1","issuer1","100"]}'
```

#### Activate a gift card

```bash
peer chaincode invoke -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  -C mychannel -n giftCard \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
  -c '{"function":"ActivateGiftCard","Args":["card1"]}'
```

#### Transfer a gift card to customer ownership

```bash
peer chaincode invoke -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  -C mychannel -n giftCard \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
  -c '{"function":"TransferGiftCard","Args":["card1","customer1"]}'
```

#### Switch to Org2 for customer actions

```bash
setGlobals 2
```

#### Redeem part of the balance

```bash
peer chaincode invoke -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  -C mychannel -n giftCard \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
  -c '{"function":"RedeemGiftCard","Args":["card1","25"]}'
```

#### Query card details

```bash
peer chaincode query -C mychannel -n giftCard \
  -c '{"function":"GetGiftCard","Args":["card1"]}'
```

#### Query current balance

```bash
peer chaincode query -C mychannel -n giftCard \
  -c '{"function":"GetCurrentBalance","Args":["card1"]}'
```

#### Query event history

```bash
peer chaincode query -C mychannel -n giftCard \
  -c '{"function":"GetGiftCardHistory","Args":["card1"]}'
```

#### Switch back to Org1 for admin actions

```bash
setGlobals 1
```

#### Suspend a card

```bash
peer chaincode invoke -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  -C mychannel -n giftCard \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
  -c '{"function":"SuspendGiftCard","Args":["card1","suspicious activity"]}'
```

#### Reactivate a card

```bash
peer chaincode invoke -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
  -C mychannel -n giftCard \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
  -c '{"function":"ReactivateGiftCard","Args":["card1"]}'
```