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
- Hyperledger Fabric Samples
    ```bash
    # Bootstrap Fabric environment (once)
    mkdir ~/fabric-dev && cd ~/fabric-dev
    curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/master/scripts/bootstrap.sh | bash -s
    ```

## How To Deploy (WIP)
- **Build Go Files**
    ```bash
    cd chaincode
    go mod tidy
    go build
    ```
- **Start Test Network**
    ```bash
    cd ~/fabric-dev/fabric-samples/test-network
    ./network.sh down # Remove anything previously generated
    ./network.sh up
    ./network.sh createChannel
    ```

- **Package Smart Contract (Chaincode)**
    ```bash
    # Set Environment Path Variables
    export PATH=${PWD}/../bin:$PATH
    export FABRIC_CFG_PATH=${PWD}/../config/

    # Package Chaincode
    peer lifecycle chaincode package giftCard.tar.gz \
    --path <path/to/chaincode> \
    --lang golang \
    --label giftCard
    ```

- **Install Chaincode Package**
    ```bash
    # Set Environment Variables for Org1
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051

    # Install The Chaincode
    peer lifecycle chaincode install giftCard.tar.gz

    # Set Environment Variables for Org2
    export CORE_PEER_LOCALMSPID="Org2MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
    export CORE_PEER_ADDRESS=localhost:9051

    # Install The Chaincode
    peer lifecycle chaincode install giftCard.tar.gz
    ```

- **Approve Chaincode Definition**
    ```bash
    # Run Command And Copy Package ID
    peer lifecycle chaincode queryinstalled

    # Export Package ID
    export CC_PACKAGE_ID= # Your Copied Package ID

    # Approve Chaincode Definition As Org2
    peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name giftCard --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

    # Set Environment Variables To Operate As The Org1 Admin
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_ADDRESS=localhost:7051

    # Approve Chaincode Definition As Org1
    peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name giftCard --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
    ```

- **Commit Chaincode Definition To Channel**
    ```bash
    # Commit Chaincode
    peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name giftCard --version 1.0 --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
    ```

- **Invoking The Chaincode**
    ```bash
    # General Syntax for Invoking a Function
    peer chaincode invoke -o localhost:7050 \
    --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" \
    -C mychannel -n giftCard \
    --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
    --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" \
    -c '{"function":"YourFunctionName","Args":["Arg1", "Arg2"]}'
    ```

- **Invoke Functions To Use**
    - `InitLedger`
    - `CreateGiftCard`
    - `TransferGiftCard`
    - `ActivateGiftCard`
    - `RedeemGiftCard`
    - `SuspendGiftCard`
    - `ReactivateGiftCard`

- **Querying The Chaincode**
    ```bash
    # General Syntax for Querying
    peer chaincode query -C mychannel -n giftCard -c '{"Args":["FunctionName", "Arg1", "Arg2"]}'
    ```

- **Query Functions To Use**
    - `GetGiftCard`
    - `GetCurrentBalance`
    - `GetGiftCardHistory`

## Web Application

The project includes a Go REST API server and a single-page browser UI that connects to the Fabric test network. Complete the chaincode deployment steps above before starting the server.

### Additional Dependencies
- Go 1.22+ (required for `net/http` path parameters)
- No Node.js required — the frontend is pure HTML/CSS/JavaScript

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `TEST_NETWORK_PATH` | *(required)* | Absolute path to `fabric-samples/test-network` |
| `CHANNEL_NAME` | `mychannel` | Fabric channel name |
| `CHAINCODE_NAME` | `giftCard` | Deployed chaincode name |
| `PORT` | `8080` | HTTP port for the API server |
| `FRONTEND_DIR` | auto-detected | Path to `application/frontend/` (override if needed) |

### Starting the API Server
```bash
# Ensure the Fabric test network is running and chaincode is deployed first, then:
export TEST_NETWORK_PATH=~/fabric-dev/fabric-samples/test-network

cd application
go mod tidy
go run ./cmd/server/
```

The server starts on `http://localhost:8080` and serves both the REST API and the frontend UI.

### Opening the Frontend

Navigate to **http://localhost:8080** in your browser. The UI has four tabs:

| Tab | Role | Available Actions |
|---|---|---|
| **Issuer** | Org1MSP | Create gift card |
| **Retailer** | Org1MSP | Activate card, Transfer card |
| **Customer** | Org2MSP | Redeem card |
| **Admin** | Org1MSP | Suspend card, Reactivate card |

All tabs share a **Look Up Gift Card** panel that shows the current card status, balance, and full transaction history.

### REST API Endpoints

| Method | Path | Role param | Description |
|---|---|---|---|
| `POST` | `/cards` | `issuer` | Create a gift card |
| `POST` | `/cards/{id}/activate` | `retailer` | Activate a card |
| `POST` | `/cards/{id}/transfer` | `retailer` | Transfer card to new owner |
| `POST` | `/cards/{id}/redeem` | `customer` | Redeem an amount from card |
| `POST` | `/cards/{id}/suspend` | `admin` | Suspend a card |
| `POST` | `/cards/{id}/reactivate` | `admin` | Reactivate a suspended card |
| `GET` | `/cards/{id}` | any | Get card details |
| `GET` | `/cards/{id}/balance` | any | Get current balance |
| `GET` | `/cards/{id}/history` | any | Get event history |

**Example — create a card:**
```bash
curl -X POST http://localhost:8080/cards?role=issuer \
  -H 'Content-Type: application/json' \
  -d '{"cardID":"GC001","issuerID":"issuer1","balance":100}'
```

**Example — look up a card:**
```bash
curl http://localhost:8080/cards/GC001?role=issuer
```
