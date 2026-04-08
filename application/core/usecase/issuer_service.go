package usecase

import (
    "context"
    "encoding/json"
    "fmt"
    "strconv"

    "giftcard/application/core/client"
    "giftcard/application/core/models"
)

type IssuerService struct {
    gateway client.Gateway
}

func NewIssuerService(gateway client.Gateway) *IssuerService {
    return &IssuerService{gateway: gateway}
}

func (s *IssuerService) IssueGiftCard(ctx context.Context, cardID, issuerID string, balance float64) (string, error) {
    payload, err := s.gateway.Submit(ctx, "CreateGiftCard", cardID, issuerID, strconv.FormatFloat(balance, 'f', -1, 64))
    if err != nil {
        return "", fmt.Errorf("issue gift card failed: %w", err)
    }
    return string(payload), nil
}

func (s *IssuerService) ActivateGiftCard(ctx context.Context, cardID string) (string, error) {
    payload, err := s.gateway.Submit(ctx, "ActivateGiftCard", cardID)
    if err != nil {
        return "", fmt.Errorf("activate gift card failed: %w", err)
    }
    return string(payload), nil
}

func (s *IssuerService) RedeemGiftCard(ctx context.Context, cardID string, amount float64) (string, error) {
    payload, err := s.gateway.Submit(ctx, "RedeemGiftCard", cardID, strconv.FormatFloat(amount, 'f', -1, 64))
    if err != nil {
        return "", fmt.Errorf("redeem gift card failed: %w", err)
    }
    return string(payload), nil
}

func (s *IssuerService) GetGiftCard(ctx context.Context, cardID string) (*models.GiftCard, error) {
    payload, err := s.gateway.Evaluate(ctx, "GetGiftCard", cardID)
    if err != nil {
        return nil, fmt.Errorf("get gift card failed: %w", err)
    }

    var card models.GiftCard
    if err := json.Unmarshal(payload, &card); err != nil {
        return nil, fmt.Errorf("decode gift card failed: %w", err)
    }
    return &card, nil
}

func (s *IssuerService) GetGiftCardHistory(ctx context.Context, cardID string) ([]models.GiftCardEvent, error) {
    payload, err := s.gateway.Evaluate(ctx, "GetGiftCardHistory", cardID)
    if err != nil {
        return nil, fmt.Errorf("get gift card history failed: %w", err)
    }

    var events []models.GiftCardEvent
    if err := json.Unmarshal(payload, &events); err != nil {
        return nil, fmt.Errorf("decode gift card history failed: %w", err)
    }
    return events, nil
}
