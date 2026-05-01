package main

import "fmt"

// Helper function to validate card ID format
func validateCardID(cardID string) error {
	if cardID == "" {
		return fmt.Errorf("card ID cannot be empty")
	}
	return nil
}

// Helper function to validate participant IDs (issuer, retailer, etc)
func validateParticipantID(id string, fieldName string) error {
	if id == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// Helper function to validate gift card amount (positive values only)
func validateAmount(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	return nil
}

// Helper function to check if the card has enough balance for redemption
func cardHasSufficientBalance(card *GiftCard, amount float64) bool {
	if card == nil {
		return false
	}
	return card.Balance >= amount
}

// Helper function to validate state transitions
func assertValidStateTransition(currentStatus GiftCardStatus, newStatus GiftCardStatus) error {
	switch currentStatus {
	case StatusCreated:
		if newStatus == StatusTransferred ||
			newStatus == StatusActivated ||
			newStatus == StatusSuspended {
			return nil
		}
	case StatusTransferred:
		if newStatus == StatusActivated ||
			newStatus == StatusSuspended {
			return nil
		}
	case StatusActivated:
		if newStatus == StatusPartiallyRedeemed ||
			newStatus == StatusRedeemed ||
			newStatus == StatusSuspended ||
			newStatus == StatusTransferred {
			return nil
		}
	case StatusPartiallyRedeemed:
		if newStatus == StatusPartiallyRedeemed ||
			newStatus == StatusRedeemed ||
			newStatus == StatusSuspended {
			return nil
		}
	case StatusSuspended:
		if newStatus == StatusActivated {
			return nil
		}
	}

	return fmt.Errorf("invalid status transition from %s to %s", currentStatus, newStatus)
}
