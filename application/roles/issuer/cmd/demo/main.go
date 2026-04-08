package main

import (
    "context"
    "fmt"
    "time"

    "giftcard/application/core/client"
    issuerrole "giftcard/application/roles/issuer/service"
)

// mockGateway is a placeholder for a real Fabric gateway adapter.
// Replace this with a Fabric SDK implementation when wiring the app.
type mockGateway struct{}

func (m *mockGateway) Submit(ctx context.Context, transactionName string, args ...string) ([]byte, error) {
    return []byte(fmt.Sprintf("mock submit: %s %v", transactionName, args)), nil
}

func (m *mockGateway) Evaluate(ctx context.Context, transactionName string, args ...string) ([]byte, error) {
    if transactionName == "GetGiftCard" {
        return []byte(`{"cardID":"GC1001","issuerID":"issuer1","ownerID":"issuer1","balance":75,"originalBalance":100,"status":"PARTIALLY_REDEEMED"}`), nil
    }
    if transactionName == "GetGiftCardHistory" {
        return []byte(`[{"eventID":"1","cardID":"GC1001","eventType":"CARD_CREATED","description":"gift card created by issuer/admin"}]`), nil
    }
    return []byte("[]"), nil
}

func main() {
    var gateway client.Gateway = &mockGateway{}
    issuer := issuerrole.New(gateway)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    msg, _ := issuer.IssueGiftCard(ctx, "GC1001", "issuer1", 100)
    fmt.Println(msg)

    msg, _ = issuer.ActivateGiftCard(ctx, "GC1001")
    fmt.Println(msg)

    msg, _ = issuer.RedeemGiftCard(ctx, "GC1001", 25)
    fmt.Println(msg)

    card, _ := issuer.GetGiftCard(ctx, "GC1001")
    fmt.Printf("Current card: %+v\n", *card)

    history, _ := issuer.GetGiftCardHistory(ctx, "GC1001")
    fmt.Printf("History entries: %d\n", len(history))
}
