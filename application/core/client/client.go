package client

import (
    "context"
)

// Gateway defines the blockchain operations used by the role apps.

type Gateway interface {
    Submit(ctx context.Context, transactionName string, args ...string) ([]byte, error)
    Evaluate(ctx context.Context, transactionName string, args ...string) ([]byte, error)
}
