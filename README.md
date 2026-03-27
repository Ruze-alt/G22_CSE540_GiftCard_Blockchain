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

## How To Deploy / Use (WIP)
- **Build Chaincode**
    ```bash
    cd chaincode
    go mod tidy
    go build
    ```
- **Start Test Network**
    ```bash
    cd ~/fabric-dev/fabric-samples/test-network
    ./network.sh up createChannel
    ```

- **Chaincode Functions**
    ```bash
    CreateGiftCard(cardId, ownerId, balance) # Issuer Creates Card
    ActivateGiftCard(cardId) # Retailer Activates
    RedeemGiftCard(cardId, amount) # Customer Redeems
    GetGiftCardHistory(cardId) # Full History
    ```
