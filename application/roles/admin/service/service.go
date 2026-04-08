package service

import "giftcard/application/core/client"

// Service is intentionally skeletal in this phase.
type Service struct {
    gateway client.Gateway
}

func New(gateway client.Gateway) *Service {
    return &Service{gateway: gateway}
}

// Future admin methods can be added here without changing the shared layer.
