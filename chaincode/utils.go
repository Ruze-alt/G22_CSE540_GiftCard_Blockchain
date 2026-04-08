package main

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// Helper function to get transaction timestamp
func getTimestamp(ctx contractapi.TransactionContextInterface) (string, error) {
	ts, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("failed to get transaction timestamp: %v", err)
	}

	t := time.Unix(ts.Seconds, int64(ts.Nanos)).UTC()
	return t.Format(time.RFC3339), nil
}

// Helper function to update card status after validation
func updateCardStatus(card *GiftCard, newStatus GiftCardStatus) error {
	err := assertValidStateTransition(card.Status, newStatus)
	if err != nil {
		return err
	}
	card.Status = newStatus
	return nil
}
