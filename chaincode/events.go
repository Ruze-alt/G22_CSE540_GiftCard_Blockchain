package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// Helper function to record events on the ledger
func recordEvent(
	ctx contractapi.TransactionContextInterface,
	cardID string,
	eventType EventType,
	actorID string,
	actorMSP string,
	actorRole string,
	details string,
) error {
	txID := ctx.GetStub().GetTxID()

	timestamp, err := getTimestamp(ctx)
	if err != nil {
		return err
	}

	eventKey, err := ctx.GetStub().CreateCompositeKey(EventObjectType, []string{cardID, txID})
	if err != nil {
		return fmt.Errorf("failed to create composite key for event: %v", err)
	}

	event := GiftCardEvent{
		EventID:   eventKey,
		CardID:    cardID,
		EventType: eventType,
		ActorID:   actorID,
		ActorMSP:  actorMSP,
		ActorRole: actorRole,
		Timestamp: timestamp,
		TxID:      txID,
		Details:   details,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	err = ctx.GetStub().PutState(eventKey, eventJSON)
	if err != nil {
		return fmt.Errorf("failed to store event: %v", err)
	}

	notificationPayload := map[string]string{
		"cardId":    cardID,
		"eventType": string(eventType),
		"txId":      txID,
	}

	payload, err := json.Marshal(notificationPayload)
	if err == nil {
		ctx.GetStub().SetEvent(string(eventType), payload)
	}

	return nil
}

// Helper function to query events by card ID
func queryEventsByCardID(ctx contractapi.TransactionContextInterface, cardID string) ([]*GiftCardEvent, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(EventObjectType, []string{cardID})
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %v", err)
	}
	defer resultsIterator.Close()

	var events []*GiftCardEvent
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate event results: %v", err)
		}

		var event GiftCardEvent
		err = json.Unmarshal(queryResult.Value, &event)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal event: %v", err)
		}

		events = append(events, &event)
	}

	return events, nil
}
