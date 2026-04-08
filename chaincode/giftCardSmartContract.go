package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type GiftCardSmartContract struct {
	contractapi.Contract
}

// InitLedger may be used for demo purposes
func (gc *GiftCardSmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
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
	ownerID string,
	retailerID string,
	balance float64,
) (*GiftCard, error) {
	err := requireMSP(ctx, []string{Org1MSP})
	if err != nil {
		return nil, err
	}

	err = validateCardID(cardID)
	if err != nil {
		return nil, err
	}

	err = validateParticipantID(ownerID, "owner ID")
	if err != nil {
		return nil, err
	}

	err = validateParticipantID(retailerID, "retailer ID")
	if err != nil {
		return nil, err
	}

	err = validateAmount(balance)
	if err != nil {
		return nil, err
	}

	exists, err := gc.GiftCardExists(ctx, cardID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("gift card %s already exists", cardID)
	}

	info, err := getClientIdentity(ctx)
	if err != nil {
		return nil, err
	}

	now, err := getTimestamp(ctx)
	if err != nil {
		return nil, err
	}

	card := &GiftCard{
		CardID:          cardID,
		OwnerID:         ownerID,
		OwnerMSP:        Org2MSP,
		IssuerID:        info.ClientID,
		IssuerMSP:       info.MSPID,
		RetailerID:      retailerID,
		RetailerMSP:     Org1MSP,
		Balance:         balance,
		OriginalBalance: balance,
		Status:          StatusCreated,
		CreatedAt:       now,
		LastUpdatedAt:   now,
	}

	err = writeGiftCardState(ctx, card)
	if err != nil {
		return nil, err
	}

	err = recordEvent(
		ctx,
		cardID,
		EventTypeCardCreated,
		info.ClientID,
		info.MSPID,
		getActorRoleByMSP(info.MSPID),
		"gift card created by issuer/admin",
	)
	if err != nil {
		return nil, err
	}

	return card, nil
}

// Function TransferGiftCard allows the current owner or an admin to transfer the gift card to a new owner
func (gc *GiftCardSmartContract) TransferGiftCard(
	ctx contractapi.TransactionContextInterface,
	cardID string,
	newOwnerID string,
) (*GiftCard, error) {
	var err error
	err = validateCardID(cardID)
	if err != nil {
		return nil, err
	}
	err = validateParticipantID(newOwnerID, "new owner ID")
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

	if info.MSPID != Org1MSP && !(info.MSPID == card.OwnerMSP && info.ClientID == card.OwnerID) {
		return nil, fmt.Errorf("caller is not authorized to transfer this gift card")
	}

	err = assertValidStateTransition(card.Status, StatusTransferred)
	if err != nil {
		return nil, err
	}

	card.OwnerID = newOwnerID
	card.OwnerMSP = Org2MSP
	card.Status = StatusTransferred

	now, err := getTimestamp(ctx)
	if err != nil {
		return nil, err
	}
	card.LastUpdatedAt = now

	err = writeGiftCardState(ctx, card)
	if err != nil {
		return nil, err
	}

	err = recordEvent(
		ctx,
		cardID,
		EventTypeCardTransferred,
		info.ClientID,
		info.MSPID,
		getActorRoleByMSP(info.MSPID),
		fmt.Sprintf("gift card transferred to new owner %s", newOwnerID),
	)
	if err != nil {
		return nil, err
	}

	return card, nil
}

// Function ActivateGiftCard allows a retailer or admin to activate a created gift card
func (gc *GiftCardSmartContract) ActivateGiftCard(
	ctx contractapi.TransactionContextInterface,
	cardID string,
) (*GiftCard, error) {
	var err error
	err = requireMSP(ctx, []string{Org1MSP})
	if err != nil {
		return nil, err
	}

	err = validateCardID(cardID)
	if err != nil {
		return nil, err
	}

	card, err := readGiftCardState(ctx, cardID)
	if err != nil {
		return nil, err
	}

	err = updateCardStatus(card, StatusActivated)
	if err != nil {
		return nil, err
	}

	now, err := getTimestamp(ctx)
	if err != nil {
		return nil, err
	}
	card.ActivatedAt = now
	card.LastUpdatedAt = now

	info, err := getClientIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = writeGiftCardState(ctx, card)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return card, nil
}

// Function RedeemGiftCard allows the current owner to redeem a gift card
func (gc *GiftCardSmartContract) RedeemGiftCard(
	ctx contractapi.TransactionContextInterface,
	cardID string,
	amount float64,
) (*GiftCard, error) {
	var err error
	err = validateCardID(cardID)
	if err != nil {
		return nil, err
	}
	err = validateAmount(amount)
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

	if info.MSPID != Org1MSP && !(info.MSPID == card.OwnerMSP && info.ClientID == card.OwnerID) {
		return nil, fmt.Errorf("caller is not authorized to redeem this gift card")
	}

	if card.Status == StatusSuspended {
		return nil, fmt.Errorf("gift card is suspended")
	}

	if card.Status != StatusActivated && card.Status != StatusPartiallyRedeemed {
		return nil, fmt.Errorf("gift card is not in a redeemable state")
	}

	if !cardHasSufficientBalance(card, amount) {
		return nil, fmt.Errorf("insufficient balance")
	}

	card.Balance -= amount
	if card.Balance == 0 {
		card.Status = StatusRedeemed
	} else {
		card.Status = StatusPartiallyRedeemed
	}

	now, err := getTimestamp(ctx)
	if err != nil {
		return nil, err
	}
	card.LastUpdatedAt = now

	err = writeGiftCardState(ctx, card)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return card, nil
}

// Function SuspendGiftCard allows an admin to suspend a gift card
func (gc *GiftCardSmartContract) SuspendGiftCard(
	ctx contractapi.TransactionContextInterface,
	cardID string,
	reason string,
) (*GiftCard, error) {
	var err error
	err = requireMSP(ctx, []string{Org1MSP})
	if err != nil {
		return nil, err
	}

	err = validateCardID(cardID)
	if err != nil {
		return nil, err
	}

	card, err := readGiftCardState(ctx, cardID)
	if err != nil {
		return nil, err
	}

	err = updateCardStatus(card, StatusSuspended)
	if err != nil {
		return nil, err
	}

	now, err := getTimestamp(ctx)
	if err != nil {
		return nil, err
	}
	card.LastUpdatedAt = now

	info, err := getClientIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = writeGiftCardState(ctx, card)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return card, nil
}

// Function ReactivateGiftCard allows an admin to reactivate a suspended gift card
func (gc *GiftCardSmartContract) ReactivateGiftCard(
	ctx contractapi.TransactionContextInterface,
	cardID string,
) (*GiftCard, error) {
	var err error
	err = requireMSP(ctx, []string{Org1MSP})
	if err != nil {
		return nil, err
	}

	err = validateCardID(cardID)
	if err != nil {
		return nil, err
	}

	card, err := readGiftCardState(ctx, cardID)
	if err != nil {
		return nil, err
	}

	if card.Status != StatusSuspended {
		return nil, fmt.Errorf("gift card is not suspended")
	}

	card.Status = StatusActivated

	now, err := getTimestamp(ctx)
	if err != nil {
		return nil, err
	}
	card.LastUpdatedAt = now

	info, err := getClientIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = writeGiftCardState(ctx, card)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return card, nil
}

// Function GetGiftCard allows the owner or an admin to view gift card details
func (gc *GiftCardSmartContract) GetGiftCard(
	ctx contractapi.TransactionContextInterface,
	cardID string,
) (*GiftCard, error) {
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

	if info.MSPID != Org1MSP && !(info.MSPID == card.OwnerMSP && info.ClientID == card.OwnerID) {
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

	if info.MSPID != Org1MSP && !(info.MSPID == card.OwnerMSP && info.ClientID == card.OwnerID) {
		return nil, fmt.Errorf("caller is not authorized to view gift card history")
	}

	return queryEventsByCardID(ctx, cardID)
}
