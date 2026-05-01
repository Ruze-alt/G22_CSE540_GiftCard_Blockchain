package service

import (
    "giftcard/application/core/client"
    "giftcard/application/core/usecase"
)

// New returns the issuer role service. Hardcoded for now.
func New(gateway client.Gateway) *usecase.IssuerService {
    return usecase.NewIssuerService(gateway)
}
