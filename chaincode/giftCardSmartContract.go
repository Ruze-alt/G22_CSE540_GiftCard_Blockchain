package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type GiftCardSmartContract struct {
	contractapi.Contract
}

// Helper function to check if a gift card exists in the ledger
func (gc *GiftCardSmartContract) GiftCardExists(
	ctx contractapi.TransactionContextInterface,
	cardID string,
) (bool, error) {
	cardJSON, err := ctx.GetStub().GetState(cardID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return cardJSON != nil, nil
}

// Function CreateGiftCard allows an issuer or admin to create a new gift card
func (gc *GiftCardSmartContract) CreateGiftCard(
	ctx contractapi.TransactionContextInterface,
	cardID string,
	issuerID string,
	balance float64,
) (string, error) {
	err := requireMSP(ctx, Org1MSP)
	if err != nil {
		return "", err
	}

	err = validateCardID(cardID)
	if err != nil {
		return "", err
	}

	err = validateParticipantID(issuerID, "issuer ID")
	if err != nil {
		return "", err
	}

	err = validateAmount(balance)
	if err != nil {
		return "", err
	}

	exists, err := gc.GiftCardExists(ctx, cardID)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("gift card %s already exists", cardID)
	}

	info, err := getClientIdentity(ctx)
	if err != nil {
		return "", err
	}

	now, err := getTimestamp(ctx)
	if err != nil {
		return "", err
	}

	card := &GiftCard{
		CardID:          cardID,
		OwnerID:         issuerID,
		OwnerMSP:        info.MSPID,
		IssuerID:        issuerID,
		IssuerMSP:       info.MSPID,
		Balance:         balance,
		OriginalBalance: balance,
		Status:          StatusCreated,
		CreatedAt:       now,
		LastUpdatedAt:   now,
	}

	err = writeGiftCardState(ctx, card)
	if err != nil {
		return "", err
	}

	err = recordEvent(
		ctx,
		cardID,
		EventTypeCardCreated,
		info.ClientID,
		info.MSPID,
		string(RoleIssuer),
		"gift card created by issuer/admin",
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Gift card %s created successfully", cardID), nil
}

// Function TransferGiftCard allows the current owner or an admin to transfer the gift card to a new owner
func (gc *GiftCardSmartContract) TransferGiftCard(
	ctx contractapi.TransactionContextInterface,
	cardID string,
	newOwnerID string,
) (string, error) {
	err := requireMSP(ctx, Org1MSP)
	if err != nil {
		return "", err
	}

	if err := validateCardID(cardID); err != nil {
		return "", err
	}

	err = validateParticipantID(newOwnerID, "new owner ID")
	if err != nil {
		return "", err
	}

	card, err := readGiftCardState(ctx, cardID)
	if err != nil {
		return "", err
	}

	if card.Status != StatusActivated {
		return "", fmt.Errorf("gift card must be activated before transfer to customer")
	}

	info, err := getClientIdentity(ctx)
	if err != nil {
		return "", err
	}

	card.OwnerID = newOwnerID
	card.OwnerMSP = Org2MSP
	card.Status = StatusTransferred

	now, err := getTimestamp(ctx)
	if err != nil {
		return "", err
	}
	card.LastUpdatedAt = now

	err = writeGiftCardState(ctx, card)
	if err != nil {
		return "", err
	}

	err = recordEvent(
		ctx,
		cardID,
		EventTypeCardTransferred,
		info.ClientID,
		info.MSPID,
		string(RoleRetailer),
		fmt.Sprintf("gift card transferred to customer %s", newOwnerID),
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Gift card %s transferred to %s", cardID, newOwnerID), nil
}

// Function ActivateGiftCard allows a retailer or admin to activate a created gift card
func (gc *GiftCardSmartContract) ActivateGiftCard(
	ctx contractapi.TransactionContextInterface,
	cardID string,
) (string, error) {
	err := requireMSP(ctx, Org1MSP)
	if err != nil {
		return "", err
	}

	err = validateCardID(cardID)
	if err != nil {
		return "", err
	}

	card, err := readGiftCardState(ctx, cardID)
	if err != nil {
		return "", err
	}

	err = updateCardStatus(card, StatusActivated)
	if err != nil {
		return "", err
	}

	info, err := getClientIdentity(ctx)
	if err != nil {
		return "", err
	}

	now, err := getTimestamp(ctx)
	if err != nil {
		return "", err
	}

	card.RetailerID = info.ClientID
	card.RetailerMSP = info.MSPID
	card.ActivatedAt = now
	card.LastUpdatedAt = now
	card.OwnerID = "retailer1" // For demo

	err = writeGiftCardState(ctx, card)
	if err != nil {
		return "", err
	}

	err = recordEvent(
		ctx,
		cardID,
		EventTypeCardActivated,
		info.ClientID,
		info.MSPID,
		string(RoleRetailer),
		"gift card activated by retailer/admin",
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Gift card %s activated successfully", cardID), nil
}

// Function RedeemGiftCard allows the current owner to redeem a gift card
func (gc *GiftCardSmartContract) RedeemGiftCard(
	ctx contractapi.TransactionContextInterface,
	cardID string,
	amount float64,
) (string, error) {
	var err error
	err = validateCardID(cardID)
	if err != nil {
		return "", err
	}
	err = validateAmount(amount)
	if err != nil {
		return "", err
	}

	card, err := readGiftCardState(ctx, cardID)
	if err != nil {
		return "", err
	}

	info, err := getClientIdentity(ctx)
	if err != nil {
		return "", err
	}

	if info.MSPID != Org1MSP && info.MSPID != card.OwnerMSP {
		return "", fmt.Errorf("caller is not authorized to redeem this gift card")
	}

	if card.Status == StatusSuspended {
		return "", fmt.Errorf("gift card is suspended")
	}

	if card.Status != StatusActivated && card.Status != StatusTransferred && card.Status != StatusPartiallyRedeemed {
		return "", fmt.Errorf("gift card is not in a redeemable state")
	}

	if !cardHasSufficientBalance(card, amount) {
		return "", fmt.Errorf("insufficient balance")
	}

	card.Balance -= amount
	if card.Balance == 0 {
		card.Status = StatusRedeemed
	} else {
		card.Status = StatusPartiallyRedeemed
	}

	now, err := getTimestamp(ctx)
	if err != nil {
		return "", err
	}
	card.LastUpdatedAt = now

	err = writeGiftCardState(ctx, card)
	if err != nil {
		return "", err
	}

	err = recordEvent(
		ctx,
		cardID,
		EventTypeCardRedeemed,
		info.ClientID,
		info.MSPID,
		getActorRoleByMSP(info.MSPID),
		fmt.Sprintf("gift card redeemed for amount %.2f", amount),
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Gift card %s redeemed for %.2f, new balance: %.2f", cardID, amount, card.Balance), nil
}

// Function SuspendGiftCard allows an admin to suspend a gift card
func (gc *GiftCardSmartContract) SuspendGiftCard(
	ctx contractapi.TransactionContextInterface,
	cardID string,
	reason string,
) (string, error) {
	err := requireMSP(ctx, Org1MSP)
	if err != nil {
		return "", err
	}

	err = validateCardID(cardID)
	if err != nil {
		return "", err
	}

	card, err := readGiftCardState(ctx, cardID)
	if err != nil {
		return "", err
	}

	err = updateCardStatus(card, StatusSuspended)
	if err != nil {
		return "", err
	}

	now, err := getTimestamp(ctx)
	if err != nil {
		return "", err
	}
	card.LastUpdatedAt = now

	info, err := getClientIdentity(ctx)
	if err != nil {
		return "", err
	}

	err = writeGiftCardState(ctx, card)
	if err != nil {
		return "", err
	}

	err = recordEvent(
		ctx,
		cardID,
		EventTypeCardSuspended,
		info.ClientID,
		info.MSPID,
		string(RoleAdmin),
		reason,
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Gift card %s suspended: %s", cardID, reason), nil
}

// Function ReactivateGiftCard allows an admin to reactivate a suspended gift card
func (gc *GiftCardSmartContract) ReactivateGiftCard(
	ctx contractapi.TransactionContextInterface,
	cardID string,
) (string, error) {
	err := requireMSP(ctx, Org1MSP)
	if err != nil {
		return "", err
	}

	err = validateCardID(cardID)
	if err != nil {
		return "", err
	}

	card, err := readGiftCardState(ctx, cardID)
	if err != nil {
		return "", err
	}

	if card.Status != StatusSuspended {
		return "", fmt.Errorf("gift card is not suspended")
	}

	card.Status = StatusActivated

	now, err := getTimestamp(ctx)
	if err != nil {
		return "", err
	}
	card.LastUpdatedAt = now

	info, err := getClientIdentity(ctx)
	if err != nil {
		return "", err
	}

	err = writeGiftCardState(ctx, card)
	if err != nil {
		return "", err
	}

	err = recordEvent(
		ctx,
		cardID,
		EventTypeCardReactivated,
		info.ClientID,
		info.MSPID,
		string(RoleAdmin),
		"gift card reactivated by admin",
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Gift card %s reactivated", cardID), nil
}

// Function GetGiftCard allows the owner or an admin to view gift card details
func (gc *GiftCardSmartContract) GetGiftCard(
	ctx contractapi.TransactionContextInterface,
	cardID string,
) (*GiftCard, error) {
	err := validateCardID(cardID)
	if err != nil {
		return nil, err
	}

	card, err := readGiftCardState(ctx, cardID)
	if err != nil {
		return nil, err
	}

	info, err := getClientIdentity(ctx)
	if err != nil {
		return nil, err
	}

	if info.MSPID != Org1MSP && info.MSPID != card.OwnerMSP {
		return nil, fmt.Errorf("caller is not authorized to view this gift card")
	}

	return card, nil
}

// Function GetCurrentBalance allows the owner or an admin to view the current balance of the gift card
func (gc *GiftCardSmartContract) GetCurrentBalance(
	ctx contractapi.TransactionContextInterface,
	cardID string,
) (float64, error) {
	card, err := gc.GetGiftCard(ctx, cardID)
	if err != nil {
		return 0, err
	}
	return card.Balance, nil
}

// Function GetGiftCardHistory allows the owner or an admin to view the history of events related to the gift card
func (gc *GiftCardSmartContract) GetGiftCardHistory(
	ctx contractapi.TransactionContextInterface,
	cardID string,
) ([]*GiftCardEvent, error) {
	var err error
	err = validateCardID(cardID)
	if err != nil {
		return nil, err
	}

	card, err := readGiftCardState(ctx, cardID)
	if err != nil {
		return nil, err
	}

	info, err := getClientIdentity(ctx)
	if err != nil {
		return nil, err
	}

	if info.MSPID != Org1MSP && info.MSPID != card.OwnerMSP {
		return nil, fmt.Errorf("caller is not authorized to view gift card history")
	}

	return queryEventsByCardID(ctx, cardID)
}
