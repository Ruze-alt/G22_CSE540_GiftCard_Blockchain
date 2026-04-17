package models

type GiftCard struct {
    CardID          string  `json:"cardId"`
    OwnerID         string  `json:"ownerId"`
    OwnerMSP        string  `json:"ownerMsp"`
    IssuerID        string  `json:"issuerId"`
    IssuerMSP       string  `json:"issuerMsp"`
    RetailerID      string  `json:"retailerId"`
    RetailerMSP     string  `json:"retailerMsp"`
    Balance         float64 `json:"balance"`
    OriginalBalance float64 `json:"originalBalance"`
    Status          string  `json:"status"`
    CreatedAt       string  `json:"createdAt"`
    ActivatedAt     string  `json:"activatedAt"`
    LastUpdatedAt   string  `json:"lastUpdatedAt"`
}

type GiftCardEvent struct {
    EventID   string `json:"eventId"`
    CardID    string `json:"cardId"`
    EventType string `json:"eventType"`
    ActorID   string `json:"actorId"`
    ActorMSP  string `json:"actorMsp"`
    ActorRole string `json:"actorRole"`
    Details   string `json:"details"`
    TxID      string `json:"txId"`
    Timestamp string `json:"timestamp"`
}
