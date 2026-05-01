package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// Helper function to read gift card state from the ledger
func readGiftCardState(ctx contractapi.TransactionContextInterface, cardID string) (*GiftCard, error) {
	cardJSON, err := ctx.GetStub().GetState(cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from state: %v", err)
	}

	if cardJSON == nil {
		return nil, fmt.Errorf("gift card %s does not exist", cardID)
	}

	var giftCard GiftCard
	err = json.Unmarshal(cardJSON, &giftCard)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal gift card: %v", err)
	}

	return &giftCard, nil
}

// Helper function to write gift card state to the ledger
func writeGiftCardState(ctx contractapi.TransactionContextInterface, giftCard *GiftCard) error {
	cardJSON, err := json.Marshal(giftCard)
	if err != nil {
		return fmt.Errorf("failed to marshal gift card: %v", err)
	}

	return ctx.GetStub().PutState(giftCard.CardID, cardJSON)
}
