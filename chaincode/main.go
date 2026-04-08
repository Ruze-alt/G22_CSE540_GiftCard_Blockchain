package main

import (
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// Main function to start the gift card chaincode
func main() {
	chaincode, err := contractapi.NewChaincode(&GiftCardSmartContract{})
	if err != nil {
		panic("error creating gift card chaincode: " + err.Error())
	}

	err = chaincode.Start()
	if err != nil {
		panic("error starting gift card chaincode: " + err.Error())
	}
}
