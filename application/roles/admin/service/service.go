package service

import (
    "context"
    "encoding/json"
    "fmt"

    "giftcard/application/core/client"
    "giftcard/application/core/models"
)

type Service struct {
    gateway client.Gateway
}

func New(gateway client.Gateway) *Service {
    return &Service{gateway: gateway}
}

func (s *Service) SuspendGiftCard(ctx context.Context, cardID, reason string) (string, error) {
    payload, err := s.gateway.Submit(ctx, "SuspendGiftCard", cardID, reason)
    if err != nil {
        return "", fmt.Errorf("suspend gift card failed: %w", err)
    }
    return string(payload), nil
}

func (s *Service) ReactivateGiftCard(ctx context.Context, cardID string) (string, error) {
    payload, err := s.gateway.Submit(ctx, "ReactivateGiftCard", cardID)
    if err != nil {
        return "", fmt.Errorf("reactivate gift card failed: %w", err)
    }
    return string(payload), nil
}

func (s *Service) GetGiftCard(ctx context.Context, cardID string) (*models.GiftCard, error) {
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

func (s *Service) GetGiftCardHistory(ctx context.Context, cardID string) ([]models.GiftCardEvent, error) {
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
